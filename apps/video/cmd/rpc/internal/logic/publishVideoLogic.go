package logic

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/mqueue/cmd/job/jobtype"
	"go-zero-douyin/apps/video/cmd/rpc/internal/model"
	"go-zero-douyin/common/message"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xerr"
	"strings"

	"go-zero-douyin/apps/video/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/video/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PublishVideoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPublishVideoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PublishVideoLogic {
	return &PublishVideoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PublishVideoLogic) PublishVideo(in *pb.PublishVideoReq) (*pb.PublishVideoResp, error) {
	// todo: add your logic here and delete this line
	// 校验输入
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_CHECK_ERR), "Publish video empty param")
	}
	// 属于延迟发布的视频，创建延迟任务，直接返回成功
	if in.GetPublishTime() != 0 {
		if payload, err := json.Marshal(jobtype.DeferPublishVideoPayload{
			Title:     in.GetTitle(),
			OwnerId:   in.GetOwnerId(),
			OwnerName: in.GetOwnerName(),
			SectionID: in.GetSectionId(),
			TagIds:    in.GetTags(),
			PlayUrl:   in.GetPlayUrl(),
			CoverUrl:  in.GetCoverUrl(),
		}); err != nil {
			l.Logger.Errorf("创建发布视频延迟任务时，序列化消息失败, err: %v", err)
		} else {
			_, err = l.svcCtx.Asynq.EnqueueContext(l.ctx, asynq.NewTask(jobtype.DeferPublishVideo, payload), asynq.ProcessAt(utils.FromUnixTimestampToTime(in.GetPublishTime())))
			if err != nil {
				return nil, errors.Wrapf(xerr.NewErrMsg("创建发布视频的延迟任务失败"), "err: %v", err)
			}
		}
		return &pb.PublishVideoResp{}, nil
	}

	// 非延迟任务，直接插入数据
	video := &model.Video{
		Title:     in.GetTitle(),
		OwnerID:   in.GetOwnerId(),
		OwnerName: in.GetOwnerName(),
		SectionID: in.GetSectionId(),
		TagIds:    strings.Join(in.GetTags(), ","),
		PlayURL:   in.GetPlayUrl(),
		CoverURL:  in.GetCoverUrl(),
	}
	err := l.svcCtx.VideoDo.InsertVideo(l.ctx, video)

	// 插入失败
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_INSERT_ERR), "insert video failed, err: %v", err)
	}

	// 发布更新es视频文档的消息
	msg, _ := json.Marshal(message.MysqlVideoUpdateMessage{VideoId: video.ID})
	err = l.svcCtx.Rabbit.Send("", "MysqlVideoUpdateMq", msg)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_UPDATE_ERR), "req: %v, err: %v", in, err)
	}

	// 返回信息
	return &pb.PublishVideoResp{
		Video: &pb.VideoInfo{
			Id:         video.ID,
			Title:      video.Title,
			SectionId:  video.SectionID,
			Tags:       strings.Split(video.TagIds, ","),
			OwnerId:    video.OwnerID,
			OwnerName:  video.OwnerName,
			PlayUrl:    video.PlayURL,
			CoverUrl:   video.CoverURL,
			CreateTime: video.CreateTime,
			UpdateTime: video.UpdateTime,
		},
	}, nil
}

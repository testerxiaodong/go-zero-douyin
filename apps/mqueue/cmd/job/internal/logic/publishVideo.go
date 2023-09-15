package logic

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/mqueue/cmd/job/internal/svc"
	"go-zero-douyin/apps/mqueue/cmd/job/jobtype"
	"go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"
)

var ErrPublishVideo = xerr.NewErrMsg("publish video fail")

// PublishVideoHandler 延迟发布视频路由
type PublishVideoHandler struct {
	svcCtx *svc.ServiceContext
}

func NewPublishVideoHandler(svcCtx *svc.ServiceContext) *PublishVideoHandler {
	return &PublishVideoHandler{
		svcCtx: svcCtx,
	}
}

// ProcessTask defer publish video : if return err != nil , asynq will retry
func (l *PublishVideoHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {

	var p jobtype.DeferPublishVideoPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return errors.Wrapf(ErrPublishVideo, "DeferPublishVideoPayload payload err:%v, payLoad:%+v", err, t.Payload())
	}

	_, err := l.svcCtx.VideoRpc.PublishVideo(ctx, &pb.PublishVideoReq{
		Title: p.Title,
		SectionId: p.SectionID,
		Tags: p.TagIds,
		OwnerId: p.OwnerId,
		OwnerName: p.OwnerName,
		PlayUrl: p.PlayUrl,
		CoverUrl: p.CoverUrl,
	})
	if err != nil {
		return err
	}

	return nil
}

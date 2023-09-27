package like

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"
	pbVideo "go-zero-douyin/apps/video/cmd/rpc/pb"
	"go-zero-douyin/common/ctxdata"
	"go-zero-douyin/common/message"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/social/cmd/api/internal/svc"
	"go-zero-douyin/apps/social/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddLikeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddLikeLogic {
	return &AddLikeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddLikeLogic) AddLike(req *types.VideoLikeReq) error {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := l.svcCtx.Validator.Validate(req); len(validateResult) > 0 {
		return xerr.NewErrMsg(validateResult)
	}
	// 调用videorpc判断视频是否存在
	_, err := l.svcCtx.VideoRpc.GetVideoById(l.ctx, &pbVideo.GetVideoByIdReq{Id: req.VideoId})
	if err != nil {
		return errors.Wrapf(err, "req: %v", req)
	}
	// 获取当前用户id
	uid := ctxdata.GetUidFromCtx(l.ctx)
	// 发布消息到kafka
	data, _ := json.Marshal(message.VideoLikeMessage{UserId: uid, VideoId: req.VideoId})
	err = l.svcCtx.KqueueVideoLikeClient.Push(string(data))
	if err != nil {
		return errors.Wrapf(err, "req: %v", req)
	}
	return nil
}

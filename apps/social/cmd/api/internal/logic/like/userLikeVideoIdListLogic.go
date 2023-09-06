package like

import (
	"context"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"go-zero-douyin/apps/social/cmd/rpc/pb"
	"go-zero-douyin/common/xerr"

	"go-zero-douyin/apps/social/cmd/api/internal/svc"
	"go-zero-douyin/apps/social/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLikeVideoIdListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLikeVideoIdListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLikeVideoIdListLogic {
	return &UserLikeVideoIdListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLikeVideoIdListLogic) UserLikeVideoIdList(req *types.GetUserLikeVideoIdListReq) (resp *types.GetUserLikeVideoIdListResp, err error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if validateResult := l.svcCtx.Validator.Validate(req); len(validateResult) > 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg(validateResult), "req: %v", req)
	}

	// 调用likerpc
	idList, err := l.svcCtx.SocialRpc.GetUserLikeVideoIdList(l.ctx, &pb.GetUserLikeVideoIdListReq{UserId: req.UserId})
	if err != nil {
		return nil, errors.Wrapf(err, "req: %v", req)
	}
	resp = &types.GetUserLikeVideoIdListResp{VideoIdList: make([]int64, 0)}
	err = copier.Copy(resp, idList)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("copy likerpc.GetUserLikeVideoIdListResp resp to api resp failed"), "data: %v, err: %v", idList, err)
	}
	return resp, nil
}

package logic

import (
	"context"
	"github.com/pkg/errors"
	"go-zero-douyin/common/utils"
	"go-zero-douyin/common/xconst"
	"go-zero-douyin/common/xerr"
	"strconv"

	"go-zero-douyin/apps/social/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/social/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommentCountByVideoIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommentCountByVideoIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommentCountByVideoIdLogic {
	return &GetCommentCountByVideoIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommentCountByVideoIdLogic) GetCommentCountByVideoId(in *pb.GetCommentCountByVideoIdReq) (*pb.GetCommentCountByVideoIdResp, error) {
	// todo: add your logic here and delete this line
	// 参数校验
	if in == nil {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get video comment count with empty param")
	}
	if in.GetVideoId() == 0 {
		return nil, errors.Wrap(xerr.NewErrCode(xerr.PB_LOGIC_CHECK_ERR), "get video comment count with empty video_id")
	}

	// 从redis中获取数据
	result, err := l.svcCtx.Redis.ExistsCtx(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisVideoCommentPrefix, in.GetVideoId()))
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "get redis video comment count key exist failed: %v", err)
	}
	if result == true {
		val, err := l.svcCtx.Redis.GetCtx(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisVideoCommentPrefix, in.GetVideoId()))
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_SEARCH_ERR), "get redis video comment count failed: %v", err)
		}
		count, err := strconv.Atoi(val)
		if err != nil {
			return nil, errors.Wrapf(xerr.NewErrCode(xerr.RPC_INSERT_ERR), "redis value invalid, except int, err: %v", err)
		}
		return &pb.GetCommentCountByVideoIdResp{Count: int64(count)}, nil
	}

	// 从mysql中获取数据
	key := strconv.Itoa(int(in.GetVideoId()))
	count, err := l.svcCtx.SingleFlight.Do(key, func() (any, error) {
		return l.GetVideoCommentFromDb(in.GetVideoId())
	})
	if err != nil {
		return nil, err
	}
	countInt64, ok := count.(int64)
	if !ok {
		return &pb.GetCommentCountByVideoIdResp{}, errors.Wrap(xerr.NewErrCode(xerr.PB_CHECK_ERR), "type assert failed")
	}

	// 重新构建缓存
	l.BuildVideoCommentCountCache(in.VideoId, countInt64)

	return &pb.GetCommentCountByVideoIdResp{Count: countInt64}, nil
}

func (l *GetCommentCountByVideoIdLogic) GetVideoCommentFromDb(videoId int64) (int64, error) {
	commentQuery := l.svcCtx.Query.Comment
	count, err := commentQuery.WithContext(l.ctx).Where(commentQuery.VideoID.Eq(videoId)).Count()
	if err != nil {
		return 0, errors.Wrapf(xerr.NewErrCode(xerr.DB_SEARCH_ERR), "get mysql video commnet count failed: %v", err)
	}
	return count, nil
}

func (l *GetCommentCountByVideoIdLogic) BuildVideoCommentCountCache(videoId int64, commentCount int64) {
	commentCountStr := strconv.Itoa(int(commentCount))
	err := l.svcCtx.Redis.SetCtx(l.ctx, utils.GetRedisKeyWithPrefix(xconst.RedisVideoCommentPrefix, videoId), commentCountStr)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("set video comment cache failed, video_id: %d", videoId)
		return
	}
}

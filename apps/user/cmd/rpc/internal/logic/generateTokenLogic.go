package logic

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"go-zero-douyin/common/ctxdata"
	"go-zero-douyin/common/xerr"
	"time"

	"go-zero-douyin/apps/user/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/user/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrGenerateTokenError = xerr.NewErrMsg("生成token失败")

type GenerateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateTokenLogic {
	return &GenerateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GenerateTokenLogic) GenerateToken(in *pb.GenerateTokenReq) (*pb.GenerateTokenResp, error) {
	// todo: add your logic here and delete this line
	now := time.Now().Unix()
	accessExpire := l.svcCtx.Config.JwtAuth.AccessExpire
	accessToken, err := l.getJwtToken(l.svcCtx.Config.JwtAuth.AccessSecret, now, accessExpire, in.GetUserId())
	if err != nil {
		return nil, errors.Wrapf(ErrGenerateTokenError, "Generate token failed, user_id: %d", in.GetUserId())
	}
	return &pb.GenerateTokenResp{
		AccessToken:  accessToken,
		ExpireTime:   now + accessExpire,
		RefreshAfter: now + accessExpire/2,
	}, nil
}

func (l *GenerateTokenLogic) getJwtToken(secretKey string, iat, expire int64, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + expire
	claims["iat"] = iat
	claims[ctxdata.CtxKeyJwtUserId] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

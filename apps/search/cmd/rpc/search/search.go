// Code generated by goctl. DO NOT EDIT.
// Source: search.proto

package search

import (
	"context"

	"go-zero-douyin/apps/search/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	CompleteVideoReq  = pb.CompleteVideoReq
	CompleteVideoResp = pb.CompleteVideoResp
	SearchUserReq     = pb.SearchUserReq
	SearchUserResp    = pb.SearchUserResp
	SearchVideoReq    = pb.SearchVideoReq
	SearchVideoResp   = pb.SearchVideoResp
	User              = pb.User
	Video             = pb.Video

	Search interface {
		// 视频相关功能
		SearchVideo(ctx context.Context, in *SearchVideoReq, opts ...grpc.CallOption) (*SearchVideoResp, error)
		CompleteVideo(ctx context.Context, in *CompleteVideoReq, opts ...grpc.CallOption) (*CompleteVideoResp, error)
		// 用户相关功能
		SearchUser(ctx context.Context, in *SearchUserReq, opts ...grpc.CallOption) (*SearchUserResp, error)
	}

	defaultSearch struct {
		cli zrpc.Client
	}
)

func NewSearch(cli zrpc.Client) Search {
	return &defaultSearch{
		cli: cli,
	}
}

// 视频相关功能
func (m *defaultSearch) SearchVideo(ctx context.Context, in *SearchVideoReq, opts ...grpc.CallOption) (*SearchVideoResp, error) {
	client := pb.NewSearchClient(m.cli.Conn())
	return client.SearchVideo(ctx, in, opts...)
}

func (m *defaultSearch) CompleteVideo(ctx context.Context, in *CompleteVideoReq, opts ...grpc.CallOption) (*CompleteVideoResp, error) {
	client := pb.NewSearchClient(m.cli.Conn())
	return client.CompleteVideo(ctx, in, opts...)
}

// 用户相关功能
func (m *defaultSearch) SearchUser(ctx context.Context, in *SearchUserReq, opts ...grpc.CallOption) (*SearchUserResp, error) {
	client := pb.NewSearchClient(m.cli.Conn())
	return client.SearchUser(ctx, in, opts...)
}

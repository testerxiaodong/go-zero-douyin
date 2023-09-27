// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: rpc/pb/recommend.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Recommend_VideoRecommendSection_FullMethodName = "/pb.recommend/VideoRecommendSection"
)

// RecommendClient is the client API for Recommend service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RecommendClient interface {
	VideoRecommendSection(ctx context.Context, in *VideoRecommendSectionReq, opts ...grpc.CallOption) (*VideoRecommendSectionResp, error)
}

type recommendClient struct {
	cc grpc.ClientConnInterface
}

func NewRecommendClient(cc grpc.ClientConnInterface) RecommendClient {
	return &recommendClient{cc}
}

func (c *recommendClient) VideoRecommendSection(ctx context.Context, in *VideoRecommendSectionReq, opts ...grpc.CallOption) (*VideoRecommendSectionResp, error) {
	out := new(VideoRecommendSectionResp)
	err := c.cc.Invoke(ctx, Recommend_VideoRecommendSection_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RecommendServer is the server API for Recommend service.
// All implementations must embed UnimplementedRecommendServer
// for forward compatibility
type RecommendServer interface {
	VideoRecommendSection(context.Context, *VideoRecommendSectionReq) (*VideoRecommendSectionResp, error)
	mustEmbedUnimplementedRecommendServer()
}

// UnimplementedRecommendServer must be embedded to have forward compatible implementations.
type UnimplementedRecommendServer struct {
}

func (UnimplementedRecommendServer) VideoRecommendSection(context.Context, *VideoRecommendSectionReq) (*VideoRecommendSectionResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VideoRecommendSection not implemented")
}
func (UnimplementedRecommendServer) mustEmbedUnimplementedRecommendServer() {}

// UnsafeRecommendServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RecommendServer will
// result in compilation errors.
type UnsafeRecommendServer interface {
	mustEmbedUnimplementedRecommendServer()
}

func RegisterRecommendServer(s grpc.ServiceRegistrar, srv RecommendServer) {
	s.RegisterService(&Recommend_ServiceDesc, srv)
}

func _Recommend_VideoRecommendSection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VideoRecommendSectionReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecommendServer).VideoRecommendSection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Recommend_VideoRecommendSection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecommendServer).VideoRecommendSection(ctx, req.(*VideoRecommendSectionReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Recommend_ServiceDesc is the grpc.ServiceDesc for Recommend service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Recommend_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.recommend",
	HandlerType: (*RecommendServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "VideoRecommendSection",
			Handler:    _Recommend_VideoRecommendSection_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpc/pb/recommend.proto",
}
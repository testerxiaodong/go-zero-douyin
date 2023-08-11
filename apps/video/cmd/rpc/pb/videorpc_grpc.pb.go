// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: rpc/pb/videorpc.proto

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
	Videorpc_PublishVideo_FullMethodName = "/pb.videorpc/PublishVideo"
)

// VideorpcClient is the client API for Videorpc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VideorpcClient interface {
	PublishVideo(ctx context.Context, in *PublishVideoReq, opts ...grpc.CallOption) (*PublishVideoResp, error)
}

type videorpcClient struct {
	cc grpc.ClientConnInterface
}

func NewVideorpcClient(cc grpc.ClientConnInterface) VideorpcClient {
	return &videorpcClient{cc}
}

func (c *videorpcClient) PublishVideo(ctx context.Context, in *PublishVideoReq, opts ...grpc.CallOption) (*PublishVideoResp, error) {
	out := new(PublishVideoResp)
	err := c.cc.Invoke(ctx, Videorpc_PublishVideo_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VideorpcServer is the server API for Videorpc service.
// All implementations must embed UnimplementedVideorpcServer
// for forward compatibility
type VideorpcServer interface {
	PublishVideo(context.Context, *PublishVideoReq) (*PublishVideoResp, error)
	mustEmbedUnimplementedVideorpcServer()
}

// UnimplementedVideorpcServer must be embedded to have forward compatible implementations.
type UnimplementedVideorpcServer struct {
}

func (UnimplementedVideorpcServer) PublishVideo(context.Context, *PublishVideoReq) (*PublishVideoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishVideo not implemented")
}
func (UnimplementedVideorpcServer) mustEmbedUnimplementedVideorpcServer() {}

// UnsafeVideorpcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VideorpcServer will
// result in compilation errors.
type UnsafeVideorpcServer interface {
	mustEmbedUnimplementedVideorpcServer()
}

func RegisterVideorpcServer(s grpc.ServiceRegistrar, srv VideorpcServer) {
	s.RegisterService(&Videorpc_ServiceDesc, srv)
}

func _Videorpc_PublishVideo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishVideoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VideorpcServer).PublishVideo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Videorpc_PublishVideo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VideorpcServer).PublishVideo(ctx, req.(*PublishVideoReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Videorpc_ServiceDesc is the grpc.ServiceDesc for Videorpc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Videorpc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.videorpc",
	HandlerType: (*VideorpcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PublishVideo",
			Handler:    _Videorpc_PublishVideo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpc/pb/videorpc.proto",
}
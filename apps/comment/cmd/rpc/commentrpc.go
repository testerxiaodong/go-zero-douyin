package main

import (
	"flag"
	"fmt"
	"go-zero-douyin/common/interceptor/rpcServer"

	"go-zero-douyin/apps/comment/cmd/rpc/internal/config"
	"go-zero-douyin/apps/comment/cmd/rpc/internal/server"
	"go-zero-douyin/apps/comment/cmd/rpc/internal/svc"
	"go-zero-douyin/apps/comment/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/commentrpc.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterCommentrpcServer(grpcServer, server.NewCommentrpcServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	//rpc log
	s.AddUnaryInterceptors(rpcServer.ErrTransLogInterceptor)

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

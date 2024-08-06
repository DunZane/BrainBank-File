package main

import (
	"flag"
	"fmt"

	"github.com/dunzane/brainbank-file/rpc/fileInfo/internal/config"
	"github.com/dunzane/brainbank-file/rpc/fileInfo/internal/server"
	"github.com/dunzane/brainbank-file/rpc/fileInfo/internal/svc"
	"github.com/dunzane/brainbank-file/rpc/fileInfo/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/fileinfo.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterFileInfoServer(grpcServer, server.NewFileInfoServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

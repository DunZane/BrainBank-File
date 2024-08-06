package svc

import (
	"github.com/dunzane/brainbank-file/api/internal/config"
	"github.com/dunzane/brainbank-file/api/internal/middleware"
	"github.com/dunzane/brainbank-file/rpc/fileInfo/fileinfo"
	"github.com/dunzane/brainbank-file/rpc/fileOps/fileops"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config        config.Config
	JwtMiddleware rest.Middleware
	FileInfo      fileinfo.FileInfo
	FileOps       fileops.FileOps
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		JwtMiddleware: middleware.NewJwtMiddleware(c).Handle,
		FileInfo:      fileinfo.NewFileInfo(zrpc.MustNewClient(c.FileInfo)),
		FileOps:       fileops.NewFileOps(zrpc.MustNewClient(c.FileOps)),
	}
}

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	FileInfo zrpc.RpcClientConf
	FileOps  zrpc.RpcClientConf
	Auth     struct {
		AccessSecret string
		AccessExpire int64
	}
}

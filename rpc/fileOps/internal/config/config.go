package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Datasource string
	Cache      cache.CacheConf
	Minio      MinioConf    `json:"Minio"`
	RabbitMQ   RabbitMQConf `json:"RabbitMQ"`
}

type MinioConf struct {
	EndPoint           string `json:"Endpoint"`
	AccessKeyID        string `json:"AccessKeyID"`
	SecretAccessKey    string `json:"SecretAccessKey"`
	PresignedURLExpire int    `json:"PresignedURLExpire"`
	BucketName         string `json:"BucketName"`
}

type RabbitMQConf struct {
	Protocol string `json:"Protocol"`
	Username string `json:"Username"`
	Password string `json:"Password"`
	Host     string `json:"Host"`
	Port     int    `json:"Port"`
}

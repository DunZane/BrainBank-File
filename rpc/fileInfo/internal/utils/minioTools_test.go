package utils

import (
	"github.com/dunzane/brainbank-file/rpc/fileInfo/internal/config"
	"github.com/dunzane/brainbank-file/rpc/fileInfo/internal/svc"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zeromicro/go-zero/zrpc"
	"testing"
	"time"
)

func NewMinioClient(c config.Config) *minio.Client {
	endpoint := c.Minio.EndPoint
	accessKeyID := c.Minio.AccessKeyID
	secretAccessKey := c.Minio.SecretAccessKey

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: false,
	})

	if err != nil {
		panic("err init minio client")
	}

	return minioClient
}

var c = config.Config{
	RpcServerConf: zrpc.RpcServerConf{},
	Datasource:    "",
	Minio: config.MinioConf{
		EndPoint:           "127.0.0.1:9000",
		AccessKeyID:        "TjWQtWLO23PfqNCs2Hmg",
		SecretAccessKey:    "3hueJG1GVrPubRPiqZXdZatfr8HygRZowyW5aqfg",
		PresignedURLExpire: 3600,
		BucketName:         "file",
	},
}

var ctx = svc.ServiceContext{
	Config: c,
	Minio:  NewMinioClient(c),
}

func TestMinioTools_PresignedPutObject(t *testing.T) {
	m := MinioTools{}

	userId := 37461266
	fileId := "test.pdf"

	url, expires, err := m.PresignedGetObject(&ctx, int64(userId), fileId)
	if err != nil {
		t.Errorf("err to get the oeration of put presigned url:{%v}", err)
	}

	expiresAt := time.Now().Add(expires)

	t.Logf("success to get the oeration of put presigned url:{%s},expires at{%s}",
		url.String(), expiresAt.Format("2006-01-02 15:04:05"))
}

package svc

import (
	"github.com/dunzane/brainbank-file/rpc/fileInfo/internal/config"
	"github.com/dunzane/brainbank-file/rpc/model"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config    config.Config
	FileModel model.FileModel
	Minio     *minio.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		FileModel: model.NewFileModel(sqlx.NewMysql(c.Datasource), c.Cache),
		Minio:     NewMinioClient(c),
	}
}

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

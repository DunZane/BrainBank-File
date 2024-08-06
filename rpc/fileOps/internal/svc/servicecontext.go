package svc

import (
	"github.com/dunzane/brainbank-file/rpc/fileOps/internal/config"
	"github.com/dunzane/brainbank-file/rpc/model"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"strconv"
	"strings"
)

type ServiceContext struct {
	Config    config.Config
	FileModel model.FileModel
	Minio     *minio.Client
	MQConn    *amqp.Connection
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		FileModel: model.NewFileModel(sqlx.NewMysql(c.Datasource), c.Cache),
		Minio:     NewMinioClient(c),
		MQConn:    NewRabbitMQConn(c),
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

func NewRabbitMQConn(c config.Config) *amqp.Connection {
	var builder strings.Builder
	builder.WriteString(c.RabbitMQ.Protocol)
	builder.WriteString("://")
	builder.WriteString(c.RabbitMQ.Username)
	builder.WriteString(":")
	builder.WriteString(c.RabbitMQ.Password)
	builder.WriteString("@")
	builder.WriteString(c.RabbitMQ.Host)
	builder.WriteString(":")
	builder.WriteString(strconv.Itoa(c.RabbitMQ.Port))
	dns := builder.String()

	conn, err := amqp.Dial(dns)
	if err != nil {
		panic("err dial to rabbitmq")
	}

	return conn
}

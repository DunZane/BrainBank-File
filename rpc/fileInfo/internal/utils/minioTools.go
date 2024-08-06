package utils

import (
	"context"
	"github.com/dunzane/brainbank-file/rpc/fileInfo/internal/svc"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type MinioTools struct{}

func (mt *MinioTools) PresignedGetObject(svcCtx *svc.ServiceContext,
	userID int64, fileId string) (*url.URL, time.Duration, error) {
	// Retrieve MinIO configuration
	minioConfig := svcCtx.Config.Minio
	bucketName := minioConfig.BucketName
	objectPath := strings.Join([]string{strconv.FormatInt(userID, 10), fileId}, "/")

	// Set the expiration time for the presigned URL
	urlExpire := minioConfig.PresignedURLExpire
	expiration := time.Duration(urlExpire) * time.Second

	// Set request parameters for response headers
	reqParams := make(url.Values)
	reqParams.Set("response-cache-control", "no-cache")
	reqParams.Set("response-content-type", "application/octet-stream")

	// Generate the presigned URL
	presignedURL, err := svcCtx.Minio.PresignedGetObject(context.Background(),
		bucketName, objectPath, expiration, reqParams)
	if err != nil {
		return nil, 0, err
	}
	return presignedURL, expiration, nil
}

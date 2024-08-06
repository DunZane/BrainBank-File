package utils

import (
	"context"
	"github.com/dunzane/brainbank-file/rpc/fileOps/internal/svc"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// MinioTools provides methods for interacting with MinIO.
type MinioTools struct{}

// PresignedPutObject generates a presigned URL for uploading an object to MinIO.
func (mt *MinioTools) PresignedPutObject(svcCtx *svc.ServiceContext,
	userID int64, fileName string) (*url.URL, time.Duration, error) {
	// Retrieve MinIO configuration
	minioConfig := svcCtx.Config.Minio
	bucketName := minioConfig.BucketName
	objectPath := strings.Join([]string{strconv.FormatInt(userID, 10), fileName}, "/")

	// Set the expiration time for the presigned URL
	urlExpire := minioConfig.PresignedURLExpire
	expiration := time.Duration(urlExpire) * time.Second

	// Generate the presigned URL
	presignedURL, err := svcCtx.Minio.PresignedPutObject(context.Background(), bucketName, objectPath, expiration)
	if err != nil {
		return nil, 0, err
	}

	return presignedURL, expiration, nil
}

func (mt *MinioTools) PresignedGetObject(svcCtx *svc.ServiceContext,
	userID, fileName string) (*url.URL, time.Duration, error) {
	// Retrieve MinIO configuration
	minioConfig := svcCtx.Config.Minio
	bucketName := minioConfig.BucketName
	objectPath := strings.Join([]string{userID, fileName}, "/")

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

func (mt *MinioTools) PresignedCopyObject(svcCtx *svc.ServiceContext,
	userID int64, fileName, srcObject string) (*url.URL, time.Duration, error) {
	// Retrieve MinIO configuration
	minioConfig := svcCtx.Config.Minio
	bucketName := minioConfig.BucketName
	objectPath := strings.Join([]string{strconv.FormatInt(userID, 10), fileName}, "/")

	// Set the expiration time for the presigned URL
	urlExpire := minioConfig.PresignedURLExpire
	expiration := time.Duration(urlExpire) * time.Second

	// Set request parameters for COPY operation
	reqParams := url.Values{}
	reqParams.Set("x-amz-copy-source", "/"+bucketName+"/"+srcObject)

	// Generate the presigned URL for COPY operation
	presignedURL, err := svcCtx.Minio.Presign(context.Background(), http.MethodPut,
		bucketName, objectPath, expiration, reqParams)
	if err != nil {
		return nil, 0, err
	}

	return presignedURL, expiration, nil
}

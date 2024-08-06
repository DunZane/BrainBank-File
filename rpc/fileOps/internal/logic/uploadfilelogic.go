package logic

import (
	"context"
	"database/sql"
	"errors"
	"github.com/dunzane/brainbank-file/rpc/fileOps/fileops"
	"github.com/dunzane/brainbank-file/rpc/fileOps/internal/constant"
	"github.com/dunzane/brainbank-file/rpc/fileOps/internal/svc"
	"github.com/dunzane/brainbank-file/rpc/fileOps/internal/utils"
	"github.com/dunzane/brainbank-file/rpc/fileOps/pb"
	"github.com/dunzane/brainbank-file/rpc/model"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileLogic {
	return &UploadFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UploadFileLogic) UploadFile(in *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	// 默认状态返回成功
	resp := &fileops.UploadFileResponse{Code: constant.StatusOK}
	tools := utils.MinioTools{}

	// 获取文件信息
	userId := in.UserId
	fileName := in.File.FileName
	fileMd5 := in.File.FileMd5
	l.Logger.Infof("UploadFile request received - UserID: %d, FileName: %s, FileMD5: %s", userId, fileName, fileMd5)

	// 校验文件
	var presignedURL *url.URL
	var expires time.Duration
	file, err := l.svcCtx.FileModel.FindOneByMD5(l.ctx, fileMd5)
	if err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			l.Logger.Errorf("Error occurred while finding file by MD5: %v", err)
			resp.Code = constant.StatusInternalServerError
			return resp, nil
		}

		// 返回put操作
		l.Logger.Infof("File with MD5 %s not found. Generating presigned PUT URL for file: %s", fileMd5, fileName)
		presignedURL, expires, err = tools.PresignedPutObject(l.svcCtx, userId, fileName)
		if err != nil {
			l.Logger.Errorf("Error when trying to generate a presigned PUT URL: %v", err)
			resp.Code = constant.ErrorGenPresignedURL
			return resp, nil
		}

	} else {
		l.Logger.Infof("File with MD5 %s found. Generating presigned COPY URL for file: %s", fileMd5, fileName)

		// 目标文件信息
		sourceObject := strings.Join([]string{strconv.Itoa(int(file.OwnerId)), file.Name}, "/")
		l.Logger.Infof("Source object path for COPY operation: %s", sourceObject)

		// 返回copy操作
		presignedURL, expires, err = tools.PresignedCopyObject(l.svcCtx, userId, fileName, sourceObject)
		if err != nil {
			l.Logger.Errorf("Error when trying to generate a presigned COPY URL: %v", err)
			resp.Code = constant.ErrorGenPresignedURL
			return resp, nil
		}
	}

	l.Logger.Infof("Returning presigned URL: %s", presignedURL.String())

	// 插入文件对象到数据库 todo：数据放在缓存中
	idGen := utils.IdGenerator{}
	fileId := idGen.GenerateFileID()
	m := &model.File{
		Id:       fileId,
		Name:     fileName,
		Type:     in.File.FileType,
		Size:     in.File.FileSize,
		Path:     strings.Join([]string{strconv.FormatInt(userId, 10), "/", fileId}, ""),
		OwnerId:  userId,
		Status:   "pending",
		Checksum: sql.NullString{String: fileMd5, Valid: fileMd5 != ""},
	}

	_, err = l.svcCtx.FileModel.Insert(l.ctx, m)
	if err != nil {
		l.Logger.Errorf("Error occurred while inserting a file object: %v", err)
		resp.Code = constant.StatusInternalServerError
		return resp, nil
	}

	// 绑定返回数据
	resp.PresignedUrl = presignedURL.String()
	expiresAt := time.Now().Add(expires)
	resp.ExpiresAt = expiresAt.Format("2006-01-02 15:04:05")
	resp.FileId = fileId
	return resp, nil
}

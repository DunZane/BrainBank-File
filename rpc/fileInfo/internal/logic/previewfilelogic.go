package logic

import (
	"context"
	"errors"
	"github.com/dunzane/brainbank-file/rpc/fileInfo/fileinfo"
	"github.com/dunzane/brainbank-file/rpc/fileInfo/internal/constant"
	"github.com/dunzane/brainbank-file/rpc/fileInfo/internal/utils"
	"github.com/dunzane/brainbank-file/rpc/model"
	"time"

	"github.com/dunzane/brainbank-file/rpc/fileInfo/internal/svc"
	"github.com/dunzane/brainbank-file/rpc/fileInfo/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PreviewFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPreviewFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreviewFileLogic {
	return &PreviewFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PreviewFileLogic) PreviewFile(in *pb.PreviewFileRequest) (*pb.PreviewFileResponse, error) {
	// 默认状态返回成功
	resp := &fileinfo.PreviewFileResponse{Code: constant.StatusOK}
	tools := utils.MinioTools{}

	l.Logger.Infof("Received PreviewFile request: %+v", in)

	// 查询出该文件记录
	file, err := l.svcCtx.FileModel.FindOne(l.ctx, in.FileId)
	if err != nil {
		// 查询出错
		if !errors.Is(err, model.ErrNotFound) {
			l.Logger.Errorf("Error occurred while finding file with file_id %s: %v", in.FileId, err)
			resp.Code = constant.StatusInternalServerError
			return resp, nil
		}

		// 空记录
		l.Logger.Errorf("File with file_id %s not found, proceeding with update operation", in.FileId)
		resp.Code = constant.FileNotFound
		return resp, nil
	}
	l.Logger.Infof("Found file record: %+v", file)

	// 校验身份
	if file.OwnerId != in.UserId {
		l.Logger.Errorf("User with user_id %s is not the owner of the file_id %s", in.UserId, in.FileId)
		resp.Code = constant.FileNotFound
		return resp, nil
	}

	presignedURL, expires, err := tools.PresignedGetObject(l.svcCtx, in.UserId, in.FileId)
	if err != nil {
		l.Logger.Errorf("Error when trying to generate a presigned GET URL: %v", err)
		resp.Code = constant.ErrorGenPresignedURL
		return nil, err
	}

	resp.PreSignedUrl = presignedURL.String()
	expiresAt := time.Now().Add(expires)
	resp.ExpiresAt = expiresAt.Format("2006-01-02 15:04:05")
	return resp, nil
}

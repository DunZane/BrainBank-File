package logic

import (
	"context"
	"errors"
	"github.com/dunzane/brainbank-file/rpc/fileOps/fileops"
	"github.com/dunzane/brainbank-file/rpc/fileOps/internal/constant"
	"github.com/dunzane/brainbank-file/rpc/model"

	"github.com/dunzane/brainbank-file/rpc/fileOps/internal/svc"
	"github.com/dunzane/brainbank-file/rpc/fileOps/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileLogic {
	return &DeleteFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteFileLogic) DeleteFile(in *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	// 默认状态返回成功
	resp := &fileops.DeleteFileResponse{Code: constant.StatusOK}

	l.Logger.Infof("Received DeleteFileRequest: %+v", in)

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

	// 删除文件
	file.Status = "deleted"
	err = l.svcCtx.FileModel.Update(l.ctx, file)
	if err != nil {
		l.Logger.Errorf("Error occurred while updating file record with file_id %s: %v", in.FileId, err)
		resp.Code = constant.StatusInternalServerError
		return resp, nil
	}

	l.Logger.Infof("Successfully marked file with file_id %s as deleted", in.FileId)
	return resp, nil
}

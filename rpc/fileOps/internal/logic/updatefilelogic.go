package logic

import (
	"context"
	"errors"
	"github.com/dunzane/brainbank-file/rpc/fileOps/fileops"
	"github.com/dunzane/brainbank-file/rpc/fileOps/internal/constant"
	"github.com/dunzane/brainbank-file/rpc/fileOps/internal/utils"
	"github.com/dunzane/brainbank-file/rpc/model"

	"github.com/dunzane/brainbank-file/rpc/fileOps/internal/svc"
	"github.com/dunzane/brainbank-file/rpc/fileOps/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFileLogic {
	return &UpdateFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateFileLogic) UpdateFile(in *pb.UpdateFileRequest) (*pb.UpdateFileResponse, error) {
	// 默认状态返回成功
	resp := &fileops.UpdateFileResponse{Code: constant.StatusOK}

	// 查询出该文件记录
	file, err := l.svcCtx.FileModel.FindOne(l.ctx, in.FileId)
	if err != nil {
		// 查询出错
		if !errors.Is(err, model.ErrNotFound) {
			l.Logger.Errorf("Error occurred while finding file by file_id: %v", err)
			resp.Code = constant.StatusInternalServerError
			return resp, nil
		}

		// 空记录
		l.Logger.Errorf("File with file_id %s not found, proceeding with update operation", in.FileId)
		resp.Code = constant.FileNotFound
		return resp, nil
	} else {

		// 等待被更新的
		l.Logger.Infof("Found file record: %+v", file)
	}

	// 更新对象
	file.StorageProvider = in.StoreProvider
	file.Status = in.Status
	l.Logger.Infof("Updated file object: %+v", file)

	// 新对象重新入库
	err = l.svcCtx.FileModel.Update(l.ctx, file)
	if err != nil {
		// 更新出错
		l.Logger.Errorf("Error occurred while updating file record: %v", err)
		resp.Code = constant.StatusInternalServerError
		return resp, nil
	}

	l.Logger.Infof("Successfully updated file record with file_id: %s", in.FileId)

	// 获取消息队列的tools
	tool, err := utils.NewRabbitMQTools(l.svcCtx)
	if err != nil {
		resp.Code = constant.StatusInternalServerError
		l.Logger.Errorf("Error occurred while initializing RabbitMQ tools: %v", err)
		return resp, nil
	}

	// 消息队列
	err = tool.SendMessage(in.UserId, in.FileId)
	if err != nil {
		resp.Code = constant.StatusInternalServerError
		l.Logger.Errorf("Error occurred while sending message to RabbitMQ: %v", err)
		return resp, err
	}

	l.Logger.Infof("Message sent to RabbitMQ for user_id: %d and file_id: %s", in.UserId, in.FileId)
	return resp, nil
}

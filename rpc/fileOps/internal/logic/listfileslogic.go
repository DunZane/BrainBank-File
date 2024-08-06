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

type ListFilesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListFilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFilesLogic {
	return &ListFilesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListFilesLogic) ListFiles(in *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	// 默认状态返回成功
	resp := &fileops.ListFilesResponse{Code: constant.StatusOK}

	l.Logger.Infof("Received ListFilesRequest: %+v", in)

	// 获取总数
	count, err := l.svcCtx.FileModel.CountFiles(l.ctx, in.UserId)
	if err != nil {
		l.Logger.Errorf("Error occurred while counting files for userId %s: %v", in.UserId, err)
		resp.Code = constant.StatusInternalServerError
		return resp, nil
	}
	l.Logger.Infof("Total file count for userId %d: %d", in.UserId, count)

	// 绑定总数
	resp.Total = int32(count)
	if resp.Total == 0 {
		l.Logger.Infof("No files found for userId %s", in.UserId)
		return resp, nil
	}

	// 列表查询
	files, err := l.svcCtx.FileModel.ListFiles(l.ctx, in.UserId, int(in.Limit), int(in.Offset))
	if err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			l.Logger.Errorf("Error occurred while listing files for userId %d: %v", in.UserId, err)
			resp.Code = constant.StatusInternalServerError
			return resp, nil
		}
		l.Logger.Infof("No files found for userId %d with limit %d and offset %d", in.UserId, in.Limit, in.Offset)
	}
	l.Logger.Infof("Number of files retrieved for userId %d: %d", in.UserId, len(files))

	// 绑定返回数据
	data := make([]*fileops.FileObject, 0, len(files))
	for _, file := range files {
		if file.Status != "active" {
			continue
		}
		data = append(data, &fileops.FileObject{
			FileName: file.Name,
			FileMd5:  file.Checksum.String,
			FileType: file.Type,
			FileSize: file.Size,
		})
	}
	resp.Files = data

	l.Logger.Infof("Returning %d files for userId %d", len(data), in.UserId)
	return resp, nil
}

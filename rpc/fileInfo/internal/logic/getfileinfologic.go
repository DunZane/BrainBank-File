package logic

import (
	"context"
	"errors"
	"github.com/dunzane/brainbank-file/rpc/fileInfo/fileinfo"
	"github.com/dunzane/brainbank-file/rpc/fileInfo/internal/constant"
	"github.com/dunzane/brainbank-file/rpc/model"

	"github.com/dunzane/brainbank-file/rpc/fileInfo/internal/svc"
	"github.com/dunzane/brainbank-file/rpc/fileInfo/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFileInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFileInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFileInfoLogic {
	return &GetFileInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetFileInfoLogic) GetFileInfo(in *pb.GetFileInfoRequest) (*pb.GetFileInfoResponse, error) {
	// 默认状态返回成功
	resp := &fileinfo.GetFileInfoResponse{Code: constant.StatusOK}

	l.Logger.Infof("Received GetFileInfo request: %+v", in)

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

	// 返回数据
	resp.File = &pb.FileObject{
		FileName: file.Name,
		FileMd5:  file.Checksum.String,
		FileType: file.Type,
		FileSize: file.Size,
	}

	l.Logger.Infof("File info successfully bound to response: %+v", resp.File)
	return resp, nil
}

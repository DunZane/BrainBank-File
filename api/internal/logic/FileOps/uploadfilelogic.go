package FileOps

import (
	"context"
	"errors"
	"fmt"
	"github.com/dunzane/brainbank-file/api/internal/constant"
	"github.com/dunzane/brainbank-file/rpc/fileOps/fileops"

	"github.com/dunzane/brainbank-file/api/internal/svc"
	"github.com/dunzane/brainbank-file/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadFileLogic {
	return &UploadFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadFileLogic) UploadFile(req *types.UploadFileRequst) (resp *types.UploadFileResponse, err error) {
	// 默认请求成功
	resp = &types.UploadFileResponse{}
	resp.Code = constant.StatusOK
	resp.Msg = constant.StatusText[resp.Code]

	// 解析JWT，获取关键数据
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		err = errors.New("can't get user info from JWT token which store in context")
		l.Logger.Error(err)
		return nil, err
	}
	l.Logger.Infof("Parsed userID from JWT token: %d", userId)

	// 创建上传文件rpc请求体
	rpcReq := &fileops.UploadFileRequest{
		File: &fileops.FileObject{
			FileName: req.File.FileName,
			FileMd5:  req.File.FileMD5,
			FileType: req.File.FileType,
			FileSize: req.File.FileSize,
		},
		UserId: userId,
	}

	// 调用远程函数
	l.Logger.Debugf("Preparing to call UploadFile RPC with request body: %+v", rpcReq)
	rpcResp, err := l.svcCtx.FileOps.UploadFile(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Error calling UploadFile RPC: %v", err)
		return nil, fmt.Errorf("file.ops.rpc return err:%w", err)
	}
	l.Logger.Debugf("Received response from UploadFile RPC: %+v", rpcResp)

	// 检查请求是否满足预期
	if rpcResp.Code != constant.StatusOK {
		l.Logger.Errorf("UploadFile RPC request unsuccessful, code: %d, message: %s",
			rpcResp.Code, constant.StatusText[int(rpcResp.Code)])
		return nil, fmt.Errorf("UploadFile RPC request unsuccessful, code: %d, message: %s",
			rpcResp.Code, constant.StatusText[int(rpcResp.Code)])
	}

	// 绑定返回数据
	resp.Code = int(rpcResp.Code)
	resp.Msg = constant.StatusText[resp.Code]
	resp.PresignedURL = rpcResp.PresignedUrl
	resp.ExpiresAt = rpcResp.ExpiresAt
	resp.FileId = rpcResp.FileId
	l.Logger.Infof("Upload file request successful, received presigned URL: %s", resp.PresignedURL)
	return
}

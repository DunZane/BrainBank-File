package FileInfo

import (
	"context"
	"errors"
	"fmt"
	"github.com/dunzane/brainbank-file/api/internal/constant"
	"github.com/dunzane/brainbank-file/rpc/fileInfo/fileinfo"

	"github.com/dunzane/brainbank-file/api/internal/svc"
	"github.com/dunzane/brainbank-file/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileInfoLogic {
	return &FileInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileInfoLogic) FileInfo(req *types.GetFileInfoRequest) (resp *types.GetFileInfoResponse, err error) {
	// 默认请求成功
	resp = &types.GetFileInfoResponse{}
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

	// 创建上传文件rpc请求
	rpcReq := &fileinfo.GetFileInfoRequest{
		FileId: req.FileId,
		UserId: userId,
	}

	// 调用远程函数
	l.Logger.Debugf("Preparing to call FileInfo RPC with request body: %+v", rpcReq)
	rpcResp, err := l.svcCtx.FileInfo.GetFileInfo(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Error calling FileInfo RPC: %v", err)
		return nil, fmt.Errorf("file.info.rpc return err:%w", err)
	}
	l.Logger.Debugf("Received response from FileInfo RPC: %+v", rpcResp)

	// 检查请求是否满足预期
	if rpcResp.Code != constant.StatusOK {
		l.Logger.Errorf("FileInfo RPC request unsuccessful, code: %d, message: %s",
			rpcResp.Code, constant.StatusText[int(rpcResp.Code)])
		return nil, fmt.Errorf("FileInfo RPC request unsuccessful, code: %d, message: %s",
			rpcResp.Code, constant.StatusText[int(rpcResp.Code)])
	}

	// 绑定返回的数据
	resp.Code = int(rpcResp.Code)
	resp.Msg = constant.StatusText[resp.Code]
	resp.File = types.FileObject{
		FileName: rpcResp.File.FileName,
		FileMD5:  rpcResp.File.FileMd5,
		FileType: rpcResp.File.FileType,
		FileSize: rpcResp.File.FileSize,
	}
	l.Logger.Infof("File info request successful, received obj : %+v", rpcResp)
	return
}

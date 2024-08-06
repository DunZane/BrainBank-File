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

type ListFilesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListFilesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListFilesLogic {
	return &ListFilesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListFilesLogic) ListFiles(req *types.ListFilesRequest) (resp *types.ListFilesResponse, err error) {
	// 默认请求成功
	resp = &types.ListFilesResponse{}
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

	// 创建列表文件rpc请求体
	rpcReq := &fileops.ListFilesRequest{
		UserId: userId,
		Limit:  int32(req.Limit),
		Offset: int32(req.Offset),
	}

	// 调用远程函数
	l.Logger.Infof("Preparing to call ListFile RPC with request body: %+v", rpcReq)
	rpcResp, err := l.svcCtx.FileOps.ListFiles(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Error calling UploadFile RPC: %v", err)
		return nil, fmt.Errorf("file.ops.rpc return err:%w", err)
	}
	l.Logger.Infof("Received response from ListFile RPC: %+v", rpcResp)

	// 检查请求是否满足预期
	if rpcResp.Code != constant.StatusOK {
		l.Logger.Errorf("ListFile RPC request unsuccessful, code: %d, message: %s",
			rpcResp.Code, constant.StatusText[int(rpcResp.Code)])
		return nil, fmt.Errorf("ListFile RPC request unsuccessful, code: %d, message: %s",
			rpcResp.Code, constant.StatusText[int(rpcResp.Code)])
	}

	// 绑定返回数据
	resp.Code = int(rpcResp.Code)
	resp.Msg = constant.StatusText[resp.Code]
	resp.Total = int(rpcResp.Total)
	files := make([]*types.FileObject, 0, len(resp.Files))
	for _, file := range rpcResp.Files {
		files = append(files, &types.FileObject{
			FileName: file.FileName,
			FileMD5:  file.FileMd5,
			FileType: file.FileType,
			FileSize: file.FileSize,
		})
	}
	resp.Files = files
	return
}

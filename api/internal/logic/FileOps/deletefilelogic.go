package FileOps

import (
	"context"
	"errors"
	"fmt"
	"github.com/dunzane/brainbank-file/api/internal/constant"
	"github.com/dunzane/brainbank-file/api/internal/svc"
	"github.com/dunzane/brainbank-file/api/internal/types"
	"github.com/dunzane/brainbank-file/rpc/fileOps/fileops"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteFileLogic {
	return &DeleteFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteFileLogic) DeleteFile(req *types.DeleteFileRequest) (resp *types.Base, err error) {
	// 默认请求成功
	resp = &types.Base{}
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
	rpcReq := &fileops.DeleteFileRequest{
		FileId: req.FileId,
		UserId: userId,
	}

	// 调用远程函数
	l.Logger.Infof("Preparing to call DeleteFile RPC with request body: %+v", rpcReq)
	rpcResp, err := l.svcCtx.FileOps.DeleteFile(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Error calling DeleteFile RPC: %v", err)
		return nil, fmt.Errorf("file.ops.rpc return err:%w", err)
	}
	l.Logger.Infof("Received response from DeleteFile RPC: %+v", rpcResp)

	// 检查请求是否满足预期
	if rpcResp.Code != constant.StatusOK {
		l.Logger.Errorf("DeleteFile RPC request unsuccessful, code: %d, message: %s",
			rpcResp.Code, constant.StatusText[int(rpcResp.Code)])
		return nil, fmt.Errorf("DeleteFile RPC request unsuccessful, code: %d, message: %s",
			rpcResp.Code, constant.StatusText[int(rpcResp.Code)])
	}

	// 绑定返回数据
	resp.Code = int(rpcResp.Code)
	resp.Msg = constant.StatusText[resp.Code]
	return
}

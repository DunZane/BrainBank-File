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

type UpdateFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFileLogic {
	return &UpdateFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateFileLogic) UpdateFile(req *types.UpdateFileRequest) (resp *types.Base, err error) {
	// 默认请求成功
	resp = &types.Base{}
	resp.Code = constant.StatusOK
	resp.Msg = constant.StatusText[resp.Code]

	// 解析JWT，获取关键数据
	userId, ok := l.ctx.Value("userId").(int64)
	if !ok {
		err = errors.New("can't parse user info from JWT token")
		l.Logger.Error(err)
		return nil, err
	}
	l.Logger.Infof("Parsed userID from JWT token: %d", userId)

	// 创建上传文件rpc请求体
	rpcReq := &fileops.UpdateFileRequest{
		UserId:        userId,
		Status:        req.Status,
		StoreProvider: req.StorageProvider,
		FileId:        req.FileId,
	}

	// 调用远程函数
	l.Logger.Debugf("Preparing to call UpdateFile RPC with request body: %+v", rpcReq)
	rpcResp, err := l.svcCtx.FileOps.UpdateFile(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Error calling UpdateFile RPC: %v", err)
		return nil, fmt.Errorf("file.ops.rpc return err:%w", err)
	}
	l.Logger.Debugf("Received response from UpdateFile RPC: %+v", rpcResp)

	// 检查请求是否满足预期
	if rpcResp.Code != constant.StatusOK {
		l.Logger.Errorf("UpdateFile RPC request unsuccessful, code: %d, message: %s",
			rpcResp.Code, constant.StatusText[int(rpcResp.Code)])
		return nil, fmt.Errorf("UpdateFile RPC request unsuccessful, code: %d, message: %s",
			rpcResp.Code, constant.StatusText[int(rpcResp.Code)])
	}

	// 绑定返回数据
	resp.Code = int(rpcResp.Code)
	resp.Msg = constant.StatusText[int(rpcResp.Code)]

	l.Logger.Infof("UpdateFile operation successful. Response code: %d, message: %s", resp.Code, resp.Msg)
	return
}

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

type PreviewFileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPreviewFileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreviewFileLogic {
	return &PreviewFileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PreviewFileLogic) PreviewFile(req *types.PreviewFileRequest) (resp *types.PreviewFileResponse, err error) {
	// 默认请求成功
	resp = &types.PreviewFileResponse{}
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

	// 创建预览文件rpc请求
	rpcReq := &fileinfo.PreviewFileRequest{
		FileId: req.FileId,
		UserId: userId,
	}

	// 调用远程函数
	l.Logger.Debugf("Preparing to call PreviewFile RPC with request body: %+v", rpcReq)
	rpcResp, err := l.svcCtx.FileInfo.PreviewFile(l.ctx, rpcReq)
	if err != nil {
		l.Logger.Errorf("Error calling PreviewFile RPC: %v", err)
		return nil, fmt.Errorf("file.info.rpc return err:%w", err)
	}
	l.Logger.Debugf("Received response from PreviewFile RPC: %+v", rpcResp)

	// 检查请求是否满足预期
	if rpcResp.Code != constant.StatusOK {
		l.Logger.Errorf("PreviewFile RPC request unsuccessful, code: %d, message: %s",
			rpcResp.Code, constant.StatusText[int(rpcResp.Code)])
		return nil, fmt.Errorf("PreviewFile RPC request unsuccessful, code: %d, message: %s",
			rpcResp.Code, constant.StatusText[int(rpcResp.Code)])
	}

	// 返回绑定的数据
	resp.Code = int(rpcResp.Code)
	resp.Msg = constant.StatusText[resp.Code]
	resp.PresignedURL = rpcResp.PreSignedUrl
	resp.ExpiresAt = rpcResp.ExpiresAt
	return
}

// Code generated by goctl. DO NOT EDIT.
// Source: fileOps.proto

package fileops

import (
	"context"

	"github.com/dunzane/brainbank-file/rpc/fileOps/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	DeleteFileRequest  = pb.DeleteFileRequest
	DeleteFileResponse = pb.DeleteFileResponse
	FileObject         = pb.FileObject
	ListFilesRequest   = pb.ListFilesRequest
	ListFilesResponse  = pb.ListFilesResponse
	UpdateFileRequest  = pb.UpdateFileRequest
	UpdateFileResponse = pb.UpdateFileResponse
	UploadFileRequest  = pb.UploadFileRequest
	UploadFileResponse = pb.UploadFileResponse

	FileOps interface {
		UploadFile(ctx context.Context, in *UploadFileRequest, opts ...grpc.CallOption) (*UploadFileResponse, error)
		UpdateFile(ctx context.Context, in *UpdateFileRequest, opts ...grpc.CallOption) (*UpdateFileResponse, error)
		ListFiles(ctx context.Context, in *ListFilesRequest, opts ...grpc.CallOption) (*ListFilesResponse, error)
		DeleteFile(ctx context.Context, in *DeleteFileRequest, opts ...grpc.CallOption) (*DeleteFileResponse, error)
	}

	defaultFileOps struct {
		cli zrpc.Client
	}
)

func NewFileOps(cli zrpc.Client) FileOps {
	return &defaultFileOps{
		cli: cli,
	}
}

func (m *defaultFileOps) UploadFile(ctx context.Context, in *UploadFileRequest, opts ...grpc.CallOption) (*UploadFileResponse, error) {
	client := pb.NewFileOpsClient(m.cli.Conn())
	return client.UploadFile(ctx, in, opts...)
}

func (m *defaultFileOps) UpdateFile(ctx context.Context, in *UpdateFileRequest, opts ...grpc.CallOption) (*UpdateFileResponse, error) {
	client := pb.NewFileOpsClient(m.cli.Conn())
	return client.UpdateFile(ctx, in, opts...)
}

func (m *defaultFileOps) ListFiles(ctx context.Context, in *ListFilesRequest, opts ...grpc.CallOption) (*ListFilesResponse, error) {
	client := pb.NewFileOpsClient(m.cli.Conn())
	return client.ListFiles(ctx, in, opts...)
}

func (m *defaultFileOps) DeleteFile(ctx context.Context, in *DeleteFileRequest, opts ...grpc.CallOption) (*DeleteFileResponse, error) {
	client := pb.NewFileOpsClient(m.cli.Conn())
	return client.DeleteFile(ctx, in, opts...)
}

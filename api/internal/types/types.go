// Code generated by goctl. DO NOT EDIT.
package types

type UploadFileRequst struct {
	File FileObject `json:"file"`
}

type UploadFileResponse struct {
	Base
	PresignedURL string `json:"presigned_url"`
	ExpiresAt    string `json:"expires_at"`
	FileId       string `json:"file_id"`
}

type UpdateFileRequest struct {
	Status          string `json:"status"`
	StorageProvider string `json:"store_provider"`
	FileId          string `json:"file_id"`
}

type ListFilesRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ListFilesResponse struct {
	Base
	Files []*FileObject `json:"files"`
	Total int           `json:"total""`
}

type DeleteFileRequest struct {
	FileId string `json:"file_id"`
}

type PreviewFileRequest struct {
	FileId string `json:"file_id"`
}

type PreviewFileResponse struct {
	Base
	PresignedURL string `json:"presigned_url"`
	ExpiresAt    string `json:"expires_at"`
}

type GetFileInfoRequest struct {
	FileId string `json:"file_id"`
}

type GetFileInfoResponse struct {
	Base
	File FileObject `json:"file"`
}

type Base struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type FileObject struct {
	FileName string `json:"file_name" validate:"required"`
	FileMD5  string `json:"file_md5" validate:"required"`
	FileType string `json:"file_type" validate:"required"`
	FileSize int64  `json:"file_size" validate:"required"`
}
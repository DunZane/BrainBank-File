package constant

// 通用状态码
const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusNotFound            = 404
	StatusInternalServerError = 500
)

// 业务状态码
const (
	ErrorGenPresignedURL = 30001
	FileNotFound         = 30002
)

var StatusText = map[int]string{
	StatusOK:                  "OK",
	StatusBadRequest:          "Bad Request",
	StatusUnauthorized:        "Unauthorized",
	StatusNotFound:            "Not Found",
	StatusInternalServerError: "Internal Server Error",

	ErrorGenPresignedURL: "Error when try to gen a presigned url",
	FileNotFound:         "The user does not have this file object",
}

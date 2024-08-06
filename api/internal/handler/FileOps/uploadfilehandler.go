package FileOps

import (
	"net/http"

	"github.com/dunzane/brainbank-file/api/internal/logic/FileOps"
	"github.com/dunzane/brainbank-file/api/internal/svc"
	"github.com/dunzane/brainbank-file/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadFileRequst
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := FileOps.NewUploadFileLogic(r.Context(), svcCtx)
		resp, err := l.UploadFile(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

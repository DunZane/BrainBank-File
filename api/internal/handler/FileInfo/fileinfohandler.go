package FileInfo

import (
	"net/http"

	"github.com/dunzane/brainbank-file/api/internal/logic/FileInfo"
	"github.com/dunzane/brainbank-file/api/internal/svc"
	"github.com/dunzane/brainbank-file/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetFileInfoRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := FileInfo.NewFileInfoLogic(r.Context(), svcCtx)
		resp, err := l.FileInfo(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

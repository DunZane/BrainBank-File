package FileInfo

import (
	"net/http"

	"github.com/dunzane/brainbank-file/api/internal/logic/FileInfo"
	"github.com/dunzane/brainbank-file/api/internal/svc"
	"github.com/dunzane/brainbank-file/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func PreviewFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PreviewFileRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := FileInfo.NewPreviewFileLogic(r.Context(), svcCtx)
		resp, err := l.PreviewFile(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

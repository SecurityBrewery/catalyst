package catalyst

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/SecurityBrewery/catalyst/ui"
)

func static(ctx *gin.Context) {
	fsys, _ := fs.Sub(ui.UI, "dist")

	upath := strings.TrimPrefix(ctx.Request.URL.Path, "/")

	if _, err := fs.Stat(fsys, upath); err != nil {
		ctx.Request.URL.Path = "/"
		ctx.Request.URL.RawPath = "/"
	}

	http.FileServer(http.FS(fsys)).ServeHTTP(ctx.Writer, ctx.Request)
}

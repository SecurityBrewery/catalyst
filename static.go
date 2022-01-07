package catalyst

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/SecurityBrewery/catalyst/ui"
)

func static(w http.ResponseWriter, r *http.Request) {
	fsys, _ := fs.Sub(ui.UI, "dist")

	upath := strings.TrimPrefix(r.URL.Path, "/")

	if _, err := fs.Stat(fsys, upath); err != nil {
		r.URL.Path = "/"
		r.URL.RawPath = "/"
	}

	http.FileServer(http.FS(fsys)).ServeHTTP(w, r)
}

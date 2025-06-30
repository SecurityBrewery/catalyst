package router

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/SecurityBrewery/catalyst/ui"
)

func staticFiles(w http.ResponseWriter, r *http.Request) {
	if devServer := os.Getenv("UI_DEVSERVER"); devServer != "" {
		u, _ := url.Parse(devServer)

		r.Host = r.URL.Host

		httputil.NewSingleHostReverseProxy(u).ServeHTTP(w, r)

		return
	}

	vueStatic(w, r)
}

func vueStatic(w http.ResponseWriter, r *http.Request) {
	handler := http.FileServer(http.FS(ui.UI()))

	if strings.HasPrefix(r.URL.Path, "/ui/assets/") {
		handler = http.StripPrefix("/ui", handler)
	} else {
		r.URL.Path = "/"
	}

	handler.ServeHTTP(w, r)
}

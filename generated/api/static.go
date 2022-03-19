package api

import (
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func VueStatic(fsys fs.FS) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler := http.FileServer(http.FS(fsys))

		if strings.HasPrefix(r.URL.Path, "/static/") {
			handler = http.StripPrefix("/static/", handler)
		} else {
			r.URL.Path = "/"
		}

		handler.ServeHTTP(w, r)
	}
}

func Static(fsys fs.FS) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		http.FileServer(http.FS(fsys)).ServeHTTP(w, r)
	}
}

func Proxy(dest string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		u, _ := url.Parse(dest)
		proxy := httputil.NewSingleHostReverseProxy(u)

		r.Host = r.URL.Host

		proxy.ServeHTTP(w, r)
	}
}

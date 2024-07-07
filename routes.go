package main

import (
	"embed"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

//go:embed ui/dist/*
var ui embed.FS

func dev() bool {
	return strings.HasPrefix(os.Args[0], os.TempDir())
}

func addRoutes() func(*core.ServeEvent) error {
	return func(e *core.ServeEvent) error {
		e.Router.GET("/", func(c echo.Context) error {
			return c.Redirect(http.StatusFound, "/ui/")
		})
		e.Router.GET("/ui/*", staticFiles())
		e.Router.GET("/api/config", func(c echo.Context) error {
			flags, err := flags(e.App)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, map[string]any{
				"flags": flags,
			})
		})

		return nil
	}
}

func staticFiles() func(echo.Context) error {
	return func(c echo.Context) error {
		if dev() {
			u, _ := url.Parse("http://localhost:3000/")
			proxy := httputil.NewSingleHostReverseProxy(u)

			c.Request().Host = c.Request().URL.Host

			proxy.ServeHTTP(c.Response(), c.Request())

			return nil
		}

		fsys, _ := fs.Sub(ui, "ui/dist")

		return apis.StaticDirectoryHandler(fsys, true)(c)
	}
}
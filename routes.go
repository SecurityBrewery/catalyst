package main

import (
	"embed"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"github.com/SecurityBrewery/catalyst/ff"
)

//go:embed ui/dist/*
var ui embed.FS

func addRoutes() func(*core.ServeEvent) error {
	return func(e *core.ServeEvent) error {
		e.Router.GET("/", func(c echo.Context) error {
			return c.Redirect(http.StatusFound, "/ui/")
		})
		e.Router.GET("/ui/*", staticFiles())
		e.Router.GET("/api/flags", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]any{
				"flags": ff.Flags(),
			})
		})

		return nil
	}
}

func staticFiles() func(echo.Context) error {
	return func(c echo.Context) error {
		if ff.HasDevFlag() {
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

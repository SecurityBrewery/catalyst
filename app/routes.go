package app

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"github.com/SecurityBrewery/catalyst/ui"
)

func addRoutes() func(*core.ServeEvent) error {
	return func(e *core.ServeEvent) error {
		e.Router.GET("/", func(c echo.Context) error {
			return c.Redirect(http.StatusFound, "/ui/")
		})
		e.Router.GET("/ui/*", staticFiles())

		e.Router.GET("/api/config", func(c echo.Context) error {
			flags, err := Flags(e.App)
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

			c.Request().Host = c.Request().URL.Host

			httputil.NewSingleHostReverseProxy(u).ServeHTTP(c.Response(), c.Request())

			return nil
		}

		return apis.StaticDirectoryHandler(ui.UI(), true)(c)
	}
}

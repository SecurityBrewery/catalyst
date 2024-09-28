package app

import (
	"github.com/SecurityBrewery/catalyst-analysis/cmd/server/service"
	"github.com/SecurityBrewery/catalyst-analysis/generated/api"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"github.com/SecurityBrewery/catalyst/analysis"
	"github.com/SecurityBrewery/catalyst/ui"
)

func addRoutes(engine *analysis.Engine) func(*core.ServeEvent) error {
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

		e.Router.GET("/api/analysis/*", func(c echo.Context) error {
			if err := engine.SetDao(e.App.Dao()); err != nil {
				return err
			}

			apiServer, err := api.NewServer(service.New(engine.Engine()))
			if err != nil {
				return err
			}

			return echo.WrapHandler(http.StripPrefix("/api/analysis", apiServer))(c)
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

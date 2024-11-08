package app

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"github.com/SecurityBrewery/catalyst/ui"
)

func setupServer(appURL string, flags []string) func(e *core.ServeEvent) error {
	return func(e *core.ServeEvent) error {
		if err := SetFlags(e.App, flags); err != nil {
			return err
		}

		if HasFlag(e.App, "demo") {
			bindDemoHooks(e.App)
		}

		if appURL != "" {
			s := e.App.Settings()
			s.Meta.AppUrl = appURL

			if err := e.App.Dao().SaveSettings(s); err != nil {
				return err
			}
		}

		e.Router.GET("/", func(c echo.Context) error {
			return c.Redirect(http.StatusFound, "/ui/")
		})
		e.Router.GET("/ui/*", staticFiles())
		e.Router.GET("/health", func(c echo.Context) error {
			if _, err := Flags(e.App); err != nil {
				return err
			}

			return c.String(http.StatusOK, "OK")
		})

		e.Router.GET("/api/config", func(c echo.Context) error {
			flags, err := Flags(e.App)
			if err != nil {
				return err
			}

			return c.JSON(http.StatusOK, map[string]any{
				"flags": flags,
			})
		})

		return e.App.RefreshSettings()
	}
}

func bindDemoHooks(app core.App) {
	app.OnRecordBeforeCreateRequest("files", "reactions").Add(func(e *core.RecordCreateEvent) error {
		return fmt.Errorf("cannot create %s in demo mode", e.Record.Collection().Name)
	})
	app.OnRecordBeforeUpdateRequest("files", "reactions").Add(func(e *core.RecordUpdateEvent) error {
		return fmt.Errorf("cannot update %s in demo mode", e.Record.Collection().Name)
	})
	app.OnRecordBeforeDeleteRequest("files", "reactions").Add(func(e *core.RecordDeleteEvent) error {
		return fmt.Errorf("cannot delete %s in demo mode", e.Record.Collection().Name)
	})
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

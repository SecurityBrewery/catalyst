package testing

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

type RecordingServer struct {
	server *echo.Echo

	Entries []string
}

func NewRecordingServer() *RecordingServer {
	e := echo.New()

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"status": "ok",
		})
	})
	e.Any("/*", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"test": true,
		})
	})

	return &RecordingServer{
		server: e,
	}
}

func (s *RecordingServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Entries = append(s.Entries, r.URL.Path)

	s.server.ServeHTTP(w, r)
}

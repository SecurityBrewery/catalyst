package testing

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type RecordingServer struct {
	server chi.Router

	Entries []string
}

func NewRecordingServer() *RecordingServer {
	e := chi.NewRouter()

	e.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(map[string]any{
			"status": "ok",
		})
		if err != nil {
			http.Error(w, "failed to marshal response", http.StatusInternalServerError)

			return
		}

		_, err = w.Write(b)
	})
	e.Handle("/*", func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(map[string]any{
			"test": true,
		})
		if err != nil {
			http.Error(w, "failed to marshal response", http.StatusInternalServerError)

			return
		}

		_, err = w.Write(b)
	})

	return &RecordingServer{
		server: e,
	}
}

func (s *RecordingServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Entries = append(s.Entries, r.URL.Path)

	s.server.ServeHTTP(w, r)
}

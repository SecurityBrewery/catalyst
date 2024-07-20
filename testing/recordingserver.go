package testing

import "net/http"

type RecordingServer struct {
	Entries []string
}

func NewRecordingServer() *RecordingServer {
	return &RecordingServer{}
}

func (s *RecordingServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Entries = append(s.Entries, r.URL.Path)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"test":true}`)) //nolint:errcheck
}

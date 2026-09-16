package handler

import (
	"net/http"

	"sandbox/pkg/httpx"
)

func (s *Server) registerExportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/export/snapshot", s.exportSnapshot)
}

func (s *Server) exportSnapshot(w http.ResponseWriter, r *http.Request) {
	snapshot, err := s.svc.ExportSnapshot()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, snapshot)
}

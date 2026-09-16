package handler

import (
	"net/http"

	"sandbox/internal/model"
	"sandbox/pkg/httpx"
)

func (s *Server) registerArtifactRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/artifacts", s.createArtifact)
	mux.HandleFunc("GET /api/artifacts", s.listArtifacts)
	mux.HandleFunc("GET /api/artifacts/{id}", s.getArtifact)
	mux.HandleFunc("PUT /api/artifacts/{id}", s.updateArtifact)
	mux.HandleFunc("DELETE /api/artifacts/{id}", s.deleteArtifact)
}

type createArtifactRequest struct {
	TaskID    string `json:"task_id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	SizeBytes int64  `json:"size_bytes"`
	Status    string `json:"status"`
}

func (s *Server) createArtifact(w http.ResponseWriter, r *http.Request) {
	var req createArtifactRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.CreateArtifact(model.Artifact{
		TaskID:    req.TaskID,
		Name:      req.Name,
		Path:      req.Path,
		SizeBytes: req.SizeBytes,
		Status:    req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, a)
}

func (s *Server) listArtifacts(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ArtifactFilter{
		TaskID: r.URL.Query().Get("task_id"),
		Status: r.URL.Query().Get("status"),
		Name:   r.URL.Query().Get("name"),
	}
	items, total, err := s.svc.ListArtifacts(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getArtifact(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := s.svc.GetArtifact(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

type updateArtifactRequest struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	SizeBytes int64  `json:"size_bytes"`
	Status    string `json:"status"`
}

func (s *Server) updateArtifact(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateArtifactRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	a, err := s.svc.UpdateArtifact(id, model.Artifact{
		Name:      req.Name,
		Path:      req.Path,
		SizeBytes: req.SizeBytes,
		Status:    req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, a)
}

func (s *Server) deleteArtifact(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteArtifact(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

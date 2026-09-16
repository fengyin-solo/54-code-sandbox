package handler

import (
	"net/http"

	"sandbox/internal/model"
	"sandbox/pkg/httpx"
)

func (s *Server) registerSandboxRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/sandboxes", s.createSandbox)
	mux.HandleFunc("GET /api/sandboxes", s.listSandboxes)
	mux.HandleFunc("GET /api/sandboxes/{id}", s.getSandbox)
	mux.HandleFunc("PUT /api/sandboxes/{id}", s.updateSandbox)
	mux.HandleFunc("DELETE /api/sandboxes/{id}", s.deleteSandbox)
}

type createSandboxRequest struct {
	Name     string `json:"name"`
	Language string `json:"language"`
	Image    string `json:"image"`
	Status   string `json:"status"`
}

func (s *Server) createSandbox(w http.ResponseWriter, r *http.Request) {
	var req createSandboxRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sb, err := s.svc.CreateSandbox(model.Sandbox{Name: req.Name, Language: req.Language, Image: req.Image, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, sb)
}

func (s *Server) listSandboxes(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.SandboxFilter{
		Language: r.URL.Query().Get("language"),
		Status:   r.URL.Query().Get("status"),
		Keyword:  r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListSandboxes(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSandbox(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sb, err := s.svc.GetSandbox(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sb)
}

type updateSandboxRequest struct {
	Name     string `json:"name"`
	Language string `json:"language"`
	Image    string `json:"image"`
	Status   string `json:"status"`
}

func (s *Server) updateSandbox(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateSandboxRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sb, err := s.svc.UpdateSandbox(id, model.Sandbox{Name: req.Name, Language: req.Language, Image: req.Image, Status: req.Status})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sb)
}

func (s *Server) deleteSandbox(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteSandbox(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

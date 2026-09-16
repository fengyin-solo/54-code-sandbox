package handler

import (
	"net/http"

	"sandbox/internal/model"
	"sandbox/pkg/httpx"
)

func (s *Server) registerExecutionTaskRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/tasks", s.createExecutionTask)
	mux.HandleFunc("GET /api/tasks", s.listExecutionTasks)
	mux.HandleFunc("GET /api/tasks/{id}", s.getExecutionTask)
	mux.HandleFunc("PUT /api/tasks/{id}", s.updateExecutionTask)
	mux.HandleFunc("DELETE /api/tasks/{id}", s.deleteExecutionTask)
	mux.HandleFunc("POST /api/tasks/{id}/run", s.runTask)
	mux.HandleFunc("POST /api/tasks/{id}/complete", s.completeTask)
	mux.HandleFunc("POST /api/tasks/{id}/fail", s.failTask)
	mux.HandleFunc("POST /api/tasks/{id}/timeout", s.timeoutTask)
	mux.HandleFunc("POST /api/tasks/{id}/kill", s.killTask)
}

type createExecutionTaskRequest struct {
	SandboxID string `json:"sandbox_id"`
	Code      string `json:"code"`
	Language  string `json:"language"`
	Stdin     string `json:"stdin"`
}

func (s *Server) createExecutionTask(w http.ResponseWriter, r *http.Request) {
	var req createExecutionTaskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.CreateExecutionTask(model.ExecutionTask{
		SandboxID: req.SandboxID,
		Code:      req.Code,
		Language:  req.Language,
		Stdin:     req.Stdin,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, t)
}

func (s *Server) listExecutionTasks(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ExecutionTaskFilter{
		SandboxID: r.URL.Query().Get("sandbox_id"),
		Language:  r.URL.Query().Get("language"),
		Status:    r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListExecutionTasks(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getExecutionTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := s.svc.GetExecutionTask(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

type updateExecutionTaskRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
	Stdin    string `json:"stdin"`
	Status   string `json:"status"`
}

func (s *Server) updateExecutionTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateExecutionTaskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.UpdateExecutionTask(id, model.ExecutionTask{
		Code:     req.Code,
		Language: req.Language,
		Stdin:    req.Stdin,
		Status:   req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) deleteExecutionTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteExecutionTask(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) runTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := s.svc.RunTask(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

type completeTaskRequest struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exit_code"`
}

func (s *Server) completeTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req completeTaskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.CompleteTask(id, req.Stdout, req.Stderr, req.ExitCode)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

type failTaskRequest struct {
	Stderr string `json:"stderr"`
}

func (s *Server) failTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req failTaskRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.FailTask(id, req.Stderr)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) timeoutTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := s.svc.TimeoutTask(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) killTask(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := s.svc.KillTask(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

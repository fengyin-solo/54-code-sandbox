package handler

import (
	"net/http"

	"sandbox/internal/model"
	"sandbox/pkg/httpx"
)

func (s *Server) registerExecutionLogRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/logs", s.createExecutionLog)
	mux.HandleFunc("GET /api/logs", s.listExecutionLogs)
	mux.HandleFunc("GET /api/logs/{id}", s.getExecutionLog)
	mux.HandleFunc("DELETE /api/logs/{id}", s.deleteExecutionLog)
}

type createExecutionLogRequest struct {
	TaskID  string `json:"task_id"`
	Level   string `json:"level"`
	Message string `json:"message"`
}

func (s *Server) createExecutionLog(w http.ResponseWriter, r *http.Request) {
	var req createExecutionLogRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	l, err := s.svc.CreateExecutionLog(model.ExecutionLog{
		TaskID:  req.TaskID,
		Level:   req.Level,
		Message: req.Message,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, l)
}

func (s *Server) listExecutionLogs(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ExecutionLogFilter{
		TaskID: r.URL.Query().Get("task_id"),
		Level:  r.URL.Query().Get("level"),
	}
	items, total, err := s.svc.ListExecutionLogs(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getExecutionLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	l, err := s.svc.GetExecutionLog(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, l)
}

func (s *Server) deleteExecutionLog(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteExecutionLog(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

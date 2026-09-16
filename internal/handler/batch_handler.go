package handler

import (
	"net/http"

	"sandbox/internal/model"
	"sandbox/pkg/httpx"
)

func (s *Server) registerBatchRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/tasks/batch", s.batchCreateExecutionTasks)
}

func (s *Server) batchCreateExecutionTasks(w http.ResponseWriter, r *http.Request) {
	var reqs []model.ExecutionTask
	if err := httpx.Decode(r, &reqs); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	results, err := s.svc.BatchCreateExecutionTasks(reqs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, results)
}

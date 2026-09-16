package handler

import (
	"net/http"

	"sandbox/internal/model"
	"sandbox/pkg/httpx"
)

func (s *Server) registerSubmissionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/submissions", s.createSubmission)
	mux.HandleFunc("GET /api/submissions", s.listSubmissions)
	mux.HandleFunc("GET /api/submissions/{id}", s.getSubmission)
	mux.HandleFunc("PUT /api/submissions/{id}", s.updateSubmission)
	mux.HandleFunc("DELETE /api/submissions/{id}", s.deleteSubmission)
}

type createSubmissionRequest struct {
	TaskID     string `json:"task_id"`
	Submitter  string `json:"submitter"`
	SourceCode string `json:"source_code"`
	Status     string `json:"status"`
}

func (s *Server) createSubmission(w http.ResponseWriter, r *http.Request) {
	var req createSubmissionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sub, err := s.svc.CreateSubmission(model.Submission{
		TaskID:     req.TaskID,
		Submitter:  req.Submitter,
		SourceCode: req.SourceCode,
		Status:     req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, sub)
}

func (s *Server) listSubmissions(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.SubmissionFilter{
		TaskID:    r.URL.Query().Get("task_id"),
		Submitter: r.URL.Query().Get("submitter"),
		Status:    r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListSubmissions(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSubmission(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sub, err := s.svc.GetSubmission(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sub)
}

type updateSubmissionRequest struct {
	Submitter  string `json:"submitter"`
	SourceCode string `json:"source_code"`
	Status     string `json:"status"`
}

func (s *Server) updateSubmission(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateSubmissionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sub, err := s.svc.UpdateSubmission(id, model.Submission{
		Submitter:  req.Submitter,
		SourceCode: req.SourceCode,
		Status:     req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sub)
}

func (s *Server) deleteSubmission(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteSubmission(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

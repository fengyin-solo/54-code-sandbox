package handler

import (
	"net/http"

	"sandbox/internal/model"
	"sandbox/pkg/httpx"
)

func (s *Server) registerScheduleRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/schedules", s.createSchedule)
	mux.HandleFunc("GET /api/schedules", s.listSchedules)
	mux.HandleFunc("GET /api/schedules/{id}", s.getSchedule)
	mux.HandleFunc("PUT /api/schedules/{id}", s.updateSchedule)
	mux.HandleFunc("DELETE /api/schedules/{id}", s.deleteSchedule)
	mux.HandleFunc("POST /api/schedules/{id}/run", s.runSchedule)
	mux.HandleFunc("POST /api/schedules/{id}/pause", s.pauseSchedule)
	mux.HandleFunc("POST /api/schedules/{id}/complete", s.completeSchedule)
}

type createScheduleRequest struct {
	SandboxID string `json:"sandbox_id"`
	Code      string `json:"code"`
	CronExpr  string `json:"cron_expr"`
}

func (s *Server) createSchedule(w http.ResponseWriter, r *http.Request) {
	var req createScheduleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sch, err := s.svc.CreateSchedule(model.Schedule{
		SandboxID: req.SandboxID,
		Code:      req.Code,
		CronExpr:  req.CronExpr,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, sch)
}

func (s *Server) listSchedules(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ScheduleFilter{
		SandboxID: r.URL.Query().Get("sandbox_id"),
		Status:    r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListSchedules(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sch, err := s.svc.GetSchedule(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sch)
}

type updateScheduleRequest struct {
	Code     string `json:"code"`
	CronExpr string `json:"cron_expr"`
	Status   string `json:"status"`
}

func (s *Server) updateSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateScheduleRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sch, err := s.svc.UpdateSchedule(id, model.Schedule{
		Code:     req.Code,
		CronExpr: req.CronExpr,
		Status:   req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sch)
}

func (s *Server) deleteSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteSchedule(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) runSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sch, err := s.svc.RunSchedule(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sch)
}

func (s *Server) pauseSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sch, err := s.svc.PauseSchedule(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sch)
}

func (s *Server) completeSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sch, err := s.svc.CompleteSchedule(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sch)
}

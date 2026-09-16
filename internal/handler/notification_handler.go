package handler

import (
	"net/http"

	"sandbox/internal/model"
	"sandbox/pkg/httpx"
)

func (s *Server) registerNotificationRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/notifications", s.createNotification)
	mux.HandleFunc("GET /api/notifications", s.listNotifications)
	mux.HandleFunc("GET /api/notifications/{id}", s.getNotification)
	mux.HandleFunc("PUT /api/notifications/{id}", s.updateNotification)
	mux.HandleFunc("DELETE /api/notifications/{id}", s.deleteNotification)
	mux.HandleFunc("POST /api/notifications/{id}/send", s.sendNotification)
	mux.HandleFunc("POST /api/notifications/{id}/deliver", s.deliverNotification)
	mux.HandleFunc("POST /api/notifications/{id}/fail", s.failNotification)
}

type createNotificationRequest struct {
	WebhookID string `json:"webhook_id"`
	Event     string `json:"event"`
	Payload   string `json:"payload"`
}

func (s *Server) createNotification(w http.ResponseWriter, r *http.Request) {
	var req createNotificationRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	n, err := s.svc.CreateNotification(model.Notification{
		WebhookID: req.WebhookID,
		Event:     req.Event,
		Payload:   req.Payload,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, n)
}

func (s *Server) listNotifications(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.NotificationFilter{
		WebhookID: r.URL.Query().Get("webhook_id"),
		Event:     r.URL.Query().Get("event"),
		Status:    r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListNotifications(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getNotification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	n, err := s.svc.GetNotification(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, n)
}

type updateNotificationRequest struct {
	Status     string `json:"status"`
	RetryCount int    `json:"retry_count"`
	Payload    string `json:"payload"`
}

func (s *Server) updateNotification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateNotificationRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	n, err := s.svc.UpdateNotification(id, model.Notification{
		Status:     req.Status,
		RetryCount: req.RetryCount,
		Payload:    req.Payload,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, n)
}

func (s *Server) deleteNotification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteNotification(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) sendNotification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	n, err := s.svc.MarkNotificationSent(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, n)
}

func (s *Server) deliverNotification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	n, err := s.svc.MarkNotificationDelivered(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, n)
}

func (s *Server) failNotification(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	n, err := s.svc.MarkNotificationFailed(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, n)
}

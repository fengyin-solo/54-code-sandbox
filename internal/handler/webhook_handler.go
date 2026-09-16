package handler

import (
	"net/http"

	"sandbox/internal/model"
	"sandbox/pkg/httpx"
)

func (s *Server) registerWebhookRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/webhooks", s.createWebhook)
	mux.HandleFunc("GET /api/webhooks", s.listWebhooks)
	mux.HandleFunc("GET /api/webhooks/{id}", s.getWebhook)
	mux.HandleFunc("PUT /api/webhooks/{id}", s.updateWebhook)
	mux.HandleFunc("DELETE /api/webhooks/{id}", s.deleteWebhook)
}

type createWebhookRequest struct {
	Name   string   `json:"name"`
	URL    string   `json:"url"`
	Secret string   `json:"secret"`
	Events []string `json:"events"`
	Status string   `json:"status"`
}

func (s *Server) createWebhook(w http.ResponseWriter, r *http.Request) {
	var req createWebhookRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	wh, err := s.svc.CreateWebhook(model.Webhook{
		Name:   req.Name,
		URL:    req.URL,
		Secret: req.Secret,
		Events: req.Events,
		Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, wh)
}

func (s *Server) listWebhooks(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.WebhookFilter{
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListWebhooks(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getWebhook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	wh, err := s.svc.GetWebhook(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, wh)
}

type updateWebhookRequest struct {
	Name   string   `json:"name"`
	URL    string   `json:"url"`
	Secret string   `json:"secret"`
	Events []string `json:"events"`
	Status string   `json:"status"`
}

func (s *Server) updateWebhook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateWebhookRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	wh, err := s.svc.UpdateWebhook(id, model.Webhook{
		Name:   req.Name,
		URL:    req.URL,
		Secret: req.Secret,
		Events: req.Events,
		Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, wh)
}

func (s *Server) deleteWebhook(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteWebhook(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

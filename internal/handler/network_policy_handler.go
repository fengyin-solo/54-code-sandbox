package handler

import (
	"net/http"

	"sandbox/internal/model"
	"sandbox/pkg/httpx"
)

func (s *Server) registerNetworkPolicyRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/network-policies", s.createNetworkPolicy)
	mux.HandleFunc("GET /api/network-policies", s.listNetworkPolicies)
	mux.HandleFunc("GET /api/network-policies/{id}", s.getNetworkPolicy)
	mux.HandleFunc("PUT /api/network-policies/{id}", s.updateNetworkPolicy)
	mux.HandleFunc("DELETE /api/network-policies/{id}", s.deleteNetworkPolicy)
}

type createNetworkPolicyRequest struct {
	SandboxID    string   `json:"sandbox_id"`
	AllowedHosts []string `json:"allowed_hosts"`
	BlockedHosts []string `json:"blocked_hosts"`
	Status       string   `json:"status"`
}

func (s *Server) createNetworkPolicy(w http.ResponseWriter, r *http.Request) {
	var req createNetworkPolicyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	n, err := s.svc.CreateNetworkPolicy(model.NetworkPolicy{
		SandboxID:    req.SandboxID,
		AllowedHosts: req.AllowedHosts,
		BlockedHosts: req.BlockedHosts,
		Status:       req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, n)
}

func (s *Server) listNetworkPolicies(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.NetworkPolicyFilter{
		SandboxID: r.URL.Query().Get("sandbox_id"),
		Status:    r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListNetworkPolicies(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getNetworkPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	n, err := s.svc.GetNetworkPolicy(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, n)
}

type updateNetworkPolicyRequest struct {
	AllowedHosts []string `json:"allowed_hosts"`
	BlockedHosts []string `json:"blocked_hosts"`
	Status       string   `json:"status"`
}

func (s *Server) updateNetworkPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateNetworkPolicyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	n, err := s.svc.UpdateNetworkPolicy(id, model.NetworkPolicy{
		AllowedHosts: req.AllowedHosts,
		BlockedHosts: req.BlockedHosts,
		Status:       req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, n)
}

func (s *Server) deleteNetworkPolicy(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteNetworkPolicy(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

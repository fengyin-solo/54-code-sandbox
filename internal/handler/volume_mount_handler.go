package handler

import (
	"net/http"

	"sandbox/internal/model"
	"sandbox/pkg/httpx"
)

func (s *Server) registerVolumeMountRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/volumes", s.createVolumeMount)
	mux.HandleFunc("GET /api/volumes", s.listVolumeMounts)
	mux.HandleFunc("GET /api/volumes/{id}", s.getVolumeMount)
	mux.HandleFunc("PUT /api/volumes/{id}", s.updateVolumeMount)
	mux.HandleFunc("DELETE /api/volumes/{id}", s.deleteVolumeMount)
}

type createVolumeMountRequest struct {
	SandboxID     string `json:"sandbox_id"`
	HostPath      string `json:"host_path"`
	ContainerPath string `json:"container_path"`
	ReadOnly      bool   `json:"read_only"`
	Status        string `json:"status"`
}

func (s *Server) createVolumeMount(w http.ResponseWriter, r *http.Request) {
	var req createVolumeMountRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	v, err := s.svc.CreateVolumeMount(model.VolumeMount{
		SandboxID:     req.SandboxID,
		HostPath:      req.HostPath,
		ContainerPath: req.ContainerPath,
		ReadOnly:      req.ReadOnly,
		Status:        req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, v)
}

func (s *Server) listVolumeMounts(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.VolumeMountFilter{
		SandboxID: r.URL.Query().Get("sandbox_id"),
		Status:    r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListVolumeMounts(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getVolumeMount(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	v, err := s.svc.GetVolumeMount(id)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, v)
}

type updateVolumeMountRequest struct {
	HostPath      string `json:"host_path"`
	ContainerPath string `json:"container_path"`
	ReadOnly      bool   `json:"read_only"`
	Status        string `json:"status"`
}

func (s *Server) updateVolumeMount(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req updateVolumeMountRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	v, err := s.svc.UpdateVolumeMount(id, model.VolumeMount{
		HostPath:      req.HostPath,
		ContainerPath: req.ContainerPath,
		ReadOnly:      req.ReadOnly,
		Status:        req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, v)
}

func (s *Server) deleteVolumeMount(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.svc.DeleteVolumeMount(id); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

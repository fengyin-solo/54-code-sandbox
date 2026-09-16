package model

import (
	"strings"
	"time"
)

const (
	VolumeActive   = "active"
	VolumeDisabled = "disabled"
)

type VolumeMount struct {
	ID             string    `json:"id"`
	SandboxID      string    `json:"sandbox_id"`
	HostPath       string    `json:"host_path"`
	ContainerPath  string    `json:"container_path"`
	ReadOnly       bool      `json:"read_only"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (v *VolumeMount) Validate() error {
	v.SandboxID = strings.TrimSpace(v.SandboxID)
	v.HostPath = strings.TrimSpace(v.HostPath)
	v.ContainerPath = strings.TrimSpace(v.ContainerPath)
	if v.SandboxID == "" {
		return NewValidationError("sandbox_id", "沙箱ID不能为空")
	}
	if v.HostPath == "" {
		return NewValidationError("host_path", "主机路径不能为空")
	}
	if v.ContainerPath == "" {
		return NewValidationError("container_path", "容器路径不能为空")
	}
	if v.Status == "" {
		v.Status = VolumeActive
	}
	if v.Status != VolumeActive && v.Status != VolumeDisabled {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

type VolumeMountFilter struct {
	SandboxID string
	Status    string
}

func (f VolumeMountFilter) Match(v *VolumeMount) bool {
	if f.SandboxID != "" && v.SandboxID != f.SandboxID {
		return false
	}
	if f.Status != "" && v.Status != f.Status {
		return false
	}
	return true
}

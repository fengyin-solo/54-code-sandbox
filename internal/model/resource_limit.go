package model

import (
	"strings"
	"time"
)

const (
	ResourceLimitActive   = "active"
	ResourceLimitDisabled = "disabled"
)

type ResourceLimit struct {
	ID           string    `json:"id"`
	SandboxID    string    `json:"sandbox_id"`
	CPUQuota     float64   `json:"cpu_quota"`
	MemoryMB     int       `json:"memory_mb"`
	DiskMB       int       `json:"disk_mb"`
	MaxProcesses int       `json:"max_processes"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (r *ResourceLimit) Validate() error {
	r.SandboxID = strings.TrimSpace(r.SandboxID)
	if r.SandboxID == "" {
		return NewValidationError("sandbox_id", "沙箱ID不能为空")
	}
	if r.CPUQuota <= 0 {
		return NewValidationError("cpu_quota", "CPU配额必须大于0")
	}
	if r.MemoryMB <= 0 {
		return NewValidationError("memory_mb", "内存限制必须大于0")
	}
	if r.DiskMB <= 0 {
		return NewValidationError("disk_mb", "磁盘限制必须大于0")
	}
	if r.MaxProcesses <= 0 {
		return NewValidationError("max_processes", "最大进程数必须大于0")
	}
	if r.Status == "" {
		r.Status = ResourceLimitActive
	}
	if r.Status != ResourceLimitActive && r.Status != ResourceLimitDisabled {
		return NewValidationError("status", "资源限制状态不合法")
	}
	return nil
}

type ResourceLimitFilter struct {
	SandboxID string
	Status    string
}

func (f ResourceLimitFilter) Match(r *ResourceLimit) bool {
	if f.SandboxID != "" && r.SandboxID != f.SandboxID {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	return true
}

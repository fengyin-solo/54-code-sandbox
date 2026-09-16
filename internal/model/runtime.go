package model

import (
	"strings"
	"time"
)

const (
	RuntimeActive   = "active"
	RuntimeDisabled = "disabled"
)

type Runtime struct {
	ID           string    `json:"id"`
	SandboxID    string    `json:"sandbox_id"`
	Name         string    `json:"name"`
	Version      string    `json:"version"`
	Command      string    `json:"command"`
	MemoryLimitMB int      `json:"memory_limit_mb"`
	CPULimit     float64   `json:"cpu_limit"`
	TimeoutMs    int       `json:"timeout_ms"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (r *Runtime) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.Version = strings.TrimSpace(r.Version)
	r.Command = strings.TrimSpace(r.Command)
	r.SandboxID = strings.TrimSpace(r.SandboxID)
	if r.SandboxID == "" {
		return NewValidationError("sandbox_id", "所属沙箱ID不能为空")
	}
	if r.Name == "" {
		return NewValidationError("name", "运行时名称不能为空")
	}
	if r.Version == "" {
		return NewValidationError("version", "版本不能为空")
	}
	if r.Command == "" {
		return NewValidationError("command", "启动命令不能为空")
	}
	if r.MemoryLimitMB <= 0 {
		return NewValidationError("memory_limit_mb", "内存限制必须大于0")
	}
	if r.CPULimit <= 0 {
		return NewValidationError("cpu_limit", "CPU限制必须大于0")
	}
	if r.TimeoutMs <= 0 {
		return NewValidationError("timeout_ms", "超时时间必须大于0")
	}
	if r.Status == "" {
		r.Status = RuntimeActive
	}
	if r.Status != RuntimeActive && r.Status != RuntimeDisabled {
		return NewValidationError("status", "运行时状态不合法")
	}
	return nil
}

type RuntimeFilter struct {
	SandboxID string
	Status    string
	Keyword   string
}

func (f RuntimeFilter) Match(r *Runtime) bool {
	if f.SandboxID != "" && r.SandboxID != f.SandboxID {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(r.Name), k) {
			return false
		}
	}
	return true
}

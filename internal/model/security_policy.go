package model

import (
	"strings"
	"time"
)

const (
	PolicyActive   = "active"
	PolicyDisabled = "disabled"
)

type SecurityPolicy struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	AllowNetwork     bool      `json:"allow_network"`
	AllowFileSystem  bool      `json:"allow_file_system"`
	MaxMemoryMB      int       `json:"max_memory_mb"`
	MaxCPUMs         int       `json:"max_cpu_ms"`
	MaxTimeoutMs     int       `json:"max_timeout_ms"`
	BlockedSyscalls  []string  `json:"blocked_syscalls"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (p *SecurityPolicy) Validate() error {
	p.Name = strings.TrimSpace(p.Name)
	if p.Name == "" {
		return NewValidationError("name", "策略名称不能为空")
	}
	if p.MaxMemoryMB <= 0 {
		return NewValidationError("max_memory_mb", "最大内存限制必须大于0")
	}
	if p.MaxCPUMs <= 0 {
		return NewValidationError("max_cpu_ms", "最大CPU时间限制必须大于0")
	}
	if p.MaxTimeoutMs <= 0 {
		return NewValidationError("max_timeout_ms", "最大超时限制必须大于0")
	}
	if p.Status == "" {
		p.Status = PolicyActive
	}
	if p.Status != PolicyActive && p.Status != PolicyDisabled {
		return NewValidationError("status", "策略状态不合法")
	}
	return nil
}

type SecurityPolicyFilter struct {
	Status  string
	Keyword string
}

func (f SecurityPolicyFilter) Match(p *SecurityPolicy) bool {
	if f.Status != "" && p.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(p.Name), k) {
			return false
		}
	}
	return true
}

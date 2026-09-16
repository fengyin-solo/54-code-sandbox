package model

import (
	"strings"
	"time"
)

const (
	EnvVarActive   = "active"
	EnvVarDisabled = "disabled"
)

type EnvironmentVariable struct {
	ID        string    `json:"id"`
	SandboxID string    `json:"sandbox_id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e *EnvironmentVariable) Validate() error {
	e.SandboxID = strings.TrimSpace(e.SandboxID)
	e.Key = strings.TrimSpace(e.Key)
	if e.SandboxID == "" {
		return NewValidationError("sandbox_id", "沙箱ID不能为空")
	}
	if e.Key == "" {
		return NewValidationError("key", "变量名不能为空")
	}
	if e.Status == "" {
		e.Status = EnvVarActive
	}
	if e.Status != EnvVarActive && e.Status != EnvVarDisabled {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

type EnvironmentVariableFilter struct {
	SandboxID string
	Key       string
	Status    string
}

func (f EnvironmentVariableFilter) Match(e *EnvironmentVariable) bool {
	if f.SandboxID != "" && e.SandboxID != f.SandboxID {
		return false
	}
	if f.Key != "" && e.Key != f.Key {
		return false
	}
	if f.Status != "" && e.Status != f.Status {
		return false
	}
	return true
}

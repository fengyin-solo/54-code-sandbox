package model

import (
	"strings"
	"time"
)

const (
	SandboxActive   = "active"
	SandboxDisabled = "disabled"
)

var sandboxLanguages = map[string]bool{"go": true, "python": true, "js": true, "java": true, "cpp": true}

type Sandbox struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Language  string    `json:"language"`
	Image     string    `json:"image"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Sandbox) Validate() error {
	s.Name = strings.TrimSpace(s.Name)
	s.Language = strings.TrimSpace(s.Language)
	s.Image = strings.TrimSpace(s.Image)
	if s.Name == "" {
		return NewValidationError("name", "沙箱名称不能为空")
	}
	if s.Language == "" {
		return NewValidationError("language", "语言不能为空")
	}
	if !sandboxLanguages[s.Language] {
		return NewValidationError("language", "不支持的编程语言")
	}
	if s.Image == "" {
		return NewValidationError("image", "镜像不能为空")
	}
	if s.Status == "" {
		s.Status = SandboxActive
	}
	if s.Status != SandboxActive && s.Status != SandboxDisabled {
		return NewValidationError("status", "沙箱状态不合法")
	}
	return nil
}

type SandboxFilter struct {
	Language string
	Status   string
	Keyword  string
}

func (f SandboxFilter) Match(s *Sandbox) bool {
	if f.Language != "" && s.Language != f.Language {
		return false
	}
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(s.Name), k) {
			return false
		}
	}
	return true
}

package model

import (
	"strings"
	"time"
)

const (
	TemplateActive   = "active"
	TemplateDisabled = "disabled"
)

type Template struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Language    string    `json:"language"`
	BaseImage   string    `json:"base_image"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (t *Template) Validate() error {
	t.Name = strings.TrimSpace(t.Name)
	t.Language = strings.TrimSpace(t.Language)
	t.BaseImage = strings.TrimSpace(t.BaseImage)
	if t.Name == "" {
		return NewValidationError("name", "模板名称不能为空")
	}
	if t.Language == "" {
		return NewValidationError("language", "语言不能为空")
	}
	if t.BaseImage == "" {
		return NewValidationError("base_image", "基础镜像不能为空")
	}
	if t.Status == "" {
		t.Status = TemplateActive
	}
	if t.Status != TemplateActive && t.Status != TemplateDisabled {
		return NewValidationError("status", "模板状态不合法")
	}
	return nil
}

type TemplateFilter struct {
	Language string
	Status   string
	Keyword  string
}

func (f TemplateFilter) Match(t *Template) bool {
	if f.Language != "" && t.Language != f.Language {
		return false
	}
	if f.Status != "" && t.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(t.Name), k) {
			return false
		}
	}
	return true
}

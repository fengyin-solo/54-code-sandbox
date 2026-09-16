package model

import (
	"strings"
	"time"
)

const (
	WebhookActive   = "active"
	WebhookDisabled = "disabled"
)

type Webhook struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	Secret      string    `json:"secret"`
	Events      []string  `json:"events"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (w *Webhook) Validate() error {
	w.Name = strings.TrimSpace(w.Name)
	w.URL = strings.TrimSpace(w.URL)
	if w.Name == "" {
		return NewValidationError("name", "Webhook名称不能为空")
	}
	if w.URL == "" {
		return NewValidationError("url", "URL不能为空")
	}
	if !strings.HasPrefix(w.URL, "http://") && !strings.HasPrefix(w.URL, "https://") {
		return NewValidationError("url", "URL必须以 http:// 或 https:// 开头")
	}
	if w.Status == "" {
		w.Status = WebhookActive
	}
	if w.Status != WebhookActive && w.Status != WebhookDisabled {
		return NewValidationError("status", "Webhook状态不合法")
	}
	return nil
}

func (w *Webhook) SupportsEvent(event string) bool {
	for _, e := range w.Events {
		if e == event || e == "*" {
			return true
		}
	}
	return false
}

type WebhookFilter struct {
	Status  string
	Keyword string
}

func (f WebhookFilter) Match(w *Webhook) bool {
	if f.Status != "" && w.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(w.Name), k) {
			return false
		}
	}
	return true
}

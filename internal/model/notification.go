package model

import (
	"strings"
	"time"
)

const (
	NotificationPending    = "pending"
	NotificationSent       = "sent"
	NotificationFailed     = "failed"
	NotificationDelivered  = "delivered"
)

type Notification struct {
	ID         string    `json:"id"`
	WebhookID  string    `json:"webhook_id"`
	Event      string    `json:"event"`
	Payload    string    `json:"payload"`
	Status     string    `json:"status"`
	RetryCount int       `json:"retry_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (n *Notification) Validate() error {
	n.WebhookID = strings.TrimSpace(n.WebhookID)
	n.Event = strings.TrimSpace(n.Event)
	if n.WebhookID == "" {
		return NewValidationError("webhook_id", "Webhook ID不能为空")
	}
	if n.Event == "" {
		return NewValidationError("event", "事件类型不能为空")
	}
	if n.Status == "" {
		n.Status = NotificationPending
	}
	if !isValidNotificationStatus(n.Status) {
		return NewValidationError("status", "通知状态不合法")
	}
	return nil
}

func isValidNotificationStatus(s string) bool {
	switch s {
	case NotificationPending, NotificationSent, NotificationFailed, NotificationDelivered:
		return true
	}
	return false
}

type NotificationFilter struct {
	WebhookID string
	Event     string
	Status    string
}

func (f NotificationFilter) Match(n *Notification) bool {
	if f.WebhookID != "" && n.WebhookID != f.WebhookID {
		return false
	}
	if f.Event != "" && n.Event != f.Event {
		return false
	}
	if f.Status != "" && n.Status != f.Status {
		return false
	}
	return true
}

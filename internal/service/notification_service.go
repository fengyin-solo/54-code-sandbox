package service

import (
	"sort"
	"time"

	"sandbox/internal/model"
	"sandbox/pkg/idgen"
)

func (s *Service) CreateNotification(input model.Notification) (*model.Notification, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetWebhook(input.WebhookID); err != nil {
		return nil, model.NewValidationError("webhook_id", "关联Webhook不存在")
	}
	now := time.Now()
	n := &model.Notification{
		ID:         idgen.Hex(),
		WebhookID:  input.WebhookID,
		Event:      input.Event,
		Payload:    input.Payload,
		Status:     model.NotificationPending,
		RetryCount: 0,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.store.CreateNotification(n); err != nil {
		return nil, err
	}
	s.log.Infof("创建通知: %s", n.ID)
	return n, nil
}

func (s *Service) GetNotification(id string) (*model.Notification, error) {
	return s.store.GetNotification(id)
}

func (s *Service) ListNotifications(filter model.NotificationFilter, page, size int) ([]*model.Notification, int, error) {
	all := s.store.ListNotifications()
	matched := make([]*model.Notification, 0, len(all))
	for _, n := range all {
		if filter.Match(n) {
			matched = append(matched, n)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Notification{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateNotification(id string, input model.Notification) (*model.Notification, error) {
	n, err := s.store.GetNotification(id)
	if err != nil {
		return nil, err
	}
	if input.Status != "" {
		n.Status = input.Status
	}
	if input.RetryCount >= 0 {
		n.RetryCount = input.RetryCount
	}
	if input.Payload != "" {
		n.Payload = input.Payload
	}
	n.UpdatedAt = time.Now()
	if err := n.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateNotification(n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) DeleteNotification(id string) error {
	return s.store.DeleteNotification(id)
}

func (s *Service) MarkNotificationSent(id string) (*model.Notification, error) {
	n, err := s.store.GetNotification(id)
	if err != nil {
		return nil, err
	}
	n.Status = model.NotificationSent
	n.UpdatedAt = time.Now()
	if err := s.store.UpdateNotification(n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) MarkNotificationDelivered(id string) (*model.Notification, error) {
	n, err := s.store.GetNotification(id)
	if err != nil {
		return nil, err
	}
	n.Status = model.NotificationDelivered
	n.UpdatedAt = time.Now()
	if err := s.store.UpdateNotification(n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) MarkNotificationFailed(id string) (*model.Notification, error) {
	n, err := s.store.GetNotification(id)
	if err != nil {
		return nil, err
	}
	n.Status = model.NotificationFailed
	n.RetryCount++
	n.UpdatedAt = time.Now()
	if err := s.store.UpdateNotification(n); err != nil {
		return nil, err
	}
	return n, nil
}

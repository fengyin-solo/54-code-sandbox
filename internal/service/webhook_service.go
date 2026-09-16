package service

import (
	"sort"
	"time"

	"sandbox/internal/model"
	"sandbox/pkg/idgen"
)

func (s *Service) CreateWebhook(input model.Webhook) (*model.Webhook, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	w := &model.Webhook{
		ID:        idgen.Hex(),
		Name:      input.Name,
		URL:       input.URL,
		Secret:    input.Secret,
		Events:    input.Events,
		Status:    input.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.store.CreateWebhook(w); err != nil {
		return nil, err
	}
	s.log.Infof("创建Webhook: %s", w.ID)
	return w, nil
}

func (s *Service) GetWebhook(id string) (*model.Webhook, error) {
	return s.store.GetWebhook(id)
}

func (s *Service) ListWebhooks(filter model.WebhookFilter, page, size int) ([]*model.Webhook, int, error) {
	all := s.store.ListWebhooks()
	matched := make([]*model.Webhook, 0, len(all))
	for _, w := range all {
		if filter.Match(w) {
			matched = append(matched, w)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Webhook{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateWebhook(id string, input model.Webhook) (*model.Webhook, error) {
	w, err := s.store.GetWebhook(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		w.Name = input.Name
	}
	if input.URL != "" {
		w.URL = input.URL
	}
	if input.Secret != "" {
		w.Secret = input.Secret
	}
	if input.Events != nil {
		w.Events = input.Events
	}
	if input.Status != "" {
		w.Status = input.Status
	}
	w.UpdatedAt = time.Now()
	if err := w.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateWebhook(w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *Service) DeleteWebhook(id string) error {
	return s.store.DeleteWebhook(id)
}

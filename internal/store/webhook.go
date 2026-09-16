package store

import "sandbox/internal/model"

func (s *MemoryStore) CreateWebhook(w *model.Webhook) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.webhooks {
		if exist.URL == w.URL {
			return ErrConflict
		}
	}
	s.webhooks[w.ID] = w
	return nil
}

func (s *MemoryStore) GetWebhook(id string) (*model.Webhook, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	w, ok := s.webhooks[id]
	if !ok {
		return nil, ErrNotFound
	}
	return w, nil
}

func (s *MemoryStore) ListWebhooks() []*model.Webhook {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Webhook, 0, len(s.webhooks))
	for _, w := range s.webhooks {
		list = append(list, w)
	}
	return list
}

func (s *MemoryStore) UpdateWebhook(w *model.Webhook) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.webhooks[w.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.webhooks {
		if exist.ID != w.ID && exist.URL == w.URL {
			return ErrConflict
		}
	}
	s.webhooks[w.ID] = w
	return nil
}

func (s *MemoryStore) DeleteWebhook(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.webhooks[id]; !ok {
		return ErrNotFound
	}
	delete(s.webhooks, id)
	return nil
}

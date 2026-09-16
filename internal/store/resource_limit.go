package store

import "sandbox/internal/model"

func (s *MemoryStore) CreateResourceLimit(r *model.ResourceLimit) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.resourceLimits {
		if exist.SandboxID == r.SandboxID {
			return ErrConflict
		}
	}
	s.resourceLimits[r.ID] = r
	return nil
}

func (s *MemoryStore) GetResourceLimit(id string) (*model.ResourceLimit, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.resourceLimits[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) GetResourceLimitBySandbox(sandboxID string) (*model.ResourceLimit, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, r := range s.resourceLimits {
		if r.SandboxID == sandboxID {
			return r, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListResourceLimits() []*model.ResourceLimit {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ResourceLimit, 0, len(s.resourceLimits))
	for _, r := range s.resourceLimits {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateResourceLimit(r *model.ResourceLimit) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.resourceLimits[r.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.resourceLimits {
		if exist.ID != r.ID && exist.SandboxID == r.SandboxID {
			return ErrConflict
		}
	}
	s.resourceLimits[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteResourceLimit(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.resourceLimits[id]; !ok {
		return ErrNotFound
	}
	delete(s.resourceLimits, id)
	return nil
}

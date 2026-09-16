package store

import "sandbox/internal/model"

func (s *MemoryStore) CreateRuntime(r *model.Runtime) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.runtimes[r.ID] = r
	return nil
}

func (s *MemoryStore) GetRuntime(id string) (*model.Runtime, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.runtimes[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

func (s *MemoryStore) ListRuntimes() []*model.Runtime {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Runtime, 0, len(s.runtimes))
	for _, r := range s.runtimes {
		list = append(list, r)
	}
	return list
}

func (s *MemoryStore) UpdateRuntime(r *model.Runtime) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.runtimes[r.ID]; !ok {
		return ErrNotFound
	}
	s.runtimes[r.ID] = r
	return nil
}

func (s *MemoryStore) DeleteRuntime(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.runtimes[id]; !ok {
		return ErrNotFound
	}
	delete(s.runtimes, id)
	return nil
}

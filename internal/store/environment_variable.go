package store

import "sandbox/internal/model"

func (s *MemoryStore) CreateEnvironmentVariable(e *model.EnvironmentVariable) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.environmentVariables {
		if exist.SandboxID == e.SandboxID && exist.Key == e.Key {
			return ErrConflict
		}
	}
	s.environmentVariables[e.ID] = e
	return nil
}

func (s *MemoryStore) GetEnvironmentVariable(id string) (*model.EnvironmentVariable, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.environmentVariables[id]
	if !ok {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *MemoryStore) ListEnvironmentVariables() []*model.EnvironmentVariable {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.EnvironmentVariable, 0, len(s.environmentVariables))
	for _, e := range s.environmentVariables {
		list = append(list, e)
	}
	return list
}

func (s *MemoryStore) UpdateEnvironmentVariable(e *model.EnvironmentVariable) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.environmentVariables[e.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.environmentVariables {
		if exist.ID != e.ID && exist.SandboxID == e.SandboxID && exist.Key == e.Key {
			return ErrConflict
		}
	}
	s.environmentVariables[e.ID] = e
	return nil
}

func (s *MemoryStore) DeleteEnvironmentVariable(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.environmentVariables[id]; !ok {
		return ErrNotFound
	}
	delete(s.environmentVariables, id)
	return nil
}

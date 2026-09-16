package store

import "sandbox/internal/model"

func (s *MemoryStore) CreateSandbox(sb *model.Sandbox) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.sandboxes {
		if exist.Name == sb.Name {
			return ErrConflict
		}
	}
	s.sandboxes[sb.ID] = sb
	return nil
}

func (s *MemoryStore) GetSandbox(id string) (*model.Sandbox, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sb, ok := s.sandboxes[id]
	if !ok {
		return nil, ErrNotFound
	}
	return sb, nil
}

func (s *MemoryStore) GetSandboxByName(name string) (*model.Sandbox, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, sb := range s.sandboxes {
		if sb.Name == name {
			return sb, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListSandboxes() []*model.Sandbox {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Sandbox, 0, len(s.sandboxes))
	for _, sb := range s.sandboxes {
		list = append(list, sb)
	}
	return list
}

func (s *MemoryStore) UpdateSandbox(sb *model.Sandbox) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sandboxes[sb.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.sandboxes {
		if exist.ID != sb.ID && exist.Name == sb.Name {
			return ErrConflict
		}
	}
	s.sandboxes[sb.ID] = sb
	return nil
}

func (s *MemoryStore) DeleteSandbox(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sandboxes[id]; !ok {
		return ErrNotFound
	}
	delete(s.sandboxes, id)
	return nil
}

package store

import "sandbox/internal/model"

func (s *MemoryStore) CreateArtifact(a *model.Artifact) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.artifacts[a.ID] = a
	return nil
}

func (s *MemoryStore) GetArtifact(id string) (*model.Artifact, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.artifacts[id]
	if !ok {
		return nil, ErrNotFound
	}
	return a, nil
}

func (s *MemoryStore) ListArtifacts() []*model.Artifact {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Artifact, 0, len(s.artifacts))
	for _, a := range s.artifacts {
		list = append(list, a)
	}
	return list
}

func (s *MemoryStore) UpdateArtifact(a *model.Artifact) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.artifacts[a.ID]; !ok {
		return ErrNotFound
	}
	s.artifacts[a.ID] = a
	return nil
}

func (s *MemoryStore) DeleteArtifact(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.artifacts[id]; !ok {
		return ErrNotFound
	}
	delete(s.artifacts, id)
	return nil
}

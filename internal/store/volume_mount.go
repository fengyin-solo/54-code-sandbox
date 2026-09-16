package store

import "sandbox/internal/model"

func (s *MemoryStore) CreateVolumeMount(v *model.VolumeMount) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.volumeMounts[v.ID] = v
	return nil
}

func (s *MemoryStore) GetVolumeMount(id string) (*model.VolumeMount, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.volumeMounts[id]
	if !ok {
		return nil, ErrNotFound
	}
	return v, nil
}

func (s *MemoryStore) ListVolumeMounts() []*model.VolumeMount {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.VolumeMount, 0, len(s.volumeMounts))
	for _, v := range s.volumeMounts {
		list = append(list, v)
	}
	return list
}

func (s *MemoryStore) UpdateVolumeMount(v *model.VolumeMount) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.volumeMounts[v.ID]; !ok {
		return ErrNotFound
	}
	s.volumeMounts[v.ID] = v
	return nil
}

func (s *MemoryStore) DeleteVolumeMount(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.volumeMounts[id]; !ok {
		return ErrNotFound
	}
	delete(s.volumeMounts, id)
	return nil
}

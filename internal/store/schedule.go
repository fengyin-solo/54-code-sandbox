package store

import "sandbox/internal/model"

func (s *MemoryStore) CreateSchedule(sch *model.Schedule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.schedules[sch.ID] = sch
	return nil
}

func (s *MemoryStore) GetSchedule(id string) (*model.Schedule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sch, ok := s.schedules[id]
	if !ok {
		return nil, ErrNotFound
	}
	return sch, nil
}

func (s *MemoryStore) ListSchedules() []*model.Schedule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Schedule, 0, len(s.schedules))
	for _, sch := range s.schedules {
		list = append(list, sch)
	}
	return list
}

func (s *MemoryStore) UpdateSchedule(sch *model.Schedule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.schedules[sch.ID]; !ok {
		return ErrNotFound
	}
	s.schedules[sch.ID] = sch
	return nil
}

func (s *MemoryStore) DeleteSchedule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.schedules[id]; !ok {
		return ErrNotFound
	}
	delete(s.schedules, id)
	return nil
}

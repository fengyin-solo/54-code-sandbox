package store

import "sandbox/internal/model"

func (s *MemoryStore) CreateExecutionLog(l *model.ExecutionLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.executionLogs[l.ID] = l
	return nil
}

func (s *MemoryStore) GetExecutionLog(id string) (*model.ExecutionLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	l, ok := s.executionLogs[id]
	if !ok {
		return nil, ErrNotFound
	}
	return l, nil
}

func (s *MemoryStore) ListExecutionLogs() []*model.ExecutionLog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ExecutionLog, 0, len(s.executionLogs))
	for _, l := range s.executionLogs {
		list = append(list, l)
	}
	return list
}

func (s *MemoryStore) UpdateExecutionLog(l *model.ExecutionLog) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.executionLogs[l.ID]; !ok {
		return ErrNotFound
	}
	s.executionLogs[l.ID] = l
	return nil
}

func (s *MemoryStore) DeleteExecutionLog(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.executionLogs[id]; !ok {
		return ErrNotFound
	}
	delete(s.executionLogs, id)
	return nil
}

package store

import "sandbox/internal/model"

func (s *MemoryStore) CreateExecutionTask(t *model.ExecutionTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.executionTasks[t.ID] = t
	return nil
}

func (s *MemoryStore) GetExecutionTask(id string) (*model.ExecutionTask, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.executionTasks[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (s *MemoryStore) ListExecutionTasks() []*model.ExecutionTask {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.ExecutionTask, 0, len(s.executionTasks))
	for _, t := range s.executionTasks {
		list = append(list, t)
	}
	return list
}

func (s *MemoryStore) UpdateExecutionTask(t *model.ExecutionTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.executionTasks[t.ID]; !ok {
		return ErrNotFound
	}
	s.executionTasks[t.ID] = t
	return nil
}

func (s *MemoryStore) DeleteExecutionTask(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.executionTasks[id]; !ok {
		return ErrNotFound
	}
	delete(s.executionTasks, id)
	return nil
}

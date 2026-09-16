package service

import (
	"sort"

	"sandbox/internal/model"
	"sandbox/pkg/idgen"
)

func (s *Service) CreateExecutionLog(input model.ExecutionLog) (*model.ExecutionLog, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetExecutionTask(input.TaskID); err != nil {
		return nil, model.NewValidationError("task_id", "关联任务不存在")
	}
	l := &model.ExecutionLog{
		ID:       idgen.Hex(),
		TaskID:   input.TaskID,
		Level:    input.Level,
		Message:  input.Message,
		LoggedAt: input.LoggedAt,
	}
	if err := s.store.CreateExecutionLog(l); err != nil {
		return nil, err
	}
	return l, nil
}

func (s *Service) GetExecutionLog(id string) (*model.ExecutionLog, error) {
	return s.store.GetExecutionLog(id)
}

func (s *Service) ListExecutionLogs(filter model.ExecutionLogFilter, page, size int) ([]*model.ExecutionLog, int, error) {
	all := s.store.ListExecutionLogs()
	matched := make([]*model.ExecutionLog, 0, len(all))
	for _, l := range all {
		if filter.Match(l) {
			matched = append(matched, l)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].LoggedAt.After(matched[j].LoggedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.ExecutionLog{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) DeleteExecutionLog(id string) error {
	return s.store.DeleteExecutionLog(id)
}

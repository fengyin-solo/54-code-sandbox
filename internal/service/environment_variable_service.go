package service

import (
	"sort"
	"time"

	"sandbox/internal/model"
	"sandbox/pkg/idgen"
)

func (s *Service) CreateEnvironmentVariable(input model.EnvironmentVariable) (*model.EnvironmentVariable, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetSandbox(input.SandboxID); err != nil {
		return nil, model.NewValidationError("sandbox_id", "所属沙箱不存在")
	}
	now := time.Now()
	e := &model.EnvironmentVariable{
		ID:        idgen.Hex(),
		SandboxID: input.SandboxID,
		Key:       input.Key,
		Value:     input.Value,
		Status:    input.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.store.CreateEnvironmentVariable(e); err != nil {
		return nil, err
	}
	s.log.Infof("创建环境变量: %s", e.ID)
	return e, nil
}

func (s *Service) GetEnvironmentVariable(id string) (*model.EnvironmentVariable, error) {
	return s.store.GetEnvironmentVariable(id)
}

func (s *Service) ListEnvironmentVariables(filter model.EnvironmentVariableFilter, page, size int) ([]*model.EnvironmentVariable, int, error) {
	all := s.store.ListEnvironmentVariables()
	matched := make([]*model.EnvironmentVariable, 0, len(all))
	for _, e := range all {
		if filter.Match(e) {
			matched = append(matched, e)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.EnvironmentVariable{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateEnvironmentVariable(id string, input model.EnvironmentVariable) (*model.EnvironmentVariable, error) {
	e, err := s.store.GetEnvironmentVariable(id)
	if err != nil {
		return nil, err
	}
	if input.Key != "" {
		e.Key = input.Key
	}
	if input.Value != "" {
		e.Value = input.Value
	}
	if input.Status != "" {
		e.Status = input.Status
	}
	e.UpdatedAt = time.Now()
	if err := e.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateEnvironmentVariable(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *Service) DeleteEnvironmentVariable(id string) error {
	return s.store.DeleteEnvironmentVariable(id)
}

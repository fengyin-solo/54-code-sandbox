package service

import (
	"sort"
	"time"

	"sandbox/internal/model"
	"sandbox/pkg/idgen"
)

func (s *Service) CreateRuntime(input model.Runtime) (*model.Runtime, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetSandbox(input.SandboxID); err != nil {
		return nil, model.NewValidationError("sandbox_id", "所属沙箱不存在")
	}
	now := time.Now()
	r := &model.Runtime{
		ID:            idgen.Hex(),
		SandboxID:     input.SandboxID,
		Name:          input.Name,
		Version:       input.Version,
		Command:       input.Command,
		MemoryLimitMB: input.MemoryLimitMB,
		CPULimit:      input.CPULimit,
		TimeoutMs:     input.TimeoutMs,
		Status:        input.Status,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := s.store.CreateRuntime(r); err != nil {
		return nil, err
	}
	s.log.Infof("创建运行时: %s", r.ID)
	return r, nil
}

func (s *Service) GetRuntime(id string) (*model.Runtime, error) {
	return s.store.GetRuntime(id)
}

func (s *Service) ListRuntimes(filter model.RuntimeFilter, page, size int) ([]*model.Runtime, int, error) {
	all := s.store.ListRuntimes()
	matched := make([]*model.Runtime, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Runtime{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateRuntime(id string, input model.Runtime) (*model.Runtime, error) {
	r, err := s.store.GetRuntime(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		r.Name = input.Name
	}
	if input.Version != "" {
		r.Version = input.Version
	}
	if input.Command != "" {
		r.Command = input.Command
	}
	if input.MemoryLimitMB > 0 {
		r.MemoryLimitMB = input.MemoryLimitMB
	}
	if input.CPULimit > 0 {
		r.CPULimit = input.CPULimit
	}
	if input.TimeoutMs > 0 {
		r.TimeoutMs = input.TimeoutMs
	}
	if input.Status != "" {
		r.Status = input.Status
	}
	r.UpdatedAt = time.Now()
	if err := r.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateRuntime(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) DeleteRuntime(id string) error {
	return s.store.DeleteRuntime(id)
}

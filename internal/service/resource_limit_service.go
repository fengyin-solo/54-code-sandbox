package service

import (
	"sort"
	"time"

	"sandbox/internal/model"
	"sandbox/pkg/idgen"
)

func (s *Service) CreateResourceLimit(input model.ResourceLimit) (*model.ResourceLimit, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetSandbox(input.SandboxID); err != nil {
		return nil, model.NewValidationError("sandbox_id", "所属沙箱不存在")
	}
	now := time.Now()
	r := &model.ResourceLimit{
		ID:           idgen.Hex(),
		SandboxID:    input.SandboxID,
		CPUQuota:     input.CPUQuota,
		MemoryMB:     input.MemoryMB,
		DiskMB:       input.DiskMB,
		MaxProcesses: input.MaxProcesses,
		Status:       input.Status,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.store.CreateResourceLimit(r); err != nil {
		return nil, err
	}
	s.log.Infof("创建资源限制: %s", r.ID)
	return r, nil
}

func (s *Service) GetResourceLimit(id string) (*model.ResourceLimit, error) {
	return s.store.GetResourceLimit(id)
}

func (s *Service) ListResourceLimits(filter model.ResourceLimitFilter, page, size int) ([]*model.ResourceLimit, int, error) {
	all := s.store.ListResourceLimits()
	matched := make([]*model.ResourceLimit, 0, len(all))
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
		return []*model.ResourceLimit{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateResourceLimit(id string, input model.ResourceLimit) (*model.ResourceLimit, error) {
	r, err := s.store.GetResourceLimit(id)
	if err != nil {
		return nil, err
	}
	if input.CPUQuota > 0 {
		r.CPUQuota = input.CPUQuota
	}
	if input.MemoryMB > 0 {
		r.MemoryMB = input.MemoryMB
	}
	if input.DiskMB > 0 {
		r.DiskMB = input.DiskMB
	}
	if input.MaxProcesses > 0 {
		r.MaxProcesses = input.MaxProcesses
	}
	if input.Status != "" {
		r.Status = input.Status
	}
	r.UpdatedAt = time.Now()
	if err := r.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateResourceLimit(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) DeleteResourceLimit(id string) error {
	return s.store.DeleteResourceLimit(id)
}

package service

import (
	"sort"
	"time"

	"sandbox/internal/model"
	"sandbox/pkg/idgen"
)

func (s *Service) CreateSecurityPolicy(input model.SecurityPolicy) (*model.SecurityPolicy, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	p := &model.SecurityPolicy{
		ID:              idgen.Hex(),
		Name:            input.Name,
		AllowNetwork:    input.AllowNetwork,
		AllowFileSystem: input.AllowFileSystem,
		MaxMemoryMB:     input.MaxMemoryMB,
		MaxCPUMs:        input.MaxCPUMs,
		MaxTimeoutMs:    input.MaxTimeoutMs,
		BlockedSyscalls: input.BlockedSyscalls,
		Status:          input.Status,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if err := s.store.CreateSecurityPolicy(p); err != nil {
		return nil, err
	}
	s.log.Infof("创建安全策略: %s", p.ID)
	return p, nil
}

func (s *Service) GetSecurityPolicy(id string) (*model.SecurityPolicy, error) {
	return s.store.GetSecurityPolicy(id)
}

func (s *Service) ListSecurityPolicies(filter model.SecurityPolicyFilter, page, size int) ([]*model.SecurityPolicy, int, error) {
	all := s.store.ListSecurityPolicies()
	matched := make([]*model.SecurityPolicy, 0, len(all))
	for _, p := range all {
		if filter.Match(p) {
			matched = append(matched, p)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.SecurityPolicy{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateSecurityPolicy(id string, input model.SecurityPolicy) (*model.SecurityPolicy, error) {
	p, err := s.store.GetSecurityPolicy(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		p.Name = input.Name
	}
	p.AllowNetwork = input.AllowNetwork
	p.AllowFileSystem = input.AllowFileSystem
	if input.MaxMemoryMB > 0 {
		p.MaxMemoryMB = input.MaxMemoryMB
	}
	if input.MaxCPUMs > 0 {
		p.MaxCPUMs = input.MaxCPUMs
	}
	if input.MaxTimeoutMs > 0 {
		p.MaxTimeoutMs = input.MaxTimeoutMs
	}
	if input.BlockedSyscalls != nil {
		p.BlockedSyscalls = input.BlockedSyscalls
	}
	if input.Status != "" {
		p.Status = input.Status
	}
	p.UpdatedAt = time.Now()
	if err := p.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateSecurityPolicy(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) DeleteSecurityPolicy(id string) error {
	return s.store.DeleteSecurityPolicy(id)
}

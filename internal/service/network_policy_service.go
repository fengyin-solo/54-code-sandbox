package service

import (
	"sort"
	"time"

	"sandbox/internal/model"
	"sandbox/pkg/idgen"
)

func (s *Service) CreateNetworkPolicy(input model.NetworkPolicy) (*model.NetworkPolicy, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetSandbox(input.SandboxID); err != nil {
		return nil, model.NewValidationError("sandbox_id", "所属沙箱不存在")
	}
	now := time.Now()
	n := &model.NetworkPolicy{
		ID:           idgen.Hex(),
		SandboxID:    input.SandboxID,
		AllowedHosts: input.AllowedHosts,
		BlockedHosts: input.BlockedHosts,
		Status:       input.Status,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := s.store.CreateNetworkPolicy(n); err != nil {
		return nil, err
	}
	s.log.Infof("创建网络策略: %s", n.ID)
	return n, nil
}

func (s *Service) GetNetworkPolicy(id string) (*model.NetworkPolicy, error) {
	return s.store.GetNetworkPolicy(id)
}

func (s *Service) ListNetworkPolicies(filter model.NetworkPolicyFilter, page, size int) ([]*model.NetworkPolicy, int, error) {
	all := s.store.ListNetworkPolicies()
	matched := make([]*model.NetworkPolicy, 0, len(all))
	for _, n := range all {
		if filter.Match(n) {
			matched = append(matched, n)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.NetworkPolicy{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateNetworkPolicy(id string, input model.NetworkPolicy) (*model.NetworkPolicy, error) {
	n, err := s.store.GetNetworkPolicy(id)
	if err != nil {
		return nil, err
	}
	if input.AllowedHosts != nil {
		n.AllowedHosts = input.AllowedHosts
	}
	if input.BlockedHosts != nil {
		n.BlockedHosts = input.BlockedHosts
	}
	if input.Status != "" {
		n.Status = input.Status
	}
	n.UpdatedAt = time.Now()
	if err := n.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateNetworkPolicy(n); err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Service) DeleteNetworkPolicy(id string) error {
	return s.store.DeleteNetworkPolicy(id)
}

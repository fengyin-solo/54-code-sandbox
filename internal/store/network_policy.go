package store

import "sandbox/internal/model"

func (s *MemoryStore) CreateNetworkPolicy(n *model.NetworkPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.networkPolicies {
		if exist.SandboxID == n.SandboxID {
			return ErrConflict
		}
	}
	s.networkPolicies[n.ID] = n
	return nil
}

func (s *MemoryStore) GetNetworkPolicy(id string) (*model.NetworkPolicy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.networkPolicies[id]
	if !ok {
		return nil, ErrNotFound
	}
	return n, nil
}

func (s *MemoryStore) ListNetworkPolicies() []*model.NetworkPolicy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.NetworkPolicy, 0, len(s.networkPolicies))
	for _, n := range s.networkPolicies {
		list = append(list, n)
	}
	return list
}

func (s *MemoryStore) UpdateNetworkPolicy(n *model.NetworkPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.networkPolicies[n.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.networkPolicies {
		if exist.ID != n.ID && exist.SandboxID == n.SandboxID {
			return ErrConflict
		}
	}
	s.networkPolicies[n.ID] = n
	return nil
}

func (s *MemoryStore) DeleteNetworkPolicy(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.networkPolicies[id]; !ok {
		return ErrNotFound
	}
	delete(s.networkPolicies, id)
	return nil
}

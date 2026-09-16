package store

import "sandbox/internal/model"

func (s *MemoryStore) CreateSecurityPolicy(p *model.SecurityPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.securityPolicies {
		if exist.Name == p.Name {
			return ErrConflict
		}
	}
	s.securityPolicies[p.ID] = p
	return nil
}

func (s *MemoryStore) GetSecurityPolicy(id string) (*model.SecurityPolicy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.securityPolicies[id]
	if !ok {
		return nil, ErrNotFound
	}
	return p, nil
}

func (s *MemoryStore) ListSecurityPolicies() []*model.SecurityPolicy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.SecurityPolicy, 0, len(s.securityPolicies))
	for _, p := range s.securityPolicies {
		list = append(list, p)
	}
	return list
}

func (s *MemoryStore) UpdateSecurityPolicy(p *model.SecurityPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.securityPolicies[p.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.securityPolicies {
		if exist.ID != p.ID && exist.Name == p.Name {
			return ErrConflict
		}
	}
	s.securityPolicies[p.ID] = p
	return nil
}

func (s *MemoryStore) DeleteSecurityPolicy(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.securityPolicies[id]; !ok {
		return ErrNotFound
	}
	delete(s.securityPolicies, id)
	return nil
}

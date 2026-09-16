package service

import (
	"sort"
	"time"

	"sandbox/internal/model"
	"sandbox/pkg/idgen"
)

func (s *Service) CreateSandbox(input model.Sandbox) (*model.Sandbox, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	sb := &model.Sandbox{
		ID:        idgen.Hex(),
		Name:      input.Name,
		Language:  input.Language,
		Image:     input.Image,
		Status:    input.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.store.CreateSandbox(sb); err != nil {
		return nil, err
	}
	s.log.Infof("创建沙箱: %s", sb.ID)
	return sb, nil
}

func (s *Service) GetSandbox(id string) (*model.Sandbox, error) {
	return s.store.GetSandbox(id)
}

func (s *Service) ListSandboxes(filter model.SandboxFilter, page, size int) ([]*model.Sandbox, int, error) {
	all := s.store.ListSandboxes()
	matched := make([]*model.Sandbox, 0, len(all))
	for _, sb := range all {
		if filter.Match(sb) {
			matched = append(matched, sb)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Sandbox{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateSandbox(id string, input model.Sandbox) (*model.Sandbox, error) {
	sb, err := s.store.GetSandbox(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		sb.Name = input.Name
	}
	if input.Language != "" {
		sb.Language = input.Language
	}
	if input.Image != "" {
		sb.Image = input.Image
	}
	if input.Status != "" {
		sb.Status = input.Status
	}
	sb.UpdatedAt = time.Now()
	if err := sb.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateSandbox(sb); err != nil {
		return nil, err
	}
	return sb, nil
}

func (s *Service) DeleteSandbox(id string) error {
	return s.store.DeleteSandbox(id)
}

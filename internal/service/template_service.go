package service

import (
	"sort"
	"time"

	"sandbox/internal/model"
	"sandbox/pkg/idgen"
)

func (s *Service) CreateTemplate(input model.Template) (*model.Template, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	t := &model.Template{
		ID:          idgen.Hex(),
		Name:        input.Name,
		Language:    input.Language,
		BaseImage:   input.BaseImage,
		Description: input.Description,
		Status:      input.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.store.CreateTemplate(t); err != nil {
		return nil, err
	}
	s.log.Infof("创建模板: %s", t.ID)
	return t, nil
}

func (s *Service) GetTemplate(id string) (*model.Template, error) {
	return s.store.GetTemplate(id)
}

func (s *Service) ListTemplates(filter model.TemplateFilter, page, size int) ([]*model.Template, int, error) {
	all := s.store.ListTemplates()
	matched := make([]*model.Template, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Template{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateTemplate(id string, input model.Template) (*model.Template, error) {
	t, err := s.store.GetTemplate(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		t.Name = input.Name
	}
	if input.Language != "" {
		t.Language = input.Language
	}
	if input.BaseImage != "" {
		t.BaseImage = input.BaseImage
	}
	if input.Description != "" {
		t.Description = input.Description
	}
	if input.Status != "" {
		t.Status = input.Status
	}
	t.UpdatedAt = time.Now()
	if err := t.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateTemplate(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) DeleteTemplate(id string) error {
	return s.store.DeleteTemplate(id)
}

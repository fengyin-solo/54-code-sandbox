package service

import (
	"sort"
	"time"

	"sandbox/internal/model"
	"sandbox/pkg/idgen"
)

func (s *Service) CreateArtifact(input model.Artifact) (*model.Artifact, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetExecutionTask(input.TaskID); err != nil {
		return nil, model.NewValidationError("task_id", "关联任务不存在")
	}
	now := time.Now()
	a := &model.Artifact{
		ID:        idgen.Hex(),
		TaskID:    input.TaskID,
		Name:      input.Name,
		Path:      input.Path,
		SizeBytes: input.SizeBytes,
		Status:    input.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.store.CreateArtifact(a); err != nil {
		return nil, err
	}
	s.log.Infof("创建产物: %s", a.ID)
	return a, nil
}

func (s *Service) GetArtifact(id string) (*model.Artifact, error) {
	return s.store.GetArtifact(id)
}

func (s *Service) ListArtifacts(filter model.ArtifactFilter, page, size int) ([]*model.Artifact, int, error) {
	all := s.store.ListArtifacts()
	matched := make([]*model.Artifact, 0, len(all))
	for _, a := range all {
		if filter.Match(a) {
			matched = append(matched, a)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Artifact{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateArtifact(id string, input model.Artifact) (*model.Artifact, error) {
	a, err := s.store.GetArtifact(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		a.Name = input.Name
	}
	if input.Path != "" {
		a.Path = input.Path
	}
	if input.SizeBytes >= 0 {
		a.SizeBytes = input.SizeBytes
	}
	if input.Status != "" {
		a.Status = input.Status
	}
	a.UpdatedAt = time.Now()
	if err := a.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateArtifact(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) DeleteArtifact(id string) error {
	return s.store.DeleteArtifact(id)
}

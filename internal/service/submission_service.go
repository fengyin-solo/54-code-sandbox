package service

import (
	"sort"
	"time"

	"sandbox/internal/model"
	"sandbox/pkg/idgen"
)

func (s *Service) CreateSubmission(input model.Submission) (*model.Submission, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetExecutionTask(input.TaskID); err != nil {
		return nil, model.NewValidationError("task_id", "关联任务不存在")
	}
	now := time.Now()
	sub := &model.Submission{
		ID:         idgen.Hex(),
		TaskID:     input.TaskID,
		Submitter:  input.Submitter,
		SourceCode: input.SourceCode,
		Status:     input.Status,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.store.CreateSubmission(sub); err != nil {
		return nil, err
	}
	s.log.Infof("创建提交: %s", sub.ID)
	return sub, nil
}

func (s *Service) GetSubmission(id string) (*model.Submission, error) {
	return s.store.GetSubmission(id)
}

func (s *Service) ListSubmissions(filter model.SubmissionFilter, page, size int) ([]*model.Submission, int, error) {
	all := s.store.ListSubmissions()
	matched := make([]*model.Submission, 0, len(all))
	for _, sub := range all {
		if filter.Match(sub) {
			matched = append(matched, sub)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Submission{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateSubmission(id string, input model.Submission) (*model.Submission, error) {
	sub, err := s.store.GetSubmission(id)
	if err != nil {
		return nil, err
	}
	if input.Submitter != "" {
		sub.Submitter = input.Submitter
	}
	if input.SourceCode != "" {
		sub.SourceCode = input.SourceCode
	}
	if input.Status != "" {
		sub.Status = input.Status
	}
	sub.UpdatedAt = time.Now()
	if err := sub.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateSubmission(sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *Service) DeleteSubmission(id string) error {
	return s.store.DeleteSubmission(id)
}

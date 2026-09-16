package model

import (
	"strings"
	"time"
)

const (
	SubmissionPending   = "pending"
	SubmissionAccepted  = "accepted"
	SubmissionRejected  = "rejected"
)

type Submission struct {
	ID         string    `json:"id"`
	TaskID     string    `json:"task_id"`
	Submitter  string    `json:"submitter"`
	SourceCode string    `json:"source_code"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (s *Submission) Validate() error {
	s.TaskID = strings.TrimSpace(s.TaskID)
	s.Submitter = strings.TrimSpace(s.Submitter)
	s.SourceCode = strings.TrimSpace(s.SourceCode)
	if s.TaskID == "" {
		return NewValidationError("task_id", "任务ID不能为空")
	}
	if s.Submitter == "" {
		return NewValidationError("submitter", "提交人不能为空")
	}
	if s.SourceCode == "" {
		return NewValidationError("source_code", "源代码不能为空")
	}
	if s.Status == "" {
		s.Status = SubmissionPending
	}
	if s.Status != SubmissionPending && s.Status != SubmissionAccepted && s.Status != SubmissionRejected {
		return NewValidationError("status", "提交状态不合法")
	}
	return nil
}

type SubmissionFilter struct {
	TaskID    string
	Submitter string
	Status    string
}

func (f SubmissionFilter) Match(s *Submission) bool {
	if f.TaskID != "" && s.TaskID != f.TaskID {
		return false
	}
	if f.Submitter != "" && s.Submitter != f.Submitter {
		return false
	}
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	return true
}

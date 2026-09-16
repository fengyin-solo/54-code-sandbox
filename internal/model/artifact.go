package model

import (
	"strings"
	"time"
)

const (
	ArtifactActive   = "active"
	ArtifactArchived = "archived"
)

type Artifact struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"task_id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	SizeBytes int64     `json:"size_bytes"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *Artifact) Validate() error {
	a.TaskID = strings.TrimSpace(a.TaskID)
	a.Name = strings.TrimSpace(a.Name)
	a.Path = strings.TrimSpace(a.Path)
	if a.TaskID == "" {
		return NewValidationError("task_id", "任务ID不能为空")
	}
	if a.Name == "" {
		return NewValidationError("name", "产物名称不能为空")
	}
	if a.Path == "" {
		return NewValidationError("path", "路径不能为空")
	}
	if a.SizeBytes < 0 {
		return NewValidationError("size_bytes", "大小不能为负数")
	}
	if a.Status == "" {
		a.Status = ArtifactActive
	}
	if a.Status != ArtifactActive && a.Status != ArtifactArchived {
		return NewValidationError("status", "产物状态不合法")
	}
	return nil
}

type ArtifactFilter struct {
	TaskID string
	Status string
	Name   string
}

func (f ArtifactFilter) Match(a *Artifact) bool {
	if f.TaskID != "" && a.TaskID != f.TaskID {
		return false
	}
	if f.Status != "" && a.Status != f.Status {
		return false
	}
	if f.Name != "" {
		k := strings.ToLower(strings.TrimSpace(f.Name))
		if k != "" && !strings.Contains(strings.ToLower(a.Name), k) {
			return false
		}
	}
	return true
}

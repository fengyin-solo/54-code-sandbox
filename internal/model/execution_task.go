package model

import (
	"strings"
	"time"
)

const (
	TaskPending   = "pending"
	TaskRunning   = "running"
	TaskCompleted = "completed"
	TaskFailed    = "failed"
	TaskTimeout   = "timeout"
	TaskKilled    = "killed"
)

var taskTransitions = map[string]map[string]bool{
	TaskPending:   {TaskRunning: true},
	TaskRunning:   {TaskCompleted: true, TaskFailed: true, TaskTimeout: true, TaskKilled: true},
	TaskCompleted: {},
	TaskFailed:    {},
	TaskTimeout:   {},
	TaskKilled:    {},
}

func TaskCanTransition(from, to string) bool {
	if m, ok := taskTransitions[from]; ok {
		return m[to]
	}
	return false
}

type ExecutionTask struct {
	ID          string    `json:"id"`
	SandboxID   string    `json:"sandbox_id"`
	Code        string    `json:"code"`
	Language    string    `json:"language"`
	Stdin       string    `json:"stdin"`
	Status      string    `json:"status"`
	ExitCode    int       `json:"exit_code"`
	Stdout      string    `json:"stdout"`
	Stderr      string    `json:"stderr"`
	DurationMs  int       `json:"duration_ms"`
	CreatedAt   time.Time `json:"created_at"`
	StartedAt   *time.Time `json:"started_at,omitempty"`
	FinishedAt  *time.Time `json:"finished_at,omitempty"`
}

func (t *ExecutionTask) Validate() error {
	t.SandboxID = strings.TrimSpace(t.SandboxID)
	t.Language = strings.TrimSpace(t.Language)
	t.Code = strings.TrimSpace(t.Code)
	if t.SandboxID == "" {
		return NewValidationError("sandbox_id", "沙箱ID不能为空")
	}
	if t.Language == "" {
		return NewValidationError("language", "语言不能为空")
	}
	if t.Code == "" {
		return NewValidationError("code", "代码不能为空")
	}
	if t.Status == "" {
		t.Status = TaskPending
	}
	if !isValidTaskStatus(t.Status) {
		return NewValidationError("status", "任务状态不合法")
	}
	return nil
}

func isValidTaskStatus(s string) bool {
	switch s {
	case TaskPending, TaskRunning, TaskCompleted, TaskFailed, TaskTimeout, TaskKilled:
		return true
	}
	return false
}

type ExecutionTaskFilter struct {
	SandboxID string
	Language  string
	Status    string
}

func (f ExecutionTaskFilter) Match(t *ExecutionTask) bool {
	if f.SandboxID != "" && t.SandboxID != f.SandboxID {
		return false
	}
	if f.Language != "" && t.Language != f.Language {
		return false
	}
	if f.Status != "" && t.Status != f.Status {
		return false
	}
	return true
}

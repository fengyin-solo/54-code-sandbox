package model

import (
	"strings"
	"time"
)

const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

type ExecutionLog struct {
	ID        string    `json:"id"`
	TaskID    string    `json:"task_id"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	LoggedAt  time.Time `json:"logged_at"`
}

func (l *ExecutionLog) Validate() error {
	l.TaskID = strings.TrimSpace(l.TaskID)
	l.Level = strings.TrimSpace(l.Level)
	l.Message = strings.TrimSpace(l.Message)
	if l.TaskID == "" {
		return NewValidationError("task_id", "任务ID不能为空")
	}
	if l.Message == "" {
		return NewValidationError("message", "日志内容不能为空")
	}
	if l.Level == "" {
		l.Level = LogLevelInfo
	}
	if l.Level != LogLevelDebug && l.Level != LogLevelInfo && l.Level != LogLevelWarn && l.Level != LogLevelError {
		return NewValidationError("level", "日志级别不合法")
	}
	return nil
}

type ExecutionLogFilter struct {
	TaskID string
	Level  string
}

func (f ExecutionLogFilter) Match(l *ExecutionLog) bool {
	if f.TaskID != "" && l.TaskID != f.TaskID {
		return false
	}
	if f.Level != "" && l.Level != f.Level {
		return false
	}
	return true
}

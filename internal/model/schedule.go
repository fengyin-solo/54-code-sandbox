package model

import (
	"strings"
	"time"
)

const (
	SchedulePending   = "pending"
	ScheduleRunning   = "running"
	SchedulePaused    = "paused"
	ScheduleCompleted = "completed"
)

var scheduleTransitions = map[string]map[string]bool{
	SchedulePending:   {ScheduleRunning: true, SchedulePaused: true, ScheduleCompleted: true},
	ScheduleRunning:   {SchedulePaused: true, ScheduleCompleted: true},
	SchedulePaused:    {ScheduleRunning: true, ScheduleCompleted: true},
	ScheduleCompleted: {},
}

func ScheduleCanTransition(from, to string) bool {
	if m, ok := scheduleTransitions[from]; ok {
		return m[to]
	}
	return false
}

type Schedule struct {
	ID         string     `json:"id"`
	SandboxID  string     `json:"sandbox_id"`
	Code       string     `json:"code"`
	CronExpr   string     `json:"cron_expr"`
	Status     string     `json:"status"`
	LastRunAt  *time.Time `json:"last_run_at,omitempty"`
	NextRunAt  *time.Time `json:"next_run_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

func (s *Schedule) Validate() error {
	s.SandboxID = strings.TrimSpace(s.SandboxID)
	s.Code = strings.TrimSpace(s.Code)
	s.CronExpr = strings.TrimSpace(s.CronExpr)
	if s.SandboxID == "" {
		return NewValidationError("sandbox_id", "沙箱ID不能为空")
	}
	if s.Code == "" {
		return NewValidationError("code", "代码不能为空")
	}
	if s.CronExpr == "" {
		return NewValidationError("cron_expr", "Cron表达式不能为空")
	}
	if s.Status == "" {
		s.Status = SchedulePending
	}
	if !isValidScheduleStatus(s.Status) {
		return NewValidationError("status", "定时任务状态不合法")
	}
	return nil
}

func isValidScheduleStatus(st string) bool {
	switch st {
	case SchedulePending, ScheduleRunning, SchedulePaused, ScheduleCompleted:
		return true
	}
	return false
}

type ScheduleFilter struct {
	SandboxID string
	Status    string
}

func (f ScheduleFilter) Match(s *Schedule) bool {
	if f.SandboxID != "" && s.SandboxID != f.SandboxID {
		return false
	}
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	return true
}

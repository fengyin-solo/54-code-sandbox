package service

import (
	"sandbox/internal/model"
)

type BatchExecutionResult struct {
	Success bool   `json:"success"`
	TaskID  string `json:"task_id,omitempty"`
	Error   string `json:"error,omitempty"`
}

func (s *Service) BatchCreateExecutionTasks(inputs []model.ExecutionTask) ([]*BatchExecutionResult, error) {
	results := make([]*BatchExecutionResult, 0, len(inputs))
	for _, input := range inputs {
		t, err := s.CreateExecutionTask(input)
		if err != nil {
			results = append(results, &BatchExecutionResult{Success: false, Error: err.Error()})
			continue
		}
		results = append(results, &BatchExecutionResult{Success: true, TaskID: t.ID})
	}
	return results, nil
}

type SnapshotExport struct {
	Sandboxes        []*model.Sandbox        `json:"sandboxes"`
	Runtimes         []*model.Runtime         `json:"runtimes"`
	ExecutionTasks   []*model.ExecutionTask   `json:"execution_tasks"`
	SecurityPolicies []*model.SecurityPolicy  `json:"security_policies"`
	ResourceLimits   []*model.ResourceLimit   `json:"resource_limits"`
	Templates        []*model.Template        `json:"templates"`
	Submissions      []*model.Submission      `json:"submissions"`
	Schedules        []*model.Schedule        `json:"schedules"`
	ExecutionLogs    []*model.ExecutionLog    `json:"execution_logs"`
	AuditLogs        []*model.AuditLog        `json:"audit_logs"`
}

func (s *Service) ExportSnapshot() (*SnapshotExport, error) {
	return &SnapshotExport{
		Sandboxes:        s.store.ListSandboxes(),
		Runtimes:         s.store.ListRuntimes(),
		ExecutionTasks:   s.store.ListExecutionTasks(),
		SecurityPolicies: s.store.ListSecurityPolicies(),
		ResourceLimits:   s.store.ListResourceLimits(),
		Templates:        s.store.ListTemplates(),
		Submissions:      s.store.ListSubmissions(),
		Schedules:        s.store.ListSchedules(),
		ExecutionLogs:    s.store.ListExecutionLogs(),
		AuditLogs:        s.store.ListAuditLogs(),
	}, nil
}

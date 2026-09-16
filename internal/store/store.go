// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"sandbox/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	CreateSandbox(s *model.Sandbox) error
	GetSandbox(id string) (*model.Sandbox, error)
	GetSandboxByName(name string) (*model.Sandbox, error)
	ListSandboxes() []*model.Sandbox
	UpdateSandbox(s *model.Sandbox) error
	DeleteSandbox(id string) error

	CreateRuntime(r *model.Runtime) error
	GetRuntime(id string) (*model.Runtime, error)
	ListRuntimes() []*model.Runtime
	UpdateRuntime(r *model.Runtime) error
	DeleteRuntime(id string) error

	CreateExecutionTask(t *model.ExecutionTask) error
	GetExecutionTask(id string) (*model.ExecutionTask, error)
	ListExecutionTasks() []*model.ExecutionTask
	UpdateExecutionTask(t *model.ExecutionTask) error
	DeleteExecutionTask(id string) error

	CreateSecurityPolicy(p *model.SecurityPolicy) error
	GetSecurityPolicy(id string) (*model.SecurityPolicy, error)
	ListSecurityPolicies() []*model.SecurityPolicy
	UpdateSecurityPolicy(p *model.SecurityPolicy) error
	DeleteSecurityPolicy(id string) error

	CreateResourceLimit(r *model.ResourceLimit) error
	GetResourceLimit(id string) (*model.ResourceLimit, error)
	GetResourceLimitBySandbox(sandboxID string) (*model.ResourceLimit, error)
	ListResourceLimits() []*model.ResourceLimit
	UpdateResourceLimit(r *model.ResourceLimit) error
	DeleteResourceLimit(id string) error

	CreateTemplate(t *model.Template) error
	GetTemplate(id string) (*model.Template, error)
	ListTemplates() []*model.Template
	UpdateTemplate(t *model.Template) error
	DeleteTemplate(id string) error

	CreateSubmission(s *model.Submission) error
	GetSubmission(id string) (*model.Submission, error)
	ListSubmissions() []*model.Submission
	UpdateSubmission(s *model.Submission) error
	DeleteSubmission(id string) error

	CreateSchedule(s *model.Schedule) error
	GetSchedule(id string) (*model.Schedule, error)
	ListSchedules() []*model.Schedule
	UpdateSchedule(s *model.Schedule) error
	DeleteSchedule(id string) error

	CreateExecutionLog(l *model.ExecutionLog) error
	GetExecutionLog(id string) (*model.ExecutionLog, error)
	ListExecutionLogs() []*model.ExecutionLog
	UpdateExecutionLog(l *model.ExecutionLog) error
	DeleteExecutionLog(id string) error

	CreateAuditLog(a *model.AuditLog) error
	GetAuditLog(id string) (*model.AuditLog, error)
	ListAuditLogs() []*model.AuditLog
	UpdateAuditLog(a *model.AuditLog) error
	DeleteAuditLog(id string) error

	CreateEnvironmentVariable(e *model.EnvironmentVariable) error
	GetEnvironmentVariable(id string) (*model.EnvironmentVariable, error)
	ListEnvironmentVariables() []*model.EnvironmentVariable
	UpdateEnvironmentVariable(e *model.EnvironmentVariable) error
	DeleteEnvironmentVariable(id string) error

	CreateNetworkPolicy(n *model.NetworkPolicy) error
	GetNetworkPolicy(id string) (*model.NetworkPolicy, error)
	ListNetworkPolicies() []*model.NetworkPolicy
	UpdateNetworkPolicy(n *model.NetworkPolicy) error
	DeleteNetworkPolicy(id string) error

	CreateVolumeMount(v *model.VolumeMount) error
	GetVolumeMount(id string) (*model.VolumeMount, error)
	ListVolumeMounts() []*model.VolumeMount
	UpdateVolumeMount(v *model.VolumeMount) error
	DeleteVolumeMount(id string) error

	CreateArtifact(a *model.Artifact) error
	GetArtifact(id string) (*model.Artifact, error)
	ListArtifacts() []*model.Artifact
	UpdateArtifact(a *model.Artifact) error
	DeleteArtifact(id string) error

	CreateWebhook(w *model.Webhook) error
	GetWebhook(id string) (*model.Webhook, error)
	ListWebhooks() []*model.Webhook
	UpdateWebhook(w *model.Webhook) error
	DeleteWebhook(id string) error

	CreateNotification(n *model.Notification) error
	GetNotification(id string) (*model.Notification, error)
	ListNotifications() []*model.Notification
	UpdateNotification(n *model.Notification) error
	DeleteNotification(id string) error
}

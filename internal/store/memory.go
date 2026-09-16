package store

import (
	"sync"

	"sandbox/internal/model"
)

type MemoryStore struct {
	mu                   sync.RWMutex
	sandboxes            map[string]*model.Sandbox
	runtimes             map[string]*model.Runtime
	executionTasks       map[string]*model.ExecutionTask
	securityPolicies     map[string]*model.SecurityPolicy
	resourceLimits       map[string]*model.ResourceLimit
	templates            map[string]*model.Template
	submissions          map[string]*model.Submission
	schedules            map[string]*model.Schedule
	executionLogs        map[string]*model.ExecutionLog
	auditLogs            map[string]*model.AuditLog
	environmentVariables map[string]*model.EnvironmentVariable
	networkPolicies      map[string]*model.NetworkPolicy
	volumeMounts         map[string]*model.VolumeMount
	artifacts            map[string]*model.Artifact
	webhooks             map[string]*model.Webhook
	notifications        map[string]*model.Notification
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		sandboxes:            make(map[string]*model.Sandbox),
		runtimes:             make(map[string]*model.Runtime),
		executionTasks:       make(map[string]*model.ExecutionTask),
		securityPolicies:     make(map[string]*model.SecurityPolicy),
		resourceLimits:       make(map[string]*model.ResourceLimit),
		templates:            make(map[string]*model.Template),
		submissions:          make(map[string]*model.Submission),
		schedules:            make(map[string]*model.Schedule),
		executionLogs:        make(map[string]*model.ExecutionLog),
		auditLogs:            make(map[string]*model.AuditLog),
		environmentVariables: make(map[string]*model.EnvironmentVariable),
		networkPolicies:      make(map[string]*model.NetworkPolicy),
		volumeMounts:         make(map[string]*model.VolumeMount),
		artifacts:            make(map[string]*model.Artifact),
		webhooks:             make(map[string]*model.Webhook),
		notifications:        make(map[string]*model.Notification),
	}
}

var _ Store = (*MemoryStore)(nil)

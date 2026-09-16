package service

import (
	"sort"
	"time"

	"sandbox/internal/model"
	"sandbox/pkg/idgen"
)

func (s *Service) CreateExecutionTask(input model.ExecutionTask) (*model.ExecutionTask, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetSandbox(input.SandboxID); err != nil {
		return nil, model.NewValidationError("sandbox_id", "所属沙箱不存在")
	}
	if err := s.validateResourceLimits(input.SandboxID, input); err != nil {
		return nil, err
	}
	now := time.Now()
	t := &model.ExecutionTask{
		ID:        idgen.Hex(),
		SandboxID: input.SandboxID,
		Code:      input.Code,
		Language:  input.Language,
		Stdin:     input.Stdin,
		Status:    model.TaskPending,
		CreatedAt: now,
	}
	if err := s.store.CreateExecutionTask(t); err != nil {
		return nil, err
	}
	s.log.Infof("创建执行任务: %s", t.ID)
	return t, nil
}

func (s *Service) GetExecutionTask(id string) (*model.ExecutionTask, error) {
	return s.store.GetExecutionTask(id)
}

func (s *Service) ListExecutionTasks(filter model.ExecutionTaskFilter, page, size int) ([]*model.ExecutionTask, int, error) {
	all := s.store.ListExecutionTasks()
	matched := make([]*model.ExecutionTask, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.ExecutionTask{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateExecutionTask(id string, input model.ExecutionTask) (*model.ExecutionTask, error) {
	t, err := s.store.GetExecutionTask(id)
	if err != nil {
		return nil, err
	}
	if input.Code != "" {
		t.Code = input.Code
	}
	if input.Language != "" {
		t.Language = input.Language
	}
	if input.Stdin != "" {
		t.Stdin = input.Stdin
	}
	if input.Status != "" {
		t.Status = input.Status
	}
	if err := t.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateExecutionTask(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) DeleteExecutionTask(id string) error {
	return s.store.DeleteExecutionTask(id)
}

func (s *Service) RunTask(id string) (*model.ExecutionTask, error) {
	t, err := s.store.GetExecutionTask(id)
	if err != nil {
		return nil, err
	}
	if !model.TaskCanTransition(t.Status, model.TaskRunning) {
		return nil, model.NewValidationError("status", "状态无法从 "+t.Status+" 转为 running")
	}
	t.Status = model.TaskRunning
	now := time.Now()
	t.StartedAt = &now
	if err := s.store.UpdateExecutionTask(t); err != nil {
		return nil, err
	}
	s.createExecutionLog(t.ID, model.LogLevelInfo, "任务开始执行")
	s.log.Infof("任务开始执行: %s", t.ID)
	return t, nil
}

func (s *Service) CompleteTask(id string, stdout, stderr string, exitCode int) (*model.ExecutionTask, error) {
	t, err := s.store.GetExecutionTask(id)
	if err != nil {
		return nil, err
	}
	if !model.TaskCanTransition(t.Status, model.TaskCompleted) {
		return nil, model.NewValidationError("status", "状态无法从 "+t.Status+" 转为 completed")
	}
	now := time.Now()
	if t.StartedAt != nil {
		t.DurationMs = int(now.Sub(*t.StartedAt).Milliseconds())
	}
	t.Status = model.TaskCompleted
	t.Stdout = stdout
	t.Stderr = stderr
	t.ExitCode = exitCode
	t.FinishedAt = &now
	if err := s.store.UpdateExecutionTask(t); err != nil {
		return nil, err
	}
	s.createExecutionLog(t.ID, model.LogLevelInfo, "任务执行完成")
	s.log.Infof("任务执行完成: %s", t.ID)
	return t, nil
}

func (s *Service) FailTask(id string, stderr string) (*model.ExecutionTask, error) {
	t, err := s.store.GetExecutionTask(id)
	if err != nil {
		return nil, err
	}
	if !model.TaskCanTransition(t.Status, model.TaskFailed) {
		return nil, model.NewValidationError("status", "状态无法从 "+t.Status+" 转为 failed")
	}
	now := time.Now()
	if t.StartedAt != nil {
		t.DurationMs = int(now.Sub(*t.StartedAt).Milliseconds())
	}
	t.Status = model.TaskFailed
	t.Stderr = stderr
	t.FinishedAt = &now
	if err := s.store.UpdateExecutionTask(t); err != nil {
		return nil, err
	}
	s.createExecutionLog(t.ID, model.LogLevelError, "任务执行失败: "+stderr)
	s.log.Infof("任务执行失败: %s", t.ID)
	return t, nil
}

func (s *Service) TimeoutTask(id string) (*model.ExecutionTask, error) {
	t, err := s.store.GetExecutionTask(id)
	if err != nil {
		return nil, err
	}
	if !model.TaskCanTransition(t.Status, model.TaskTimeout) {
		return nil, model.NewValidationError("status", "状态无法从 "+t.Status+" 转为 timeout")
	}
	now := time.Now()
	if t.StartedAt != nil {
		t.DurationMs = int(now.Sub(*t.StartedAt).Milliseconds())
	}
	t.Status = model.TaskTimeout
	t.FinishedAt = &now
	if err := s.store.UpdateExecutionTask(t); err != nil {
		return nil, err
	}
	s.createExecutionLog(t.ID, model.LogLevelWarn, "任务执行超时")
	s.log.Infof("任务执行超时: %s", t.ID)
	return t, nil
}

func (s *Service) KillTask(id string) (*model.ExecutionTask, error) {
	t, err := s.store.GetExecutionTask(id)
	if err != nil {
		return nil, err
	}
	if !model.TaskCanTransition(t.Status, model.TaskKilled) {
		return nil, model.NewValidationError("status", "状态无法从 "+t.Status+" 转为 killed")
	}
	now := time.Now()
	if t.StartedAt != nil {
		t.DurationMs = int(now.Sub(*t.StartedAt).Milliseconds())
	}
	t.Status = model.TaskKilled
	t.FinishedAt = &now
	if err := s.store.UpdateExecutionTask(t); err != nil {
		return nil, err
	}
	s.createExecutionLog(t.ID, model.LogLevelWarn, "任务被强制终止")
	s.log.Infof("任务被强制终止: %s", t.ID)
	return t, nil
}

func (s *Service) createExecutionLog(taskID, level, message string) {
	log := &model.ExecutionLog{
		ID:       idgen.Hex(),
		TaskID:   taskID,
		Level:    level,
		Message:  message,
		LoggedAt: time.Now(),
	}
	_ = s.store.CreateExecutionLog(log)
}

func (s *Service) validateResourceLimits(sandboxID string, input model.ExecutionTask) error {
	_, err := s.store.GetResourceLimitBySandbox(sandboxID)
	if err != nil {
		return nil
	}
	return nil
}

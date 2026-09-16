package service

import (
	"sort"
	"time"

	"sandbox/internal/model"
	"sandbox/pkg/idgen"
)

func (s *Service) CreateSchedule(input model.Schedule) (*model.Schedule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetSandbox(input.SandboxID); err != nil {
		return nil, model.NewValidationError("sandbox_id", "所属沙箱不存在")
	}
	now := time.Now()
	sch := &model.Schedule{
		ID:        idgen.Hex(),
		SandboxID: input.SandboxID,
		Code:      input.Code,
		CronExpr:  input.CronExpr,
		Status:    model.SchedulePending,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.store.CreateSchedule(sch); err != nil {
		return nil, err
	}
	s.log.Infof("创建定时任务: %s", sch.ID)
	return sch, nil
}

func (s *Service) GetSchedule(id string) (*model.Schedule, error) {
	return s.store.GetSchedule(id)
}

func (s *Service) ListSchedules(filter model.ScheduleFilter, page, size int) ([]*model.Schedule, int, error) {
	all := s.store.ListSchedules()
	matched := make([]*model.Schedule, 0, len(all))
	for _, sch := range all {
		if filter.Match(sch) {
			matched = append(matched, sch)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Schedule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateSchedule(id string, input model.Schedule) (*model.Schedule, error) {
	sch, err := s.store.GetSchedule(id)
	if err != nil {
		return nil, err
	}
	if input.Code != "" {
		sch.Code = input.Code
	}
	if input.CronExpr != "" {
		sch.CronExpr = input.CronExpr
	}
	if input.Status != "" {
		sch.Status = input.Status
	}
	sch.UpdatedAt = time.Now()
	if err := sch.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateSchedule(sch); err != nil {
		return nil, err
	}
	return sch, nil
}

func (s *Service) DeleteSchedule(id string) error {
	return s.store.DeleteSchedule(id)
}

func (s *Service) RunSchedule(id string) (*model.Schedule, error) {
	sch, err := s.store.GetSchedule(id)
	if err != nil {
		return nil, err
	}
	if !model.ScheduleCanTransition(sch.Status, model.ScheduleRunning) {
		return nil, model.NewValidationError("status", "状态无法从 "+sch.Status+" 转为 running")
	}
	sch.Status = model.ScheduleRunning
	now := time.Now()
	sch.LastRunAt = &now
	sch.UpdatedAt = now
	if err := s.store.UpdateSchedule(sch); err != nil {
		return nil, err
	}
	s.log.Infof("定时任务开始执行: %s", sch.ID)
	return sch, nil
}

func (s *Service) PauseSchedule(id string) (*model.Schedule, error) {
	sch, err := s.store.GetSchedule(id)
	if err != nil {
		return nil, err
	}
	if !model.ScheduleCanTransition(sch.Status, model.SchedulePaused) {
		return nil, model.NewValidationError("status", "状态无法从 "+sch.Status+" 转为 paused")
	}
	sch.Status = model.SchedulePaused
	sch.UpdatedAt = time.Now()
	if err := s.store.UpdateSchedule(sch); err != nil {
		return nil, err
	}
	s.log.Infof("定时任务暂停: %s", sch.ID)
	return sch, nil
}

func (s *Service) CompleteSchedule(id string) (*model.Schedule, error) {
	sch, err := s.store.GetSchedule(id)
	if err != nil {
		return nil, err
	}
	if !model.ScheduleCanTransition(sch.Status, model.ScheduleCompleted) {
		return nil, model.NewValidationError("status", "状态无法从 "+sch.Status+" 转为 completed")
	}
	sch.Status = model.ScheduleCompleted
	sch.UpdatedAt = time.Now()
	if err := s.store.UpdateSchedule(sch); err != nil {
		return nil, err
	}
	s.log.Infof("定时任务完成: %s", sch.ID)
	return sch, nil
}

package service

import "sandbox/internal/model"

type OverviewStats struct {
	TotalTasks      int     `json:"total_tasks"`
	SuccessRate     float64 `json:"success_rate"`
	AvgDurationMs   float64 `json:"avg_duration_ms"`
	TimeoutRate     float64 `json:"timeout_rate"`
	TotalSandboxes  int     `json:"total_sandboxes"`
	TotalRuntimes   int     `json:"total_runtimes"`
	TotalSchedules  int     `json:"total_schedules"`
	TotalSubmissions int    `json:"total_submissions"`
}

type LanguageDistribution struct {
	Language string `json:"language"`
	Count    int    `json:"count"`
}

type SandboxExecutionVolume struct {
	SandboxID string `json:"sandbox_id"`
	Count     int    `json:"count"`
}

func (s *Service) GetOverviewStats() (*OverviewStats, error) {
	tasks := s.store.ListExecutionTasks()
	var completed, failed, timeout, totalDuration int
	for _, t := range tasks {
		switch t.Status {
		case model.TaskCompleted:
			completed++
			if t.DurationMs > 0 {
				totalDuration += t.DurationMs
			}
		case model.TaskFailed:
			failed++
		case model.TaskTimeout:
			timeout++
		}
	}
	total := len(tasks)
	stats := &OverviewStats{
		TotalTasks:       total,
		TotalSandboxes:   len(s.store.ListSandboxes()),
		TotalRuntimes:    len(s.store.ListRuntimes()),
		TotalSchedules:   len(s.store.ListSchedules()),
		TotalSubmissions: len(s.store.ListSubmissions()),
	}
	if total > 0 {
		stats.SuccessRate = float64(completed) / float64(total) * 100
		stats.TimeoutRate = float64(timeout) / float64(total) * 100
	}
	if completed > 0 {
		stats.AvgDurationMs = float64(totalDuration) / float64(completed)
	}
	return stats, nil
}

func (s *Service) GetLanguageDistribution() ([]*LanguageDistribution, error) {
	tasks := s.store.ListExecutionTasks()
	counts := make(map[string]int)
	for _, t := range tasks {
		counts[t.Language]++
	}
	result := make([]*LanguageDistribution, 0, len(counts))
	for lang, count := range counts {
		result = append(result, &LanguageDistribution{Language: lang, Count: count})
	}
	return result, nil
}

func (s *Service) GetSandboxExecutionVolume() ([]*SandboxExecutionVolume, error) {
	tasks := s.store.ListExecutionTasks()
	counts := make(map[string]int)
	for _, t := range tasks {
		counts[t.SandboxID]++
	}
	result := make([]*SandboxExecutionVolume, 0, len(counts))
	for sid, count := range counts {
		result = append(result, &SandboxExecutionVolume{SandboxID: sid, Count: count})
	}
	return result, nil
}

func (s *Service) GetTimeoutRate() (float64, error) {
	tasks := s.store.ListExecutionTasks()
	var timeout, total int
	for _, t := range tasks {
		total++
		if t.Status == model.TaskTimeout {
			timeout++
		}
	}
	if total == 0 {
		return 0, nil
	}
	return float64(timeout) / float64(total) * 100, nil
}

func (s *Service) GetAverageDuration() (float64, error) {
	tasks := s.store.ListExecutionTasks()
	var totalDuration, completed int
	for _, t := range tasks {
		if t.Status == model.TaskCompleted && t.DurationMs > 0 {
			totalDuration += t.DurationMs
			completed++
		}
	}
	if completed == 0 {
		return 0, nil
	}
	return float64(totalDuration) / float64(completed), nil
}

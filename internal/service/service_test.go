package service

import (
	"testing"
	"time"

	"sandbox/internal/config"
	"sandbox/internal/model"
	"sandbox/internal/store"
	"sandbox/pkg/logger"
)

func newTestService() *Service {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

func TestCreateSandbox(t *testing.T) {
	svc := newTestService()
	sb, err := svc.CreateSandbox(model.Sandbox{Name: "sb1", Language: "go", Image: "golang:1.22"})
	if err != nil {
		t.Fatalf("create sandbox: %v", err)
	}
	if sb.ID == "" {
		t.Fatalf("id empty")
	}
	_, err = svc.CreateSandbox(model.Sandbox{Name: "sb1", Language: "go", Image: "x"})
	if err == nil {
		t.Fatalf("expected conflict")
	}
}

func TestCreateSandboxValidation(t *testing.T) {
	svc := newTestService()
	_, err := svc.CreateSandbox(model.Sandbox{Name: "", Language: "go", Image: "x"})
	if err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestSandboxListPagination(t *testing.T) {
	svc := newTestService()
	for i := 0; i < 5; i++ {
		_, _ = svc.CreateSandbox(model.Sandbox{Name: "sb" + string(rune('a'+i)), Language: "go", Image: "x"})
	}
	items, total, err := svc.ListSandboxes(model.SandboxFilter{}, 1, 2)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 5 {
		t.Fatalf("total mismatch")
	}
	if len(items) != 2 {
		t.Fatalf("page size mismatch")
	}
}

func TestExecutionTaskStateMachine(t *testing.T) {
	svc := newTestService()
	sb, _ := svc.CreateSandbox(model.Sandbox{Name: "sb-task", Language: "go", Image: "x"})
	task, err := svc.CreateExecutionTask(model.ExecutionTask{SandboxID: sb.ID, Code: "code", Language: "go"})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}
	if task.Status != model.TaskPending {
		t.Fatalf("expected pending")
	}

	_, err = svc.CompleteTask(task.ID, "out", "err", 0)
	if err == nil {
		t.Fatalf("expected error: cannot complete from pending")
	}

	running, err := svc.RunTask(task.ID)
	if err != nil {
		t.Fatalf("run task: %v", err)
	}
	if running.Status != model.TaskRunning {
		t.Fatalf("expected running")
	}
	if running.StartedAt == nil {
		t.Fatalf("expected started_at")
	}

	completed, err := svc.CompleteTask(task.ID, "stdout", "stderr", 0)
	if err != nil {
		t.Fatalf("complete task: %v", err)
	}
	if completed.Status != model.TaskCompleted {
		t.Fatalf("expected completed")
	}
	if completed.FinishedAt == nil {
		t.Fatalf("expected finished_at")
	}
	if completed.DurationMs < 0 {
		t.Fatalf("expected duration >= 0")
	}

	logs, _, _ := svc.ListExecutionLogs(model.ExecutionLogFilter{TaskID: task.ID}, 1, 10)
	if len(logs) != 2 {
		t.Fatalf("expected 2 logs, got %d", len(logs))
	}
}

func TestFailTask(t *testing.T) {
	svc := newTestService()
	sb, _ := svc.CreateSandbox(model.Sandbox{Name: "sb-fail", Language: "go", Image: "x"})
	task, _ := svc.CreateExecutionTask(model.ExecutionTask{SandboxID: sb.ID, Code: "code", Language: "go"})
	svc.RunTask(task.ID)
	failed, err := svc.FailTask(task.ID, "error msg")
	if err != nil {
		t.Fatalf("fail task: %v", err)
	}
	if failed.Status != model.TaskFailed {
		t.Fatalf("expected failed")
	}
}

func TestTimeoutTask(t *testing.T) {
	svc := newTestService()
	sb, _ := svc.CreateSandbox(model.Sandbox{Name: "sb-to", Language: "go", Image: "x"})
	task, _ := svc.CreateExecutionTask(model.ExecutionTask{SandboxID: sb.ID, Code: "code", Language: "go"})
	svc.RunTask(task.ID)
	time.Sleep(10 * time.Millisecond)
	to, err := svc.TimeoutTask(task.ID)
	if err != nil {
		t.Fatalf("timeout task: %v", err)
	}
	if to.Status != model.TaskTimeout {
		t.Fatalf("expected timeout")
	}
	if to.DurationMs < 1 {
		t.Fatalf("expected positive duration")
	}
}

func TestKillTask(t *testing.T) {
	svc := newTestService()
	sb, _ := svc.CreateSandbox(model.Sandbox{Name: "sb-kill", Language: "go", Image: "x"})
	task, _ := svc.CreateExecutionTask(model.ExecutionTask{SandboxID: sb.ID, Code: "code", Language: "go"})
	svc.RunTask(task.ID)
	killed, err := svc.KillTask(task.ID)
	if err != nil {
		t.Fatalf("kill task: %v", err)
	}
	if killed.Status != model.TaskKilled {
		t.Fatalf("expected killed")
	}
}

func TestScheduleStateMachine(t *testing.T) {
	svc := newTestService()
	sb, _ := svc.CreateSandbox(model.Sandbox{Name: "sb-sch", Language: "go", Image: "x"})
	sch, err := svc.CreateSchedule(model.Schedule{SandboxID: sb.ID, Code: "code", CronExpr: "0 0 * * *"})
	if err != nil {
		t.Fatalf("create schedule: %v", err)
	}
	if sch.Status != model.SchedulePending {
		t.Fatalf("expected pending")
	}

	running, err := svc.RunSchedule(sch.ID)
	if err != nil {
		t.Fatalf("run schedule: %v", err)
	}
	if running.Status != model.ScheduleRunning {
		t.Fatalf("expected running")
	}

	paused, err := svc.PauseSchedule(sch.ID)
	if err != nil {
		t.Fatalf("pause schedule: %v", err)
	}
	if paused.Status != model.SchedulePaused {
		t.Fatalf("expected paused")
	}

	completed, err := svc.CompleteSchedule(sch.ID)
	if err != nil {
		t.Fatalf("complete schedule: %v", err)
	}
	if completed.Status != model.ScheduleCompleted {
		t.Fatalf("expected completed")
	}

	_, err = svc.RunSchedule(sch.ID)
	if err == nil {
		t.Fatalf("expected error after completed")
	}
}

func TestResourceLimitForeignKey(t *testing.T) {
	svc := newTestService()
	_, err := svc.CreateResourceLimit(model.ResourceLimit{SandboxID: "missing", CPUQuota: 1, MemoryMB: 128, DiskMB: 512, MaxProcesses: 10})
	if err == nil {
		t.Fatalf("expected validation error for missing sandbox")
	}
}

func TestRuntimeForeignKey(t *testing.T) {
	svc := newTestService()
	_, err := svc.CreateRuntime(model.Runtime{SandboxID: "missing", Name: "r", Version: "1", Command: "c", MemoryLimitMB: 128, CPULimit: 1, TimeoutMs: 1000})
	if err == nil {
		t.Fatalf("expected validation error for missing sandbox")
	}
}

func TestBatchCreateExecutionTasks(t *testing.T) {
	svc := newTestService()
	sb, _ := svc.CreateSandbox(model.Sandbox{Name: "sb-batch", Language: "go", Image: "x"})
	inputs := []model.ExecutionTask{
		{SandboxID: sb.ID, Code: "c1", Language: "go"},
		{SandboxID: sb.ID, Code: "c2", Language: "go"},
		{SandboxID: "", Code: "c3", Language: "go"},
	}
	results, err := svc.BatchCreateExecutionTasks(inputs)
	if err != nil {
		t.Fatalf("batch: %v", err)
	}
	if len(results) != 3 {
		t.Fatalf("results count mismatch")
	}
	if !results[0].Success {
		t.Fatalf("first should succeed")
	}
	if !results[1].Success {
		t.Fatalf("second should succeed")
	}
	if results[2].Success {
		t.Fatalf("third should fail")
	}
}

func TestStatsOverview(t *testing.T) {
	svc := newTestService()
	sb, _ := svc.CreateSandbox(model.Sandbox{Name: "sb-stat", Language: "go", Image: "x"})
	task, _ := svc.CreateExecutionTask(model.ExecutionTask{SandboxID: sb.ID, Code: "code", Language: "go"})
	svc.RunTask(task.ID)
	svc.CompleteTask(task.ID, "out", "err", 0)

	stats, err := svc.GetOverviewStats()
	if err != nil {
		t.Fatalf("overview: %v", err)
	}
	if stats.TotalTasks != 1 {
		t.Fatalf("total tasks mismatch")
	}
	if stats.SuccessRate != 100 {
		t.Fatalf("success rate mismatch")
	}
}

func TestExportSnapshot(t *testing.T) {
	svc := newTestService()
	_, _ = svc.CreateSandbox(model.Sandbox{Name: "sb-exp", Language: "go", Image: "x"})
	snapshot, err := svc.ExportSnapshot()
	if err != nil {
		t.Fatalf("export: %v", err)
	}
	if len(snapshot.Sandboxes) != 1 {
		t.Fatalf("sandbox count mismatch")
	}
}

package store

import (
	"testing"
	"time"

	"sandbox/internal/model"
)

func TestMemoryStoreSandbox(t *testing.T) {
	s := NewMemoryStore()
	sb := &model.Sandbox{ID: "s1", Name: "go-sandbox", Language: "go", Image: "golang:1.22", Status: model.SandboxActive, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateSandbox(sb); err != nil {
		t.Fatalf("create sandbox: %v", err)
	}
	if err := s.CreateSandbox(&model.Sandbox{ID: "s2", Name: "go-sandbox", Language: "go", Image: "x"}); err != ErrConflict {
		t.Fatalf("expected conflict, got %v", err)
	}
	got, err := s.GetSandbox("s1")
	if err != nil {
		t.Fatalf("get sandbox: %v", err)
	}
	if got.Name != "go-sandbox" {
		t.Fatalf("name mismatch")
	}
	if _, err := s.GetSandbox("missing"); err != ErrNotFound {
		t.Fatalf("expected not found")
	}
	if len(s.ListSandboxes()) != 1 {
		t.Fatalf("list count mismatch")
	}
	sb.Name = "go-sandbox-v2"
	if err := s.UpdateSandbox(sb); err != nil {
		t.Fatalf("update: %v", err)
	}
	if err := s.DeleteSandbox("s1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
	if _, err := s.GetSandbox("s1"); err != ErrNotFound {
		t.Fatalf("expected not found after delete")
	}
}

func TestMemoryStoreRuntime(t *testing.T) {
	s := NewMemoryStore()
	r := &model.Runtime{ID: "r1", SandboxID: "s1", Name: "go-runtime", Version: "1.0", Command: "go run", MemoryLimitMB: 128, CPULimit: 1.0, TimeoutMs: 5000, Status: model.RuntimeActive, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateRuntime(r); err != nil {
		t.Fatalf("create runtime: %v", err)
	}
	got, err := s.GetRuntime("r1")
	if err != nil {
		t.Fatalf("get runtime: %v", err)
	}
	if got.Name != "go-runtime" {
		t.Fatalf("name mismatch")
	}
	if _, err := s.GetRuntime("missing"); err != ErrNotFound {
		t.Fatalf("expected not found")
	}
	r.Name = "go-runtime-v2"
	if err := s.UpdateRuntime(r); err != nil {
		t.Fatalf("update: %v", err)
	}
	if len(s.ListRuntimes()) != 1 {
		t.Fatalf("list count mismatch")
	}
	if err := s.DeleteRuntime("r1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestMemoryStoreExecutionTask(t *testing.T) {
	s := NewMemoryStore()
	task := &model.ExecutionTask{ID: "t1", SandboxID: "s1", Code: "fmt.Println(1)", Language: "go", Status: model.TaskPending, CreatedAt: time.Now()}
	if err := s.CreateExecutionTask(task); err != nil {
		t.Fatalf("create task: %v", err)
	}
	got, err := s.GetExecutionTask("t1")
	if err != nil {
		t.Fatalf("get task: %v", err)
	}
	if got.Status != model.TaskPending {
		t.Fatalf("status mismatch")
	}
	if _, err := s.GetExecutionTask("missing"); err != ErrNotFound {
		t.Fatalf("expected not found")
	}
	task.Status = model.TaskRunning
	if err := s.UpdateExecutionTask(task); err != nil {
		t.Fatalf("update: %v", err)
	}
	if len(s.ListExecutionTasks()) != 1 {
		t.Fatalf("list count mismatch")
	}
	if err := s.DeleteExecutionTask("t1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestMemoryStoreSecurityPolicy(t *testing.T) {
	s := NewMemoryStore()
	p := &model.SecurityPolicy{ID: "p1", Name: "strict", MaxMemoryMB: 128, MaxCPUMs: 1000, MaxTimeoutMs: 5000, Status: model.PolicyActive, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateSecurityPolicy(p); err != nil {
		t.Fatalf("create policy: %v", err)
	}
	if err := s.CreateSecurityPolicy(&model.SecurityPolicy{ID: "p2", Name: "strict"}); err != ErrConflict {
		t.Fatalf("expected conflict")
	}
	got, err := s.GetSecurityPolicy("p1")
	if err != nil {
		t.Fatalf("get policy: %v", err)
	}
	if got.Name != "strict" {
		t.Fatalf("name mismatch")
	}
	if _, err := s.GetSecurityPolicy("missing"); err != ErrNotFound {
		t.Fatalf("expected not found")
	}
	if err := s.DeleteSecurityPolicy("p1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestMemoryStoreResourceLimit(t *testing.T) {
	s := NewMemoryStore()
	rl := &model.ResourceLimit{ID: "rl1", SandboxID: "s1", CPUQuota: 1.0, MemoryMB: 128, DiskMB: 512, MaxProcesses: 10, Status: model.ResourceLimitActive, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateResourceLimit(rl); err != nil {
		t.Fatalf("create limit: %v", err)
	}
	if err := s.CreateResourceLimit(&model.ResourceLimit{ID: "rl2", SandboxID: "s1"}); err != ErrConflict {
		t.Fatalf("expected conflict")
	}
	got, err := s.GetResourceLimit("rl1")
	if err != nil {
		t.Fatalf("get limit: %v", err)
	}
	if got.SandboxID != "s1" {
		t.Fatalf("sandbox mismatch")
	}
	bySandbox, err := s.GetResourceLimitBySandbox("s1")
	if err != nil {
		t.Fatalf("get by sandbox: %v", err)
	}
	if bySandbox.ID != "rl1" {
		t.Fatalf("by sandbox mismatch")
	}
	if _, err := s.GetResourceLimitBySandbox("missing"); err != ErrNotFound {
		t.Fatalf("expected not found")
	}
	if err := s.DeleteResourceLimit("rl1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestMemoryStoreTemplate(t *testing.T) {
	s := NewMemoryStore()
	tpl := &model.Template{ID: "tpl1", Name: "go-template", Language: "go", BaseImage: "golang:1.22", Status: model.TemplateActive, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateTemplate(tpl); err != nil {
		t.Fatalf("create template: %v", err)
	}
	if err := s.CreateTemplate(&model.Template{ID: "tpl2", Name: "go-template"}); err != ErrConflict {
		t.Fatalf("expected conflict")
	}
	got, err := s.GetTemplate("tpl1")
	if err != nil {
		t.Fatalf("get template: %v", err)
	}
	if got.Name != "go-template" {
		t.Fatalf("name mismatch")
	}
	if err := s.DeleteTemplate("tpl1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestMemoryStoreSubmission(t *testing.T) {
	s := NewMemoryStore()
	sub := &model.Submission{ID: "sub1", TaskID: "t1", Submitter: "user1", SourceCode: "code", Status: model.SubmissionPending, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateSubmission(sub); err != nil {
		t.Fatalf("create submission: %v", err)
	}
	got, err := s.GetSubmission("sub1")
	if err != nil {
		t.Fatalf("get submission: %v", err)
	}
	if got.Submitter != "user1" {
		t.Fatalf("submitter mismatch")
	}
	if err := s.DeleteSubmission("sub1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestMemoryStoreSchedule(t *testing.T) {
	s := NewMemoryStore()
	sch := &model.Schedule{ID: "sch1", SandboxID: "s1", Code: "code", CronExpr: "* * * * *", Status: model.SchedulePending, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateSchedule(sch); err != nil {
		t.Fatalf("create schedule: %v", err)
	}
	got, err := s.GetSchedule("sch1")
	if err != nil {
		t.Fatalf("get schedule: %v", err)
	}
	if got.CronExpr != "* * * * *" {
		t.Fatalf("cron mismatch")
	}
	if err := s.DeleteSchedule("sch1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestMemoryStoreExecutionLog(t *testing.T) {
	s := NewMemoryStore()
	l := &model.ExecutionLog{ID: "l1", TaskID: "t1", Level: model.LogLevelInfo, Message: "msg", LoggedAt: time.Now()}
	if err := s.CreateExecutionLog(l); err != nil {
		t.Fatalf("create log: %v", err)
	}
	got, err := s.GetExecutionLog("l1")
	if err != nil {
		t.Fatalf("get log: %v", err)
	}
	if got.Message != "msg" {
		t.Fatalf("message mismatch")
	}
	if err := s.DeleteExecutionLog("l1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

func TestMemoryStoreAuditLog(t *testing.T) {
	s := NewMemoryStore()
	a := &model.AuditLog{ID: "a1", Operator: "admin", Action: "create", TargetType: "sandbox", TargetID: "s1", Detail: "detail", CreatedAt: time.Now()}
	if err := s.CreateAuditLog(a); err != nil {
		t.Fatalf("create audit: %v", err)
	}
	got, err := s.GetAuditLog("a1")
	if err != nil {
		t.Fatalf("get audit: %v", err)
	}
	if got.Operator != "admin" {
		t.Fatalf("operator mismatch")
	}
	if err := s.DeleteAuditLog("a1"); err != nil {
		t.Fatalf("delete: %v", err)
	}
}

async function loadStats() {
  try {
    const res = await fetch('/api/stats/overview');
    const body = await res.json();
    const data = body.data || {};
    const container = document.getElementById('stats');
    const items = [
      { label: '总任务数', value: data.total_tasks || 0 },
      { label: '成功率(%)', value: (data.success_rate || 0).toFixed(2) },
      { label: '平均耗时(ms)', value: (data.avg_duration_ms || 0).toFixed(2) },
      { label: '超时率(%)', value: (data.timeout_rate || 0).toFixed(2) },
      { label: '沙箱数', value: data.total_sandboxes || 0 },
      { label: '运行时数', value: data.total_runtimes || 0 },
      { label: '定时任务数', value: data.total_schedules || 0 },
      { label: '提交数', value: data.total_submissions || 0 },
    ];
    container.innerHTML = items.map(it => `<div class="stat-card"><div class="stat-value">${it.value}</div><div class="stat-label">${it.label}</div></div>`).join('');
  } catch (e) {
    console.error('加载统计失败', e);
  }
}

async function loadTasks() {
  try {
    const res = await fetch('/api/tasks');
    const body = await res.json();
    const tbody = document.querySelector('#list tbody');
    tbody.innerHTML = '';
    const items = (body.data && body.data.items) || [];
    items.forEach(t => {
      const tr = document.createElement('tr');
      tr.innerHTML = `<td>${t.id}</td><td>${t.sandbox_id}</td><td>${t.language}</td><td>${t.status}</td><td>${t.exit_code}</td><td>${t.duration_ms}</td><td>${new Date(t.created_at).toLocaleString()}</td>`;
      tbody.appendChild(tr);
    });
  } catch (e) {
    console.error('加载任务失败', e);
  }
}

async function loadSandboxes() {
  try {
    const res = await fetch('/api/sandboxes');
    const body = await res.json();
    const tbody = document.querySelector('#sandbox-list tbody');
    tbody.innerHTML = '';
    const items = (body.data && body.data.items) || [];
    items.forEach(s => {
      const tr = document.createElement('tr');
      tr.innerHTML = `<td>${s.id}</td><td>${s.name}</td><td>${s.language}</td><td>${s.image}</td><td>${s.status}</td>`;
      tbody.appendChild(tr);
    });
  } catch (e) {
    console.error('加载沙箱失败', e);
  }
}

async function load() {
  await loadStats();
  await loadTasks();
  await loadSandboxes();
}

load();

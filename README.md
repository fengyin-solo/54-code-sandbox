# 代码沙箱管理系统

纯 Go 标准库实现的代码沙箱/隔离执行服务后端，零第三方依赖。

## 运行方式

```bash
cd origin
go run ./cmd/server
```

默认监听 `:8080`。可通过环境变量配置：

- `PORT`：服务端口（默认 8080）
- `ADDR`：完整监听地址（覆盖 PORT）
- `MAX_PAGE_SIZE`：最大分页大小（默认 100）
- `API_KEY`：API 鉴权密钥（空则不开启）
- `RATE_LIMIT`：每分钟每 IP 最大请求数（默认 100）
- `LOG_LEVEL`：日志级别 debug/info/warn/error（默认 info）

## API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | / | 前端页面 |
| POST | /api/sandboxes | 创建沙箱 |
| GET | /api/sandboxes | 沙箱列表 |
| GET | /api/sandboxes/{id} | 获取沙箱 |
| PUT | /api/sandboxes/{id} | 更新沙箱 |
| DELETE | /api/sandboxes/{id} | 删除沙箱 |
| POST | /api/runtimes | 创建运行时 |
| GET | /api/runtimes | 运行时列表 |
| GET | /api/runtimes/{id} | 获取运行时 |
| PUT | /api/runtimes/{id} | 更新运行时 |
| DELETE | /api/runtimes/{id} | 删除运行时 |
| POST | /api/tasks | 创建执行任务 |
| GET | /api/tasks | 执行任务列表 |
| GET | /api/tasks/{id} | 获取执行任务 |
| PUT | /api/tasks/{id} | 更新执行任务 |
| DELETE | /api/tasks/{id} | 删除执行任务 |
| POST | /api/tasks/{id}/run | 开始执行任务 |
| POST | /api/tasks/{id}/complete | 完成任务 |
| POST | /api/tasks/{id}/fail | 标记失败 |
| POST | /api/tasks/{id}/timeout | 标记超时 |
| POST | /api/tasks/{id}/kill | 强制终止 |
| POST | /api/tasks/batch | 批量创建任务 |
| POST | /api/policies | 创建安全策略 |
| GET | /api/policies | 安全策略列表 |
| GET | /api/policies/{id} | 获取安全策略 |
| PUT | /api/policies/{id} | 更新安全策略 |
| DELETE | /api/policies/{id} | 删除安全策略 |
| POST | /api/limits | 创建资源限制 |
| GET | /api/limits | 资源限制列表 |
| GET | /api/limits/{id} | 获取资源限制 |
| PUT | /api/limits/{id} | 更新资源限制 |
| DELETE | /api/limits/{id} | 删除资源限制 |
| POST | /api/templates | 创建模板 |
| GET | /api/templates | 模板列表 |
| GET | /api/templates/{id} | 获取模板 |
| PUT | /api/templates/{id} | 更新模板 |
| DELETE | /api/templates/{id} | 删除模板 |
| POST | /api/submissions | 创建提交 |
| GET | /api/submissions | 提交列表 |
| GET | /api/submissions/{id} | 获取提交 |
| PUT | /api/submissions/{id} | 更新提交 |
| DELETE | /api/submissions/{id} | 删除提交 |
| POST | /api/schedules | 创建定时任务 |
| GET | /api/schedules | 定时任务列表 |
| GET | /api/schedules/{id} | 获取定时任务 |
| PUT | /api/schedules/{id} | 更新定时任务 |
| DELETE | /api/schedules/{id} | 删除定时任务 |
| POST | /api/schedules/{id}/run | 开始定时任务 |
| POST | /api/schedules/{id}/pause | 暂停定时任务 |
| POST | /api/schedules/{id}/complete | 完成定时任务 |
| POST | /api/logs | 创建执行日志 |
| GET | /api/logs | 执行日志列表 |
| GET | /api/logs/{id} | 获取执行日志 |
| DELETE | /api/logs/{id} | 删除执行日志 |
| POST | /api/audits | 创建审计日志 |
| GET | /api/audits | 审计日志列表 |
| GET | /api/audits/{id} | 获取审计日志 |
| DELETE | /api/audits/{id} | 删除审计日志 |
| GET | /api/stats/overview | 统计概览 |
| GET | /api/stats/languages | 语言分布 |
| GET | /api/stats/sandbox-volume | 沙箱执行量 |
| GET | /api/stats/timeout-rate | 超时率 |
| GET | /api/stats/average-duration | 平均耗时 |
| GET | /api/export/snapshot | 全量快照导出 |

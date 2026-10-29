# 40 · 实战二：任务管理 HTTP API

## 项目目标与范围

这个服务串起路由、JSON、输入校验、错误映射、并发共享状态和优雅停机。数据仅保存在内存中，重启清空；无认证、持久化、分页或多实例一致性，适合本地学习。默认监听 `127.0.0.1:8080`，添加生产能力时按工程章节逐项设计。

| 方法与路径 | 行为 | 成功状态 |
| --- | --- | --- |
| GET /healthz | 存活检查 | 200 |
| GET /api/tasks | 按 ID 升序列出任务 | 200，空集合为 [] |
| POST /api/tasks | 创建标题长度 1—200 rune 的任务 | 201，含 Location |
| GET /api/tasks/{id} | 查询任务 | 200 |
| PATCH /api/tasks/{id} | 显式设置 done 布尔值 | 200 |
| DELETE /api/tasks/{id} | 删除任务 | 204，无 body |

## 启动与调用

```sh
go run ./cmd/taskapi
```

另开终端，POSIX shell 示例：

```sh
curl -i -H 'Content-Type: application/json' -d '{"title":"学习 Go"}' http://127.0.0.1:8080/api/tasks
curl http://127.0.0.1:8080/api/tasks
curl -X PATCH -H 'Content-Type: application/json' -d '{"done":true}' http://127.0.0.1:8080/api/tasks/1
curl -X DELETE http://127.0.0.1:8080/api/tasks/1
```

PowerShell 可用原生对象编码，避免外部程序参数引号差异：

```powershell
$body = @{ title = '学习 Go' } | ConvertTo-Json
Invoke-RestMethod -Method Post -Uri 'http://127.0.0.1:8080/api/tasks' -ContentType 'application/json; charset=utf-8' -Body ([Text.Encoding]::UTF8.GetBytes($body))
Invoke-RestMethod -Uri 'http://127.0.0.1:8080/api/tasks'
```

## 存储的并发契约

```go
--8<-- "internal/taskapi/store.go"
```

Mutex 同时保护任务 map 与自增 ID，创建时分配 ID 和写入在同一临界区。返回 Task 副本和列表快照，避免调用者绕开锁修改内部 map。Task 字段没有可变切片或 map，浅拷贝在这里足以隔离状态。

列表排序目前在锁内完成，代码简单但大数据量时会延长持锁；可以先复制再解锁排序。进一步扩展时要考虑分页与容量，而不是只换 RWMutex。并发创建测试验证 ID 唯一、数量正确，CI race 检查实际访问路径。

## HTTP 边界

```go
--8<-- "internal/taskapi/handler.go"
```

POST/PATCH 要求 application/json，允许 charset 参数。请求体最大 1 MiB，拒绝未知字段、格式错误和尾随 JSON。标题 TrimSpace 后按 rune 计数；PATCH 用 `*bool` 让 false 与未提供区分开。

缺失资源返回 404，非法 ID 或字段返回 400，媒体类型错误 415，请求体超限 413。ServeMux 自身处理的未知路径和错误方法可返回标准库默认文本响应，业务处理器返回 JSON；若产品要求全部统一，可在后续建立外层协议包装并增加测试。

传统 encoding/json 对重复 key 等仍有兼容行为，严格输入不等于完整拒绝所有 JSON 歧义。需要签名或跨解析器授权时应增加更严格契约和测试。

## 生命周期

```go
--8<-- "cmd/taskapi/main.go"
```

Server 设置头部、读写和空闲超时，收到终止信号后使用独立 5 秒预算 Shutdown。超时再 Close，避免无限等待。服务启动错误会返回主流程，后台发送错误使用缓冲 channel，防止入口提前退出后阻塞发送。

## 验证与后续练习

```sh
go test -v ./internal/taskapi
go test -race ./internal/taskapi
```

测试覆盖完整增删改查、显式 false、空列表、非法请求、体积限制、错误方法、并发存储及本地 HTTP 往返。race 需要支持的平台与 C 工具链，仓库 CI 在 Linux 执行。

练习：增加 limit/cursor 分页；引入数据库并为迁移与事务增加集成测试；添加租户授权与幂等键；为存储定义恰好满足服务需求的小接口；保持现有协议测试通过并明确哪些接口发生了变化。

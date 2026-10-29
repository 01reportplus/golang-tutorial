# 常见问题、易错点与命令速查

## 语言语义

**为什么不能给 int32 直接加 int64？** Go 不进行一般的隐式数值转换，先确定目标范围，再显式转换。字符串解析使用 strconv，不是普通类型转换。

**Go 到底有没有引用传递？** 参数按值传递。指针、切片、map 等值能关联共享数据；复制描述或地址之后修改内容，会影响共享者，但重新赋值参数本身不一定改变调用者变量。

**nil 切片能用吗？** 可 len、range、append、copy；不能索引不存在元素。nil map 可读不可写。nil channel 的发送与接收永久阻塞，nil 函数调用 panic。

**append 为什么修改了另一个切片？** 两者共享底层数组且容量足够。需要隔离时 Clone；只限制后续 append，可用完整切片表达式，但已有元素仍共享。

**range 修改元素为什么没效果？** 第二个值通常是副本。按索引修改原元素；若元素为指针，要注意共享对象的并发访问。

**接口里是 nil，为什么 err != nil？** 接口有动态类型后就不是 nil 接口。成功分支返回字面量 nil，避免返回 typed nil。

**结构体嵌入是不是继承？** 它提升字段和方法选择器，属于组合，不建立传统子类型继承。方法集和接口实现仍要按规则判断。

**为什么 []string 不能当 []any？** 切片类型与元素存储语义不同，需要逐个转换。单个 string 可赋给 any，不代表任意容器自动协变。

**可以比较两个切片吗？** 不能直接 ==；用 slices.Equal 或按业务规则比较。函数和 map 也只允许与 nil 直接比较。

**泛型方法能写吗？** 本教程 Go 1.26 基线不能给方法增加独立类型参数；Go 1.27 已有变化，且接口方法仍有限制，见[版本附录](versions.md)。

## 并发与资源

**启动 goroutine 后等一秒够吗？** 不构成完成同步。用 WaitGroup、channel 或明确协议等待；sleep 既可能不够，也可能无谓增加延迟。

**关闭 channel 会停止生产者吗？** 不会，生产者继续发送会 panic。取消用 context 或专门信号，关闭输出由确认所有发送结束的一方执行。

**select 会优先处理取消吗？** 不保证。多个可执行 case 可能被选择；取消为协作式，要把安全停止设计在各阶段。

**加了锁为什么仍竞态？** 可能存在未加同一锁的读、泄露内部指针、复制锁或锁保护范围错误。检查所有访问路径，而非只看写入。

**GC 为什么没回收大量内存？** 数据可能仍被 map、切片、闭包或 goroutine 引用；RSS 也可能来自 Go 堆外或尚未归还的内存。用 profile 区分，不凭单个指标判断。

**有 context 就会强制停止任务吗？** 不会。处理函数和底层 I/O 必须观察取消，调用方仍需等待退出。

**为什么 defer 之后文件仍越开越多？** defer 在函数返回时执行，不在循环轮次结束时执行。把单轮操作提取成函数。

## 网络、数据库与工程

**HTTP 返回 500，为什么 err 是 nil？** Client.Do 的 error 表示请求执行层面的失败，状态码需要单独判断。无论状态如何，都要处理并关闭 body。

**sql.Open 成功是否说明连上数据库？** 通常不是。使用 PingContext 验证，连接池仍会动态创建连接。

**事务能让邮件和写库一起回滚吗？** 普通 SQL 事务不能覆盖外部投递，考虑 outbox、重试与消费者幂等。

**go.sum 是锁文件吗？** 它主要记录模块内容校验，实际版本由 go.mod、依赖图和选择规则决定。

**go test 通过说明线程安全吗？** 普通测试不执行 race 检测；race 也只覆盖实际执行路径，更不证明没有业务逻辑竞争。

**GitHub Pages 能运行 Go API 吗？** Pages 托管静态文件。教程会运行在浏览器中，API 项目要另行部署。

**为什么 Pages 没更新？** 检查推送分支是否为默认分支、Pages Source 是否为 GitHub Actions、部署环境分支规则、Actions 是否失败以及浏览器缓存，详见[发布说明](publishing.md)。

## 常用命令

| 目的 | 命令 |
| --- | --- |
| 运行/构建 | `go run ./cmd/wordcount -text "go"` / `go build ./cmd/...` |
| 格式化 | `gofmt -w cmd examples internal` |
| 静态检查 | `go vet ./...` |
| 全部测试 | `go test ./...` |
| 单项测试 | `go test ./internal/taskapi -run TestCreate -v` |
| 竞态 | `go test -race -count=1 ./...` |
| 覆盖率 | `go test -cover ./...` |
| 基准 | `go test ./internal/wordcount -run '^$' -bench . -benchmem` |
| 模糊测试 | `go test ./examples/basics -run '^$' -fuzz FuzzReverse -fuzztime 10s` |
| 文档 | `go doc net/http.Client` |
| 模块 | `go mod tidy` / `go list -m all` |
| 站点 | `python -m mkdocs serve` / `python -m mkdocs build --strict` |

## 问题报告模板

记录 Go 版本、系统与架构、项目 go 指令、精确命令、最小输入、期望、实际输出和复现频率。并发问题加上退出路径和 race 报告，性能问题加上工作负载与 profile，部署问题加上 commit 与 Actions 链接。剔除 token、连接密码和私有数据后再提交 issue。

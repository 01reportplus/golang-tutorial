# 21 · context、取消与时间预算

## context 传递调用生命周期

context 主要携带取消、截止时间和少量请求范围元数据。通常作为第一个参数，命名 ctx，不传 nil；顶层可用 Background，尚未确定来源时才用 TODO。不要把 context 当成可选参数大字典，也不要用它替代显式业务参数。

```go
ctx, cancel := context.WithTimeout(parent, 2*time.Second)
defer cancel()
result, err := call(ctx)
```

即使操作提前成功，也应调用 cancel 释放相关资源。子 context 被取消不会向上取消父 context；父 context 取消会影响派生子节点。多层调用应使用已有 ctx 派生，不要在内部随意换成 Background 丢掉请求生命周期。

## 取消是协作协议

```go
func send(ctx context.Context, out chan<- int, value int) error {
    select {
    case out <- value:
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

ctx 不会强制杀死 goroutine。CPU 密集循环需要定期检查，网络和数据库应使用支持 context 的 API。锁等待本身通常不能用 context 取消，临界区设计过长会破坏整体超时预算。

ctx.Err 返回 Canceled 或 DeadlineExceeded；Cause 可以携带更具体取消原因。不要把所有 context 错误自动记录成严重服务故障，客户端主动离开与下游超时的运营含义不同。

## 分配时间预算

一个总请求预算为 2 秒时，不能让三个串行下游各独立等待 2 秒。应在同一个父截止时间下为步骤分配预算，并保留响应编码和清理余量。重试次数也要受总预算约束，不能每次重试重新获得完整时限。

超时后写操作可能已经在服务端完成，只是响应未到达。此时“客户端失败”不等于“服务端没有执行”；支付、创建订单等操作需要幂等键和结果查询。

## WithValue 的边界

适合保存请求 ID、跟踪上下文等跨 API 的元数据。key 应用私有定义类型，避免字符串碰撞；取值失败必须有合理处理。数据库连接、业务状态和必需用户参数通常更适合显式参数或依赖字段。

value 存进去以后并不会自动变成线程安全对象。若存放可变指针，多个 goroutine 仍需同步。长时间保留 context 还可能连带保留整个 value 引用链。

## 取消原因与清理动作

```go
ctx, cancel := context.WithCancelCause(parent)
defer cancel(nil)
// 某个阶段失败：
cancel(fmt.Errorf("读取任务配置: %w", err))
// 等待方观察：
<-ctx.Done()
reason := context.Cause(ctx)
_ = reason
```

Cause 提供诊断语义，Err 仍表示通用取消状态。多方同时取消时，原因由规定的先后关系确定，不应依靠其实现严格业务优先级。需要区分哪个子任务失败时，把任务 ID 或阶段信息加入错误，避免只返回一串 context canceled。

AfterFunc 可在取消后安排回调，适合连接某些旧 API 的解除阻塞动作。回调可能与主流程并发，停止回调和回调已开始之间存在竞争，状态与资源仍需同步。不会因为回调与 ctx 关联就自动成为串行清理。

## 后台任务如何处理

HTTP 请求结束后 r.Context 通常取消，因此不能简单把它交给需要长期运行的后台任务。应明确任务归属：由应用级生命周期管理，或持久化到队列。WithoutCancel 可以去除取消链，但不等于建立了可靠任务系统，还需要自身超时、日志、重启恢复与完成跟踪。

优雅停机一般先停止接收新任务，再取消或等待进行中的任务，最后关闭资源。仅调用 cancel 而不等待完成，无法确认文件、连接和后台写入已经结束。

练习：实现带超时的 HTTP 客户端，服务端故意慢响应；分别测试主动取消、截止时间和正常完成。解释为什么超时不应直接触发无限重试。

延伸阅读：[context](https://pkg.go.dev/context)、[Go Concurrency Patterns: Context](https://go.dev/blog/context)。

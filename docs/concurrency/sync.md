# 20 · Mutex、WaitGroup、Once 与原子操作

## 锁保护不变量

锁不是保护某一行代码，而是保护共享状态的一组关系。先写明哪些字段必须在同一临界区内访问，再选择 Mutex 或其他机制。多个字段共同构成余额、状态和版本时，各自原子更新不代表整体一致。

```go
type Counter struct {
    mu sync.Mutex
    n int
}

func (c *Counter) Add(delta int) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.n += delta
}

func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.n
}
```

读也要加锁，否则写受保护而读未受保护仍然有竞态。返回内部 map 或指针后，调用者可能绕过锁修改内容；需要返回快照或设计受控访问函数。

## Mutex 与 RWMutex

Mutex 不可重入，同一 goroutine 持有后再次 Lock 会死锁。不要在持锁期间调用未知回调或执行耗时网络 I/O。统一多个锁的获取顺序，减少锁顺序反转。

RWMutex 允许多个读者或一个写者，但并非总比 Mutex 快。写者等待、读临界区开销、竞争模式都会影响效果；无法简单从 RLock 升级为 Lock，应重新设计并检查状态。它也不支持无条件递归读锁。

## WaitGroup、Once 和 Cond

WaitGroup 用于等待完成，不能代替互斥。Once 保证一个初始化动作只执行一次，失败是否重试必须设计清楚；如果函数 panic，Once 也会视为已经调用。OnceValue/OnceValues 适合缓存函数结果。

Cond 用于等待共享条件，Wait 会释放锁并在返回前重新加锁。始终在循环中检查条件，因为被唤醒不代表条件仍成立。简单生产消费可用 channel，复杂共享条件可用 Cond；根据状态模型选择，不必强行统一。

## 原子操作

```go
var requests atomic.Int64
requests.Add(1)
n := requests.Load()
_ = n
```

优先使用带类型的 atomic 值，避免旧式指针 API 的位宽与对齐错误。Go 原子操作提供顺序一致性语义，但组合的“先检查再修改”仍需 CAS 循环或锁。CAS 也要考虑重试、ABA 与对象生命周期，不能凭一段看起来无锁的代码推断更快。

`atomic.Value` 适合发布只读快照，Store 的具体类型必须保持一致，不能存 nil。发布之后不能再修改快照内部共享 map；否则读者仍会竞态。

## Pool 与 Map 的适用范围

sync.Pool 用来复用临时对象，池中对象随时可能被回收，不能用作连接池或可靠缓存。归还对象前清理敏感内容，归还后不要再使用；超大缓冲区可能造成内存保留，应设容量策略。

sync.Map 适合文档描述的特定访问模式，如读多写少且 key 稳定，或不同 goroutine 主要访问不相交 key。普通有类型 map 加锁通常更容易维护业务不变量。

## 验证

运行 `go test -race ./...` 可以发现执行路径上的数据竞态，但没有报告不等于证明所有调度都安全。死锁、饥饿、错误关闭和业务一致性还需要设计审查与专门测试。

练习：实现并发计数器和快照读取；故意去掉读取锁，让 race 检测发现问题；再解释为何两个独立 atomic 不能自动维护 `available + reserved = total`。

延伸阅读：[sync](https://pkg.go.dev/sync)、[sync/atomic](https://pkg.go.dev/sync/atomic)。

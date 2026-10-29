# 19 · channel、关闭协议与 select

## channel 是带类型的通信机制

`make(chan T)` 创建无缓冲 channel，发送与接收要相互配合；`make(chan T, n)` 创建容量 n 的缓冲 channel，未满时发送可以继续。缓冲改变等待时机，不会自动解决竞争、泄漏或无限积压。

```go
package main

import "fmt"

func main() {
    values := make(chan int)
    go func() {
        defer close(values)
        for i := range 3 { values <- i }
    }()
    for n := range values { fmt.Println(n) }
}
```

`chan<- T` 只允许发送，`<-chan T` 只允许接收。函数签名使用方向表达职责，但不产生新的 channel 或复制通信数据结构。

## nil、关闭、阻塞的行为表

| 状态 | 发送 | 接收 | close |
| --- | --- | --- | --- |
| nil | 永久阻塞 | 永久阻塞 | panic |
| 打开且无对应条件 | 等待空间/接收者 | 等待数据/发送者 | 可关闭 |
| 已关闭 | panic | 先取缓冲，之后返回零值且 ok=false | panic |

关闭表示“不会再有发送”，不是取消任何一方执行的通用命令。通常由能够证明所有发送已完成的一方关闭。多个生产者应由协调者等待全部结束后关闭结果 channel，不能每个生产者各 close 一次。

接收者可以使用 `v, ok := <-ch` 区分真实零值与结束。range 会一直读取，直到 channel 关闭且缓冲耗尽；没有 close 或取消条件就可能永远等待。

## select 选择可进行的通信

```go
select {
case out <- result:
    return nil
case <-ctx.Done():
    return ctx.Err()
}
```

多个 case 同时可执行时会选择其中一个，不保证取消优先，也不保证某业务队列公平。若必须在开始下一步前避免已取消工作，可先检查 ctx.Err，但取消仍是协作式，无法彻底消除检查后的竞争窗口。

`default` 使 select 不等待，适合明确需要“尝试一次”的操作。在循环中反复 default 可能忙等。nil channel 对应 case 永不可用，可以用于动态禁用分支，但应保证剩余分支仍有退出路径。

## 超时与定时器

简单等待可用 `time.NewTimer` 并及时 Stop。高频循环反复 `time.After` 可能产生不必要分配，优先复用 timer，但 Reset 的规则要按目标版本文档处理。Go 1.23 改变了 timer channel 语义，不要混用旧版 drain 技巧与新版假设。

超时不代表底层任务已经结束。向 context 发出取消信号后，若任务忽略信号，仍可能继续占用连接或 goroutine。等待退出与通知取消是两个步骤。

## 一个可取消的生产函数

```go
func Generate(ctx context.Context, values ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, value := range values {
            select {
            case out <- value:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}
```

调用者必须读完 out，或在提前结束时取消 ctx。比如只需要第一个结果时，在外层创建 WithCancel 并 defer cancel；否则生产者可能阻塞在第二次发送。返回只读 channel 把“调用者只接收”写进类型，但无法强制调用者遵守取消协议。

需要确认生产者已经完成资源清理时，还应提供 done 或等待句柄，不能把“我调用了 cancel”视作“后台已经结束”。实际 API 可以把启动、取消与 Wait 组合为一个任务对象，或者像综合项目一样让调用函数本身负责等待。

## 扇入时避免零值刷屏

select 监听多个 channel 时，关闭的 channel 会持续立即返回零值。若只 `continue`，可能形成高速循环。读取到 `ok=false` 后可将该 channel 变量设为 nil，禁用对应分支；所有输入关闭后退出。输入数量较多时则常用每个输入一个转发 worker，加协调者关闭汇总输出。

## 数据发送也可能共享

发送结构体复制结构体值；其中含 map、切片或指针，接收方仍可能和发送方共享底层数据。可以约定发送后转移所有权，或者复制后发送，或者通过锁保护。channel 不会自动深拷贝对象。

练习：用两个生产者合并结果，只有一个协调者关闭 out；让消费者只读一个结果就取消，证明生产者全部能够退出。

延伸阅读：[并发流水线与取消](https://go.dev/blog/pipelines)、[time](https://pkg.go.dev/time)。

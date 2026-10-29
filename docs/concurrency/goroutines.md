# 18 · goroutine、调度与生命周期

## 并发与并行不同

并发是同时组织多个进行中的任务，并行是同一时刻在多个计算资源上执行。单核也能并发；只有足够资源和可并行工作时才会得到并行加速。网络任务常在等待 I/O，所以并发可以提高吞吐，但也会增加连接数、内存和下游压力。

`go f()` 启动一个 goroutine，调用参数在启动语句执行时求值。主 goroutine 返回后进程结束，不会自动等待其他任务。因此，`time.Sleep` 只能制造偶然等待，不能表达完成关系。

```go
package main

import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup
    for i := range 3 {
        wg.Go(func() { fmt.Println(i) }) // Go 1.25+
    }
    wg.Wait()
}
```

输出顺序不固定，但三项工作都结束后 main 才返回。WaitGroup.Go 的函数不能 panic；需要处理失败时返回错误到受控通道，或使用更适合错误传播的结构。

## 经典 Add/Done 形式

```go
wg.Add(1)
go func() {
    defer wg.Done()
    work()
}()
wg.Wait()
```

Add 应在启动 goroutine 之前执行，否则 Wait 可能先观察到计数为零。WaitGroup 用后不能复制；开始新一轮等待前，要确保上一轮 Wait 已结束。它只负责等待完成，不限制并发、不传播错误，也不会取消任务。

## 调度器的理解层次

常见 G/M/P 模型中，G 代表 goroutine，M 代表操作系统线程，P 表示执行 Go 代码所需的调度资源。goroutine 可以在不同线程间调度；阻塞系统调用与网络轮询有专门协作机制。该模型帮助理解调度，不是普通代码可以依赖的执行顺序保证。

`GOMAXPROCS` 主要约束同时执行 Go 代码的调度资源数量，不等于 goroutine 上限，也不等于操作系统线程数。容器中的默认行为会随版本与资源限制变化，应按当前运行环境确认。

## 每个 goroutine 都需要退出路径

启动前写下四件事：谁拥有它、正常何时结束、失败如何报告、调用者取消后它如何退出。goroutine 泄漏会保留栈、局部变量和引用对象，GC 无法替你回收一个仍阻塞并可达的执行任务。

常见泄漏来源包括：向无人接收的 channel 发送；等待永不关闭的 channel；未设置截止时间的 I/O；后台 ticker 永远运行；调用者提前返回但生产者继续阻塞。

## goroutine 并非零成本

它通常比操作系统线程轻量，但大量任务仍占用栈、调度开销和对象。不要对百万条数据无界启动百万个 goroutine。采用工作池或信号量限制活跃任务，根据下游容量和实际测量调整并发度。

panic 只在当前 goroutine 内展开。HTTP 中间件的 recover 不能捕获它另行启动的后台任务 panic；恢复点必须与对应任务的失败策略一致。

练习：把示例里的等待去掉，观察程序可能提前退出；然后为一个后台循环画出取消、退出和等待完成的流程。

延伸阅读：[sync.WaitGroup](https://pkg.go.dev/sync#WaitGroup)、[调度器介绍](https://go.dev/src/runtime/HACKING.md)。

# 09 · 函数、多返回值、闭包与 defer

## 函数也是值

函数由参数、返回值和函数体组成，参数永远按值传递。函数类型由参数类型和返回值类型决定，参数名称不影响类型一致性。函数可以赋值给变量、作为参数传递，也可以作为返回值，用于排序、回调和依赖注入。

```go
func MapInts(xs []int, f func(int) int) []int {
    out := make([]int, len(xs))
    for i, x := range xs { out[i] = f(x) }
    return out
}
```

多返回值常用于 `value, error` 和 `value, ok`。命名返回值会在函数开始时取得零值，适合简短、含义清楚的函数；较长函数中裸 return 容易让人忘记返回什么，显式返回更直观。

## 可变参数

```go
func Sum(xs ...int) int {
    total := 0
    for _, x := range xs { total += x }
    return total
}
// nums := []int{1, 2, 3}
// total := Sum(nums...)
```

可变参数在函数内部表现为切片。传入 `nums...` 不会自动深拷贝元素，修改参数元素可能影响调用者。`[]string` 不能直接展开为 `...any`，因为两种切片类型不同，必要时逐个转换。

## 闭包捕获变量

```go
func Counter() func() int {
    n := 0
    return func() int {
        n++
        return n
    }
}
```

同一个 Counter 返回的函数共享其 n，多次调用 Counter 则产生不同状态。闭包捕获的是变量关联的状态，不是你直觉中某一时刻的快照。将同一个计数闭包交给多个 goroutine 会产生竞态，需要同步。

递归可以清楚表达树和分治，但 Go 不保证尾调用优化。输入深度不可控时，应考虑显式栈或深度限制，避免把理论简洁变成资源风险。

## defer 在什么时候做什么

defer 语句执行时立即求值被调用函数及其参数，调用本身延迟到当前函数返回前；多个 defer 后进先出。它绑定函数返回，不绑定 for 或 if 块结束。

```go
package main

import "fmt"

func main() {
    n := 1
    defer fmt.Println("参数快照", n)
    defer func() { fmt.Println("读取变量", n) }()
    n = 2
}
```

输出先是“读取变量 2”，再是“参数快照 1”。defer 中的闭包可以修改命名返回值，常用于处理关闭或提交错误，但应避免隐藏主流程。

## 资源清理的边界

```go
f, err := os.Open(path)
if err != nil { return err }
defer f.Close()
```

只在打开成功后 defer。读取场景通常优先返回读取错误；写文件时 Close 也可能报告刷新失败，需要保留它。在处理海量文件的循环中，把每次操作提取为小函数，让文件及时关闭；在大函数中对每次循环 defer 会一直保留资源到函数结束。

panic 展开栈时会执行当前 goroutine 的 defer；`os.Exit` 不会。defer 不是后台线程，也不会在语句所在处立即清理。

练习：预测三层 defer 的输出顺序；实现一个复制文件函数，同时考虑打开、复制和关闭失败，并写出你选择保留哪种错误的规则。

延伸阅读：[Defer、Panic 与 Recover](https://go.dev/blog/defer-panic-and-recover)。

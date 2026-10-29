# 17 · 泛型、类型集合与迭代器

## 从重复算法提取类型参数

泛型适合“结构相同，只是元素类型变化”的算法和数据结构。接口适合按行为替换不同实现。先写出两三个具体用例，再提取共同操作，通常比预先设计一套抽象层更准确。

```go
package main

import "fmt"

type Integer interface { ~int | ~int32 | ~int64 }

func Sum[T Integer](xs []T) T {
    var total T
    for _, x := range xs { total += x }
    return total
}

type Score int
func main() { fmt.Println(Sum([]Score{1, 2, 3})) }
```

`T` 是类型参数，Integer 是约束。`~int` 包含底层类型为 int 的定义类型，`int` 只列出该具体类型。联合类型元素描述允许的类型集合，不表示在运行时把值装成一个联合对象。

## 泛型类型与零值

```go
type Stack[T any] struct { values []T }

func (s *Stack[T]) Push(v T) { s.values = append(s.values, v) }
func (s *Stack[T]) Pop() (T, bool) {
    if len(s.values) == 0 { var zero T; return zero, false }
    i := len(s.values)-1
    v := s.values[i]
    var zero T
    s.values[i] = zero // 释放尾部引用
    s.values = s.values[:i]
    return v, true
}
```

泛型接收者声明中的 T 对应泛型类型已有参数。在 Go 1.26 基线里，方法不能另行引入自己的类型参数；可以使用独立泛型函数。Go 1.27 的变化见版本附录，避免把这个基线限制误记成永久规则。

## comparable 和推断

`comparable` 允许 `==`、`!=` 和作为 map key，但不承诺业务等价，也不意味着 `map[any]...` 永不因动态 key 而 panic。类型推断能省略很多调用处参数，失败时可以显式写 `Sum[int](...)`。

约束内只允许对类型集合所有成员都有效的操作。不要因为当前调用都传入 int，就在 `T any` 的函数里写 `a + b`。约束越窄，操作越多，但可接受类型越少。

## 使用标准泛型库

`slices` 提供排序、二分查找、查找和删除；`maps` 提供复制、键值迭代等；`cmp` 提供有序类型约束与比较工具。标准库已经提供的基础操作通常无需自己重复实现。对含浮点或自定义对象的排序，明确 NaN、相等与稳定性要求。

## range over function

Go 1.23 起，for range 可以遍历符合签名的迭代函数。`iter.Seq[T]` 相当于接收 `yield func(T) bool` 的函数，Seq2 提供两个值。yield 返回 false 表示调用方提前停止，生产者必须停止继续 yield。

```go
func Count(n int) iter.Seq[int] {
    return func(yield func(int) bool) {
        for i := 0; i < n; i++ {
            if !yield(i) { return }
        }
    }
}
// for n := range Count(3) { fmt.Println(n) }
```

迭代器不要求创建 goroutine。`iter.Pull` 把推式转换为拉式时，应按文档及时调用 stop 释放资源。若迭代底层数据库或文件，提前退出时仍需可靠关闭。

## 练习

实现保持首次出现顺序的 `Unique[T comparable]`；给 Stack 编写包含指针元素的测试；写一个可提前 break 的行迭代器，并解释资源关闭放在哪里。

延伸阅读：[泛型教程](https://go.dev/doc/tutorial/generics)、[iter 包](https://pkg.go.dev/iter)、[slices 包](https://pkg.go.dev/slices)。

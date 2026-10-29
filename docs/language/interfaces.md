# 16 · 接口、方法集与动态类型

## 接口表达需要的行为

一个类型只要具有接口要求的方法，就实现该接口，不需要显式声明。接口通常定义在使用方，用来表达最少依赖；提供方返回具体类型可以保留更清楚的能力和演进空间，但这也不是不可违反的法则。

```go
type Sender interface {
    Send(context.Context, string) error
}
```

接口越大，实现与测试成本越高。不要为了“以后可替换”在每个结构体前都创建一模一样的接口；先找到真正需要替换的边界，例如网络、时间或存储。

## 方法集决定实现关系

T 的方法集中包含以 T 为接收者的方法；`*T` 的方法集中还包含以 `*T` 为接收者的方法。如果只有 `func (*Buffer) Write(...)`，则 `*Buffer` 可实现 Writer，Buffer 不一定实现。

```go
var _ io.Writer = (*bytes.Buffer)(nil)
```

这条编译期断言不创建实际对象，表达接口实现预期。普通调用中 `b.Write(...)` 可能因为 b 可寻址而自动取地址，但把 b 放入接口时不会通过这种规则弥补方法集。

## 接口的 nil 陷阱

```go
package main

import "fmt"

type Problem struct{}
func (*Problem) Error() string { return "problem" }

func main() {
    var p *Problem = nil
    var err error = p
    fmt.Println(p == nil, err == nil) // true false
}
```

可以把接口值理解为动态类型与动态值的组合。只有二者都未设置时接口才是 nil。上例动态类型是 `*Problem`，所以 err 不等于 nil。返回 error 时，成功分支应直接 `return nil`，不要返回一个类型确定但为 nil 的错误指针。

这是一种理解模型，不意味着普通程序应依赖接口内部内存布局；内部表示是运行时实现细节。

## 类型断言与类型 switch

```go
if s, ok := value.(fmt.Stringer); ok {
    fmt.Println(s.String())
}
switch v := value.(type) {
case string:
    fmt.Println(len(v))
case int:
    fmt.Println(v + 1)
case nil:
    fmt.Println("没有动态值")
}
```

单返回值断言失败会 panic，双返回值形式更适合不确定输入。错误链应使用 errors.As，直接断言只检查当前最外层动态值。

## 接口如何让测试与组合更简单

```go
func WriteReport(w io.Writer, title string) error {
    _, err := fmt.Fprintf(w, "report: %s\n", title)
    return err
}
```

文件、标准输出和 bytes.Buffer 都能满足这个函数，而函数不需要知道它们的类型。测试传入 Buffer 后比较内容；错误测试传入一个总是返回错误的 Writer，验证函数没有吞掉写入失败。

这不是要求调用方把所有依赖都接口化。只用 Read 方法时接收 io.Reader，比要求 io.ReadWriteCloser 更能说明真实需求；确实需要 Close 时，还要明确资源由谁关闭，不能因为参数类型有 Close 就擅自接管调用方资源。

接口可作为可选能力探测，例如断言某个对象是否支持 Flush，但要注意包装层是否仍暴露该能力。HTTP 中间件包装 ResponseWriter 时若丢掉特殊接口，可能改变流式输出或连接升级行为。

## 函数与接口的选择

只有一个操作且无需多种关联方法时，函数类型通常更轻：`type Clock func() time.Time`。需要一组有共同状态的方法时，小接口更合适。两者都可以用于测试替换，选择取决于表达力，不取决于“接口一定更高级”。

## any 与约束接口

`any` 是 `interface{}` 的别名，允许保存任意类型，但操作前仍需知道类型。它不是“跳过类型系统”的许可；过度使用 map[string]any 会把很多错误推迟到运行时。

含类型集合元素的接口，如 `interface{ ~int | ~int64 }`，用于泛型约束，不能像普通值接口一样声明变量。接口多态处理运行时行为，泛型表达编译期可接受类型与操作，两者可以协作。

## 接口比较

接口与 nil 比较通常安全；两个接口互相比较时，如果动态类型相同但不可比较，会 panic。不要用接口相等当作通用深度比较。定义领域对象身份时，选择稳定 ID 往往比比较任意动态值更清楚。

练习：写一个含 `Write` 指针接收者的方法，分别把 T 和 `*T` 赋给 io.Writer；再为包装错误比较类型断言和 errors.As。

延伸阅读：[接口与反射基础](https://go.dev/blog/laws-of-reflection)、[io 接口](https://pkg.go.dev/io)。

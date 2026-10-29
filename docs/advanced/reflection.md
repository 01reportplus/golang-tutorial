# 35 · 反射、unsafe 与运行时边界

## 什么时候需要反射

当代码在编译时不知道具体类型，又要检查类型、字段或方法时，可以使用 reflect，例如序列化、表单绑定和配置工具。普通业务分支优先使用结构体、接口和泛型，它们更容易被编译器检查，也更易重构。

Type 描述类型，Value 描述一个带类型的值。reflect.TypeOf(nil) 返回 nil，reflect.ValueOf(nil) 得到无效 Value；调用具体操作前先理解 IsValid、Kind 和 IsNil 的适用范围。

```go
package main

import (
    "fmt"
    "reflect"
)

type Config struct { Port int `json:"port"` }

func main() {
    c := Config{Port: 8080}
    v := reflect.ValueOf(&c).Elem()
    field := v.FieldByName("Port")
    if field.IsValid() && field.CanSet() && field.Kind() == reflect.Int {
        field.SetInt(9090)
    }
    t := reflect.TypeOf(c)
    f, _ := t.FieldByName("Port")
    fmt.Println(c.Port, f.Tag.Get("json"))
}
```

## Kind 不等于具体类型

定义类型 UserID 的 Kind 可以是 Int64，但它的 Type 不等于 int64。反射转换要检查 ConvertibleTo，赋值要检查 AssignableTo。不能仅凭 Kind 相同就假定所有语义一致。

指针 Value 的 Elem 才能得到被指向值。CanAddr 表示能取地址，CanSet 表示能通过反射设置；二者相关但不等价。未导出字段受到访问限制，不应通过 unsafe 绕过封装以完成日常业务。

IsNil 只适用于 channel、函数、接口、map、指针、切片等支持 nil 的种类，对普通 int 调用会 panic。通用工具应逐个验证前置条件，避免把反射错误变成远离调用点的 panic。

## 构造与调用

reflect.New 创建类型指针，MakeSlice/MakeMap 创建容器，Call 动态调用函数。动态调用失去普通编译期参数检查，还增加分配与运行开销。缓存 Type 分析可减少重复工作，但缓存也要有并发与容量策略。

如果泛型可以表达所需操作，就不一定需要反射；如果操作依赖运行时字段名，泛型也不会自动提供字段遍历。根据信息何时可知来选择机制。

## unsafe 的作用与风险

unsafe 提供 Sizeof、Alignof、Offsetof 和指针相关操作，用于底层互操作或经测量必要的优化。它绕过部分类型安全，错误可能涉及越界、对象存活、对齐和 GC。

```go
type Sample struct { A byte; B int64 }
// unsafe.Sizeof(Sample{}) 可能包含 padding
// unsafe.Offsetof(Sample{}.B) 表示字段偏移，不是字段序号
```

结构体内存布局受字段类型与平台影响，不能直接作为磁盘或网络协议。两个零大小变量可能具有相同地址，不应把任意地址唯一性当作业务身份。

## uintptr 不是保活指针

uintptr 是整数，GC 不会因为整数看起来像地址就保留对象。把指针转成整数、经过若干语句后再转回，可能违反规则。runtime.KeepAlive 只解决特定存活需求，不会让错误指针运算、越界或并发访问合法化。

unsafe.String、StringData、Slice、SliceData 等 API 有严格前提：内存有效期、边界、不可变性都由调用者负责。把可变字节零拷贝变成 string 后继续修改，可能破坏 map key、缓存和并发假设。

## 检查工具与边界

go vet、race 和 checkptr 能发现部分问题，不能证明所有 unsafe 正确。需要建立平台测试、明确所有权、记录性能证据，并保留普通安全实现作对照。依赖 runtime 私有符号或内存布局会增加版本升级风险。

练习：编写只读取导出字段 tag 的反射函数，测试 nil、非结构体、结构体指针和未导出字段；解释为何一个零拷贝优化可能在保存 map key 时造成严重问题。

延伸阅读：[reflect](https://pkg.go.dev/reflect)、[unsafe](https://pkg.go.dev/unsafe)、[反射三法则](https://go.dev/blog/laws-of-reflection)。

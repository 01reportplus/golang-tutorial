# 语言规范覆盖索引

这份索引把规范主题映射到教程，用于查漏和复习。正文解释常用规则与工程含义；精确语法、极端类型推断和实现限制，以所用语言版本的[官方规范](https://go.dev/ref/spec)为准。不要把章节覆盖等同于逐条复述规范的全部例外。

## 源码、声明和类型

| 规范主题 | 教程入口 | 自检问题 |
| --- | --- | --- |
| 源码编码、字符、注释、标识符 | [语法](../language/syntax.md) | Unicode 字母与大小写导出规则是什么？ |
| token、关键字、运算符、分号 | [语法](../language/syntax.md)、[常量](../language/constants.md) | 为什么左花括号通常不能另起一行？ |
| 整数、浮点、虚数、rune、string 字面量 | [类型](../language/types.md)、[字符串](../language/strings.md) | 单引号、双引号、反引号分别表示什么？ |
| 常量、类型化/无类型、iota | [常量](../language/constants.md) | 常量为何可赋给多种数值类型？ |
| 变量、零值、作用域、空白标识符 | [语法](../language/syntax.md) | := 在什么时候遮蔽外部变量？ |
| 布尔、数值、字符串类型 | [类型](../language/types.md) | int 在哪些平台上可能有不同宽度？ |
| 数组、切片 | [数组与切片](../language/slices.md) | 长度是否参与类型，容量是否参与类型？ |
| map、可比较 key | [map](../language/maps.md) | map[any] 的动态 key 为什么可能 panic？ |
| 结构体、字段、tag、嵌入 | [结构体](../language/structs.md) | 嵌入产生了什么，未产生什么？ |
| 指针、地址与间接访问 | [指针](../language/pointers.md) | 为什么返回局部变量地址安全？ |
| 函数类型、参数、可变参数、结果 | [函数](../language/functions.md) | []string 为什么不能直接展开为 ...any？ |
| channel 类型与方向 | [channel](../concurrency/channels.md) | 单向 channel 如何表达所有权？ |
| 接口、方法集、类型集合 | [接口](../language/interfaces.md)、[泛型](../language/generics.md) | 普通接口与约束接口的区别是什么？ |
| 定义类型、底层类型、别名 | [类型](../language/types.md) | type A B 与 type A = B 的差别？ |
| 泛型类型、实例化、推断、约束 | [泛型](../language/generics.md) | ~T 扩展了哪些允许类型？ |

## 表达式、语句与运行

| 规范主题 | 教程入口 | 自检问题 |
| --- | --- | --- |
| 类型一致性、可赋值性、可转换性 | [类型](../language/types.md)、[反射](../advanced/reflection.md) | 语法允许转换是否代表业务安全？ |
| 字面量、复合字面量、函数字面量 | [结构体](../language/structs.md)、[函数](../language/functions.md) | 闭包捕获的是什么？ |
| selector、索引、切片与完整切片表达式 | [切片](../language/slices.md)、[结构体](../language/structs.md) | 三索引切片是否深拷贝？ |
| 类型断言、类型 switch | [接口](../language/interfaces.md) | 断言失败何时 panic，何时返回 ok？ |
| 函数调用、方法值、方法表达式 | [函数](../language/functions.md)、[结构体](../language/structs.md) | 方法值何时绑定接收者？ |
| 算术、位、比较、逻辑运算 | [常量](../language/constants.md) | 短路与按位操作有何区别？ |
| 接收、地址、间接访问与转换 | [channel](../concurrency/channels.md)、[指针](../language/pointers.md) | 接收关闭 channel 返回什么？ |
| 求值顺序、赋值与多重赋值 | [常量](../language/constants.md)、[语法](../language/syntax.md) | 为什么应避免依赖复杂副作用顺序？ |
| if、switch、for、range | [控制流](../language/control-flow.md) | range 的第二个值是否引用原元素？ |
| break、continue、goto、标签 | [控制流](../language/control-flow.md) | 跳转受哪些作用域限制？ |
| go、select、发送与关闭 | [并发入门](../concurrency/goroutines.md)、[channel](../concurrency/channels.md) | select 是否保证取消优先？ |
| return、defer、panic、recover | [函数](../language/functions.md)、[错误](../language/errors.md) | recover 为什么只能恢复当前 goroutine？ |
| 包、导入、初始化与程序执行 | [语法](../language/syntax.md)、[Modules](../engineering/modules.md) | init 为什么不适合复杂网络初始化？ |
| 运行时 panic、系统考虑、大小与对齐 | [错误](../language/errors.md)、[反射](../advanced/reflection.md) | 内存布局为什么不是跨平台协议？ |
| 语言版本与兼容 | [版本差异](versions.md) | 工具链版本和 module 语言版本是否相同？ |

## 全部传统关键字速记

```text
break default func interface select
case defer go map struct
chan else goto package switch
const fallthrough if range type
continue for import return var
```

`any`、`comparable`、`nil`、`true`、`false`、`iota` 及内置函数名是预声明标识符，不是关键字。语法上可以遮蔽部分预声明名字，但会降低可读性，并可能使原本的内置操作无法使用。

## 内置函数检查表

| 操作 | 关键点 |
| --- | --- |
| append、copy | 接收 append 返回值；copy 不改变长度 |
| len、cap | 长度与容量不同；对类型的适用范围不同 |
| make、new | make 初始化三种容器；new 返回指针，1.26 扩展表达式形式 |
| delete、clear | 删除 map 项；清空 map 或清零切片元素 |
| close | 发送协议的结束，通常由发送方协调关闭 |
| min、max | 有序类型，多参数，不直接接受切片展开 |
| complex、real、imag | 构造和拆分复数，注意对应精度 |
| panic、recover | 当前 goroutine 的异常展开与边界恢复 |
| print、println | 不属于稳定格式输出工具，应用输出用 fmt/log |

语言规则之外，另读[内存模型](../concurrency/memory-model.md)、[模块参考](../engineering/modules.md)和当前包文档。规范并不定义数据库隔离级别、HTTP 产品契约或全部 runtime 实现。

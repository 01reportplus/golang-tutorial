# 08 · 条件、循环、switch 与 range

## if 与提前返回

Go 的 if 条件必须是 bool，不能把非零整数当作 true。条件前可有简单初始化语句，其变量作用域包含对应 else 分支。处理错误时优先用提前返回减少嵌套，让正常流程保持清楚。

```go
if n, err := strconv.Atoi(input); err != nil {
    return fmt.Errorf("解析数量: %w", err)
} else if n < 0 {
    return errors.New("数量不能小于零")
}
```

若后续仍需使用 n，就在 if 之前声明，而不是为了缩短代码把重要状态藏进一个小作用域。

## for 是唯一的循环关键字

```go
for i := 0; i < 3; i++ { /* 计数循环 */ }
for ready { /* 类似 while */ }
for { /* 持续循环，必须设计退出条件 */ }
```

`break` 跳出相应循环或 switch/select；`continue` 进入下一轮循环。多层嵌套可以使用标签，但提取函数通常更容易读。`goto` 受作用域限制，不能跳过变量声明进入无效状态，日常业务代码很少需要它。

## switch 不默认穿透

```go
switch {
case score >= 90:
    grade = "A"
case score >= 60:
    grade = "B"
default:
    grade = "C"
}
```

Go 的普通 switch 匹配后自动结束，不用在每个 case 写 break。一个 case 可以列多个值。`fallthrough` 无条件进入下一个 case 的语句部分，不会重新检查条件，只在确实需要这种行为时使用。类型 switch 用来区分接口的动态类型，在接口章节展开。

## range 根据对象产生不同值

| 对象 | 第一个值 | 第二个值 | 注意 |
| --- | --- | --- | --- |
| 数组/切片 | 索引 | 元素副本 | 修改副本不修改原元素 |
| string | 字节索引 | rune | 非法 UTF-8 需要额外考虑 |
| map | key | value | 次序不保证 |
| channel | 接收到的值 | 无 | 关闭并取尽后退出 |
| 整数 | 0 到 n-1 | 无 | Go 1.22+；n≤0 不执行 |
| 迭代函数 | yield 提供的值 | 可选 | Go 1.23+ |

切片循环时如果要修改原元素，应使用 `s[i] = ...`，或者元素本身是指针并明确共享关系。遍历 map 的同时增加元素，新增项是否被访问没有固定保证，不要依赖这种行为。

## 循环变量的版本变化

对采用 Go 1.22+ 语言语义且使用 `:=` 声明的循环变量，每次迭代有独立变量，解决了许多闭包捕获同一循环变量的问题。但是 `for _, v = range values` 使用预先声明变量，仍然是重复赋值；同时，独立的变量并不能消除其指向对象的共享。

```go
for _, item := range items {
    item.Count++ // items 为 []Item 时，只改了副本
}
for i := range items {
    items[i].Count++
}
```

## 避免无意的忙等

无限 for 中若没有阻塞、休眠或实际有界工作，可能消耗整个 CPU 核。尤其在 select 中添加 default 后，原本等待事件的循环会立刻变成反复轮询。需要定期任务时用 timer/ticker 或明确调度机制。

练习：分别 range `"中文Go"`、`[]byte("中文Go")` 和 `[]rune("中文Go")`，比较索引与元素；解释为什么结果数量不同。

延伸阅读：[Go 1.22 循环变量说明](https://go.dev/blog/loopvar-preview)、[迭代器设计](https://go.dev/blog/range-functions)。

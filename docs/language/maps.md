# 14 · map、集合与可比较性

## map 的基本操作

map 把可比较的 key 映射到 value。读取不存在的 key 返回 value 零值，因此需要第二个返回值判断“缺失”与“存在且值为零”的区别。

```go
scores := make(map[string]int)
scores["Ada"] = 0
v, ok := scores["Ada"] // 0, true
_, missing := scores["Lin"] // false
delete(scores, "Ada")
clear(scores)
_, _ = v, ok
_ = missing
```

nil map 可以读取、遍历、查询长度和删除，但不能写入；写入会 panic。map 字面量和 make 都可以创建可写 map，容量提示不是长度，也不是固定容量上限。

## key 必须可比较

字符串、整数、指针，以及元素或字段均可比较的数组和结构体可作为 key。切片不能作为 key。组合 key 常用结构体：

```go
type Pair struct { TenantID, UserID int64 }
cache := map[Pair]string{}
cache[Pair{TenantID: 1, UserID: 2}] = "Ada"
```

不要通过 `tenant + ":" + user` 随意拼接复杂 key，否则分隔符、转义与编码规则容易出错。`map[any]V` 虽然语法可用，但放入动态类型不可比较的 key 会在运行时 panic。

## map 元素不可取地址

`&m[key]` 不合法。若 value 是结构体，不能直接 `m[key].Count++`，通常先取副本、修改后写回，或者明确使用指针值。

```go
entry := records[id]
entry.Count++
records[id] = entry
```

`map[K]*V` 简化对象原地修改，但引入 nil 指针、别名共享与并发保护问题。选择值还是指针取决于数据大小、身份语义和修改模式，不是固定性能口诀。

## 遍历与稳定输出

map 遍历次序不保证。同一程序多次遍历也不应期待固定次序。输出报告、测试快照或生成签名数据时，先收集 key，再排序后访问。Go 1.23+ 可以组合 `maps.Keys` 和 `slices.Sorted`。

删除未遍历到的 key 后，它不会再被产生；遍历中新增 key 是否被访问不固定。为了可读性，尽量把遍历和结构性修改分成清楚的阶段。

## 用 map 实现集合

```go
seen := map[string]struct{}{}
seen["go"] = struct{}{}
_, exists := seen["go"]
_ = exists
```

空结构体表达“只关心存在性”。`map[string]bool` 也可以，但要避免把 false 与不存在混淆。集合交集、并集、去重都可用 map 实现；需要稳定顺序时额外维护切片。

## 并发访问

多个 goroutine 仅进行读取且初始化已经正确发布时可以并发读；只要存在并发写或读写，就要按共享状态契约同步。runtime 有时会检测到并发 map 访问并终止程序，但不能把这种检测当成完整 race 检测。

常见方案是 Mutex 保护普通 map，或者让单个 goroutine 拥有 map。`sync.Map` 为特定读写模式设计，通常不应替代所有有类型 map。缓存还需要过期、容量与一致性策略，线程安全只解决其中一部分。

练习：实现保持输入首次出现顺序的去重函数；为“同一用户不同租户”设计一个不会冲突的组合 key；测试缺失值与零值。

延伸阅读：[Go maps](https://go.dev/blog/maps)、[maps 包](https://pkg.go.dev/maps)。

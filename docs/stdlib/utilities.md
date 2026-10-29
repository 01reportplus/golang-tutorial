# 26 · 时间、正则、排序、数学与常用容器

## 时间点与时长

time.Time 表示时间点，time.Duration 表示纳秒为单位的时长。Duration 的数值若没有单位，很容易错把秒当纳秒。配置中可使用 `time.ParseDuration("250ms")`。

```go
now := time.Now()
elapsed := time.Since(now)
deadline := now.Add(2*time.Second)
_ = elapsed
_ = deadline
```

time.Now 得到的值通常携带单调时钟读数，适合测量耗时；序列化或某些变换会移除这部分信息。跨机器传输的时间只能表达墙上时钟，不应据此获得严格的全局事件顺序。

## 格式、时区与日历运算

Go 使用参考时间 `2006-01-02 15:04:05` 定义 layout，不使用 `%Y-%m-%d`。优先采用 RFC3339 等明确格式。Parse 和 ParseInLocation 的时区默认行为不同；业务日期应明确是哪个地区的日期。

```go
loc, err := time.LoadLocation("Asia/Shanghai")
if err != nil { return err }
t, err := time.ParseInLocation("2006-01-02", "2026-10-29", loc)
if err != nil { return err }
fmt.Println(t.Format(time.RFC3339))
```

“一天”不总是 24 小时，存在夏令时的地区尤其如此。按日历推进用 AddDate，按精确时长用 Add。比较时间点通常用 Equal、Before、After，而不是 `==`，后者还涉及位置和内部表示。

容器精简镜像可能没有时区数据库，必要时导入 time/tzdata 或安装数据。定时任务要明确时区、错过执行的处理、重复执行与停机恢复。

## Timer 与 Ticker

Timer 表达一次触发，Ticker 表达周期通知。接收处理较慢时，Ticker 不保证为每个周期保存独立事件。停止 Ticker 不会关闭其 channel，因此 `for range ticker.C` 还需要独立退出信号。

## 正则表达式

regexp 采用 Go 支持的 RE2 风格语法，避免某些回溯型表达式的指数级问题，但仍需控制输入尺寸和编译成本。它不支持所有其他语言正则特性，例如反向引用。频繁使用的固定模式可以预编译，动态模式编译失败应返回错误。

```go
re, err := regexp.Compile(`^[a-z][a-z0-9_]{2,19}$`)
if err != nil { return err }
valid := re.MatchString(name)
_ = valid
```

正则适合局部格式，不适合独自验证完整 URL、安全 HTML 或所有电子邮件规则。能用 strings、net/url 或专用解析器时，优先使用表达语义的工具。

## 排序与查找

slices.Sort 处理有序类型，SortFunc 接收比较函数，SortStableFunc 保持相等项的原顺序。比较函数应满足一致的排序关系，不能用溢出风险的 `a.ID-b.ID` 直接比较大整数，优先 cmp.Compare。

BinarySearch 要求输入已按匹配规则排序。map 输出若需稳定顺序，应先排序 key；排序会修改切片底层数组，需要保留原数据时先 Clone。

## 数学与随机数

math 提供浮点计算，math/big 支持大整数和有理数，但对象经常可变，别名复用需要谨慎。随机模拟可使用 math/rand/v2，安全 token 与密钥使用 crypto/rand。伪随机种子固定适合复现测试，不适合安全用途。

## 容器与其他常用包

container/heap 支持优先队列，container/list 提供链表，container/ring 提供环结构。小集合往往用切片更简单，链表并不天然更快。flag 处理 CLI 参数，sort 仍为非泛型接口和旧代码提供工具，errors、fmt、strconv 与 strings 是最常用组合。

练习：写一个按截止时间排序的任务列表，时间相同按 ID 排序；为跨时区日期边界编写测试；解释为什么用随机 sleep 不能作为并发测试同步方式。

延伸阅读：[time](https://pkg.go.dev/time)、[regexp](https://pkg.go.dev/regexp)、[math/rand/v2](https://pkg.go.dev/math/rand/v2)。

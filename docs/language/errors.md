# 10 · error、错误链、panic 与 recover

## error 是普通接口

`error` 只要求 `Error() string` 方法。失败作为值返回，使调用者明确决定重试、降级、记录或向上返回。约定上最后一个返回值是 error，nil 表示成功。

```go
func ParseAge(s string) (int, error) {
    age, err := strconv.Atoi(s)
    if err != nil { return 0, fmt.Errorf("解析年龄 %q: %w", s, err) }
    if age < 0 || age > 150 { return 0, errors.New("年龄超出范围") }
    return age, nil
}
```

错误文字解释失败操作，不重复堆叠“error”字样；必要时包含安全的标识信息。不要把密码、完整认证头或敏感业务数据直接放入错误。

## 包装与检查

`fmt.Errorf("读取配置: %w", err)` 保留底层错误链，`%v` 只格式化文字。调用者用 `errors.Is` 检查语义，用 `errors.As` 找到具有字段的错误类型，而不是比较错误字符串。

```go
if errors.Is(err, os.ErrNotExist) {
    // 可按契约采用默认值
}
var pathErr *os.PathError
if errors.As(err, &pathErr) {
    fmt.Println(pathErr.Op, pathErr.Path)
}
```

包装一个底层错误等于把它变成可观察 API 的一部分。库如果承诺暴露某种错误，未来替换存储实现时要考虑兼容性。`errors.Join` 可以聚合多个错误，Is/As 会遍历其组成；它不意味着多个操作具有事务原子性。

## 哨兵错误与自定义类型

少量稳定条件可以用导出哨兵错误，例如 `ErrNotFound`。需要结构化信息时定义错误类型，例如包含字段名与原因的 ValidationError。不要为每一条动态文字都定义全局变量，也不要让 HTTP 状态码侵入底层通用库。

错误应在能决定行动的边界记录一次，再转换为对外响应。如果每层都 log 后原样 return，最后日志会重复，而且调用链难以阅读。HTTP 入口可把内部错误转换成稳定响应，同时将诊断信息留在服务端。

## panic 适用于什么

panic 表示当前调用无法正常继续，典型场景是不可恢复的不变量破坏或程序错误。普通用户输入、文件不存在和网络超时通常使用 error。panic 会沿当前 goroutine 栈展开，执行 defer；其他 goroutine 的 recover 不能代为捕获。

```go
func invoke(fn func()) (err error) {
    defer func() {
        if v := recover(); v != nil {
            err = fmt.Errorf("执行插件失败: %v", v)
        }
    }()
    fn()
    return nil
}
```

recover 必须在延迟调用的函数中直接使用才具有预期效果。恢复后不会从 panic 语句下一行继续执行，而是让对应函数返回。不要用一个大 recover 把损坏状态藏起来；服务边界恢复时还应记录栈，并定义请求失败后的状态处理。

## 常见陷阱

返回具体错误指针的 nil 值到 error 接口可能得到非 nil 接口，详见接口章节。忽略 `rows.Err()` 会漏掉遍历阶段错误；只检查创建文件而不检查写入会漏掉磁盘失败；随意重试非幂等写操作可能重复扣款或重复提交。

练习：定义带资源 ID 的 NotFoundError，为它实现 Error；用 errors.As 读取 ID，再通过包装证明直接类型断言和 errors.As 的区别。

延伸阅读：[错误链设计](https://go.dev/blog/go1.13-errors)、[errors 包](https://pkg.go.dev/errors)。

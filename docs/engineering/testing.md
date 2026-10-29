# 31 · 单元测试、示例、模糊测试与基准

## 测试验证可观察行为

测试文件以 `_test.go` 结尾，测试函数形如 `func TestName(t *testing.T)`。断言关注输入输出、错误契约和资源生命周期，不应只是复刻实现每一行。

```go
func TestReverse(t *testing.T) {
    tests := []struct { name, in, want string }{
        {"空串", "", ""},
        {"中文", "Go中文", "文中oG"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := Reverse(tt.in); got != tt.want {
                t.Fatalf("Reverse(%q) = %q, want %q", tt.in, got, tt.want)
            }
        })
    }
}
```

表驱动让边界并列展示；子测试便于定位。使用 `t.Helper()` 标记辅助函数，`t.Cleanup()` 注册资源清理，`t.TempDir()` 创建隔离临时目录。单元测试不要默认写用户目录或连接公网服务。

## 测试包与依赖边界

同包测试可以访问未导出成员，外部 `package foo_test` 更接近使用者视角。两者按目标组合。接口注入可替换网络或存储，但无需为纯函数强行增加接口。

mock 只能验证约定的交互，不能证明真实 SQL、TLS、序列化或事务语义。把纯逻辑单元测试、协议测试和真实环境集成测试分清，并记录后者需要哪些服务。

## HTTP 与时间测试

httptest.NewRecorder 测 Handler，NewServer 测完整本地 HTTP 往返。检查状态、响应头、body 和副作用；不要只看状态码。对超时和并发，使用可控制信号而不是“睡一小会应该够了”。

Go 1.25+ 的 testing/synctest 为部分并发代码提供隔离环境和虚拟时间，适合验证受支持的等待场景；外部 I/O 不会因此变成可控模拟。先读其阻塞判定与边界，再替换真实 sleep。

## Example 测试

Example 函数可以带 `// Output:`，既展示 API 用法又验证确定输出。无确定顺序的并发输出应先排序或避免声明固定输出，不能把调度顺序写成错误契约。

## 模糊测试

```go
func FuzzReverse(f *testing.F) {
    f.Add("Go中文")
    f.Fuzz(func(t *testing.T, s string) {
        if !utf8.ValidString(s) { t.Skip() }
        if got := Reverse(Reverse(s)); got != s {
            t.Fatalf("两次反转后 = %q, want %q", got, s)
        }
    })
}
```

```sh
go test ./examples/basics -run '^$' -fuzz FuzzReverse -fuzztime 10s
```

模糊测试适合解析器、编码、序列化与性质验证。输入域必须明确：按 rune 反转非法 UTF-8 会改变字节，所以示例限制为合法 UTF-8。失败语料应保留，形成回归测试。

## 基准与覆盖率

```go
func BenchmarkCount(b *testing.B) {
    for b.Loop() {
        Count(strings.NewReader("go go tutorial"))
    }
}
```

Go 1.24+ 的 b.Loop 简化基准循环。基准输入要代表真实工作负载，区分初始化与被测操作。使用 `-benchmem` 看分配，重复采样后比较，不能把单次运行的噪声当作优化收益。

```sh
go test -cover ./...
go test -race ./...
go test ./internal/wordcount -run '^$' -bench . -benchmem
```

覆盖率帮助找未执行路径，但不证明断言质量。并行测试用 t.Parallel 前，确认没有共享全局状态、端口或环境变量冲突。

练习：让 API 测试覆盖未知字段、超限请求、空标题、完成任务和并发创建；对流水线测试取消后退出，避免只靠延迟观察 goroutine 数量。

延伸阅读：[testing](https://pkg.go.dev/testing)、[官方 fuzz 教程](https://go.dev/doc/tutorial/fuzz)、[testing/synctest](https://pkg.go.dev/testing/synctest)。

# 39 · 实战一：文本统计 CLI

## 需求与明确的文本规则

程序支持直接文本、文件或标准输入，统计行数、词数、不同词数与词频，并支持 JSON 输出。词以 Unicode 空白分隔，统一小写，保留标点；因此 `Go` 和 `go` 合并，`go!` 是另一个词。中文没有空白时视作一个 token，这个程序不是自然语言分词器。

行数采用 bufio.ScanLines 语义：空输入 0 行，`go\n` 为 1 行，`\n\n` 为 2 行。单行 token 上限为 1 MiB，超出返回错误。词频 map 仍随不同词数量增长，大规模不受信输入还需独立的总量和不同词数量限制。

## 运行

```sh
go run ./cmd/wordcount -text "Go go 中文" -top 10
go run ./cmd/wordcount -text "Go go 中文" -json
go run ./cmd/wordcount -file examples/files/message.txt
go test ./internal/wordcount ./cmd/wordcount
```

```text
lines=1 words=3 unique=2
go      2
中文    1
```

无 `-text` 或 `-file` 时读取 stdin；`-text` 与 `-file` 互斥。`-top 0` 输出全部词频，负值属于参数错误。显式 `-text ""` 表示空文本，不会意外退回 stdin。

POSIX shell 可使用 `printf 'go go\n中文\n' | go run ./cmd/wordcount`。Windows 终端管道还受编码配置影响，初学时先使用 `-text` 或 UTF-8 文件排除终端编码干扰。

## 核心算法

```go
--8<-- "internal/wordcount/count.go"
```

Count 接收 io.Reader，因此不关心数据来自文件、内存还是管道。Scanner 逐行处理，Fields 分词，map 累计，最后转换为切片排序。输出先按次数降序、相同次数按词升序，避免 map 随机次序影响测试与用户阅读。

若总 token 数为 N，不同词数为 U，统计通常需要 O(N) 次 map 操作和 O(U) 存储，最后排序约 O(U log U)。流式读行不等于整个算法常量内存，因为词频集合仍可能增长。

## 命令入口

```go
--8<-- "cmd/wordcount/main.go"
```

将 run 与 main 分开，测试可以注入 stdin/stdout/stderr 并观察退出码，无需真的启动子进程。main 只负责连接系统资源与最终 os.Exit，run 返回前 defer 已有机会执行。

退出码 0 表示成功或帮助信息，2 表示参数用法错误，1 表示文件、扫描或输出失败。标准输出只包含结果，诊断写标准错误，便于 shell 组合与 JSON 消费。

## 测试如何覆盖风险

核心测试覆盖空输入、空行、末尾换行、大小写、Unicode 空白、标点、词频并列排序、读取中途错误和超长行。入口测试覆盖 stdin、显式空文本、参数互斥、缺失文件与 JSON 输出。这样测试关注契约，允许未来替换内部扫描算法。

```sh
go test -v ./internal/wordcount ./cmd/wordcount
go test ./internal/wordcount -run '^$' -bench . -benchmem
```

## 扩展练习

增加停用词列表，但说明大小写和编码规则；增加按 rune 总数统计，并明确是否包含换行；对超大输入用外部排序或分片统计，保留相同输出契约；为总字节数和不同词数加上限，并测试达到上限与超过上限的区别。

完成后写一段说明：哪些输入被接受、哪些被拒绝、哪些数据可能占用大量内存。可靠 CLI 的质量不只来自正常示例能输出结果。

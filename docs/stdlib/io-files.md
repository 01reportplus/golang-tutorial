# 24 · I/O、文件、缓冲与文件系统

## Reader 与 Writer 是组合基础

`io.Reader` 的 Read 返回读取字节数与错误，`io.Writer` 的 Write 返回写入字节数与错误。Reader 可能同时返回 `n > 0` 与错误，调用者应先处理已读数据，再处理错误；一次 Read 也不保证填满缓冲区。

```go
for {
    n, err := r.Read(buf)
    if n > 0 { consume(buf[:n]) }
    if err == io.EOF { break }
    if err != nil { return err }
}
```

生产代码优先使用 io.Copy、io.ReadFull、io.CopyN 等已有工具。ReadFull 适合固定长度协议字段，Copy 适合流式传输；不要为了方便对无限输入直接 ReadAll。

## 文件操作与权限

os.Open 只读打开，os.Create 创建或截断，os.OpenFile 控制打开标志。先检查错误，再安排 Close。权限模式如 `0o600` 在不同系统与 umask 下有不同实际效果，Windows 不完全采用 POSIX 权限语义。

```go
data, err := os.ReadFile(path) // 适合有大小上限的小文件
if err != nil { return fmt.Errorf("读取 %s: %w", path, err) }
if err := os.WriteFile(target, data, 0o600); err != nil { return err }
```

WriteFile 不是完整的跨平台原子更新或持久化事务。关键文件常采用同目录临时文件、写入、必要时 Sync、关闭、重命名等步骤；还要考虑目录同步、权限、目标平台与失败清理。不要把一个 rename 调用等同于所有文件系统上的断电安全。

## bufio：缓冲与行处理

bufio.Reader 适合多次小读取，Writer 需要显式 Flush，并检查 Flush 错误。Scanner 适合按行或自定义 token 扫描，但默认 token 大小有限，需要 Buffer 调整上限；超长行必须有明确处理策略。

```go
scanner := bufio.NewScanner(r)
scanner.Buffer(make([]byte, 4096), 1024*1024)
for scanner.Scan() {
    line := scanner.Text()
    process(line)
}
if err := scanner.Err(); err != nil { return err }
```

Scanner.Bytes 返回的缓冲可能被下一次 Scan 覆盖，若要长期保存应复制。Text 返回 string 更方便，但有分配成本，先按需求选而不是盲目微调。

## 路径与目录遍历

filepath 处理操作系统路径，path 处理 slash 分隔的逻辑路径，例如 URL 路径或 fs.FS 路径。filepath.Join 负责拼接，却不会自动证明用户输入没有越界。Clean 只做词法整理，不会消除符号链接带来的逃逸。

读取不可信相对路径时，要限制根目录、检查路径语义，并按版本考虑 os.Root 等支持限定目录访问的 API。仅用字符串前缀判断 `/safe` 会误判 `/safe-other`，也无法解决符号链接竞态。

filepath.WalkDir 遍历目录，回调中应处理传入 err，按需 SkipDir；大量文件处理要限制打开文件数。临时资源使用 os.CreateTemp/os.MkdirTemp，成功后安排清理，避免手工生成可预测名称。

## fs.FS 与可测试设计

io/fs 抽象只读文件系统，既能接收磁盘目录，也能接收 embed.FS 或测试内存文件系统。函数接收 fs.FS 而不是固定磁盘绝对路径，更容易测试与打包。fs.FS 中路径遵守其自身规则，不能随意传 Windows 绝对路径。

## 命令执行

os/exec 默认不通过 shell 解析，程序名与参数分开传入；这有助于避免 shell 注入。用户可控参数仍可能触发目标工具选项，应校验并按工具支持使用 `--`。CommandContext 取消进程时也要理解子进程组与平台行为。

练习：为文本统计 CLI 测试空文件、超长行、读取中途失败；实现大小受限的读取函数并区分“刚好达到限制”和“超过限制”。

延伸阅读：[io](https://pkg.go.dev/io)、[bufio](https://pkg.go.dev/bufio)、[os](https://pkg.go.dev/os)、[io/fs](https://pkg.go.dev/io/fs)。

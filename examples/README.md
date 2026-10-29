# 可运行示例

在仓库根目录执行：

| 命令 | 主题 |
| --- | --- |
| `go run ./examples/hello` | 程序入口与 Unicode 输出 |
| `go run ./examples/collections` | 切片复制、map 稳定输出 |
| `go run ./examples/interfaces` | 小接口与 Writer 组合 |
| `go run ./examples/generics` | 泛型算法与函数迭代器 |
| `go run ./examples/concurrency` | 等待完成与固定索引写入 |
| `go run ./examples/json` | 可选字段与 Go 1.26 new 表达式 |
| `go run ./examples/files` | embed、fs.FS、流式扫描 |
| `go run ./examples/pipeline` | 有界并发、按输入顺序汇总 |
| `go test -v ./examples/basics` | nil、错误链、闭包、切片与虚拟时间 |

综合项目：`cmd/wordcount`、`cmd/taskapi`。全部 Go 代码使用标准库，不依赖外部 Go 模块。数据库、cgo、RPC 和 Docker 的集成前提见对应章节，不能以本目录测试结果代替集成验证。

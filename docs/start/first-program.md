# 03 · 第一个程序与项目结构

## 一个完整程序

```go
package main

import "fmt"

func main() {
    name := "Go"
    fmt.Printf("你好，%s！\n", name)
}
```

把它保存为新实验目录里的 `main.go`，执行 `go run .`。`package main` 表示这是可执行包，`import` 引入另一个包，`main` 是入口，`:=` 在函数内声明并初始化变量，`Printf` 中 `%s` 是字符串占位符，`\n` 是换行。

`fmt.Printf` 的 `P` 大写，所以可以被其他包访问。Go 使用标识符首字母是否大写决定包级成员的导出性，不使用 `public` 或 `private` 关键字。

## run、build、install 的区别

| 命令 | 做什么 | 常见用途 |
| --- | --- | --- |
| `go run .` | 编译临时程序并执行 | 快速试验 |
| `go build -o hello .` | 生成当前平台可执行文件 | 交付和部署 |
| `go install .` | 构建并安装命令 | 个人命令工具 |
| `go test ./...` | 编译并运行当前 module 各包测试 | 持续验证 |

`go run main.go` 只包含命令行指定的文件，可能遗漏同目录其他 `.go` 文件；日常通常使用 `go run .`。程序参数放在包路径后，例如 `go run ./cmd/wordcount -text "hello go"`。

## 文件不是包边界

同目录普通 `.go` 文件通常属于同一包，声明相互可见。一个 module 可以包含许多包。测试文件以 `_test.go` 结尾；构建约束和 `_windows.go` 等文件后缀会影响参与构建的文件集合。

常见结构如下，目录名是习惯，不是所有项目都必须遵守的模板：

```text
myapp/
  go.mod
  cmd/myapp/main.go
  internal/task/service.go
  internal/task/service_test.go
```

`cmd/` 可放多个命令；`internal/` 有 Go 工具链施加的导入可见性限制。小项目完全可以从根目录一个 `main.go` 开始，复杂之后再按职责拆分，不必提前制造多层空目录。

## 输入输出和退出状态

标准输入用于接收数据，标准输出用于正常结果，标准错误用于诊断信息。命令失败应返回非零退出码。`os.Exit` 会立即终止进程，不执行尚未执行的 defer，所以适合在最外层入口处理已经完成清理的错误。

如果函数中打开文件、启动事务或分配其他外部资源，应在该函数返回前安排清理。不要把所有错误处理都推给 `main` 而丢失上下文。

## 从编译错误学习

Go 不允许同作用域重复声明同名变量，不允许未使用的导入和普通局部变量。缺少右括号可能导致下一行出现报错，先修复最早的错误再重新编译。用 `gofmt` 处理格式而不是手动排列空格。

练习：把程序拆成 `main.go` 和 `greet.go`，在后者声明返回问候语的函数；比较 `go run main.go` 与 `go run .`。再把程序改成读取一个命令行参数，对参数缺失返回明确错误。

延伸阅读：[编写 Go 代码](https://go.dev/doc/code)。

# 02 · 安装 Go 与配置开发环境

## 安装与确认

从 [Go 官方下载页](https://go.dev/dl/)选择操作系统和 CPU 架构对应的安装包。Windows 使用安装器，macOS 可使用官方安装包，Linux 按官方说明解压到单独目录。升级时不要把新旧 Go 目录混合覆盖。安装完成后打开一个新终端：

```sh
go version
go env GOROOT GOPATH GOOS GOARCH GOMOD GOTOOLCHAIN
```

`GOROOT` 是工具链安装位置，通常无需手动设置。`GOPATH` 默认保存下载模块、构建缓存相关资源和 `go install` 安装的命令；现代项目不要求放在 GOPATH 内。`GOMOD` 用于确认当前找到哪个 `go.mod`，在错误目录运行命令是初学者最常见的问题之一。

## 建立独立实验目录

不要在教程仓库里随意执行另一个 `go mod init`。为自己的练习创建同级目录：

```sh
mkdir hello-go
cd hello-go
go mod init example.com/hello-go
```

`example.com` 是示例域名，正式发布库时 module 路径应与实际导入路径一致。`go.mod` 中的 `go` 指令声明最低 Go 版本并影响语言语义；`toolchain` 指令可建议工具链。启用自动工具链选择时，Go 命令可能下载合适版本，离线环境应事先准备并检查 `GOTOOLCHAIN`。

## 代理与私有依赖

先检查当前配置：

```sh
go env GOPROXY GOPRIVATE GONOPROXY GONOSUMDB
```

`GOPROXY` 决定模块下载来源，逗号与竖线的回退行为不同。网络不可达时先判断 DNS、代理和组织策略，不要直接关闭校验。`GOPRIVATE` 为私有模块路径配置隐私和校验策略，仍需单独配置 Git 的访问凭据。不要把 token 写入 `go.mod`、代码或提交日志。

修改持久配置前记录旧值。`go env -w` 会影响后续项目；临时实验优先使用终端环境变量，结束后恢复。教程本身的 Go 示例只依赖标准库，首次运行无需下载第三方 Go 模块。

## 编辑器与格式化

VS Code 可配合官方 Go 扩展与 `gopls`；GoLand 也能识别 module 和测试。编辑器报错时，先在终端运行 `go test ./...`，区分工具链错误与索引问题。统一开启保存时格式化，学习 `go doc fmt.Println`、跳转定义和运行单个测试。

## Windows 注意事项

PowerShell 设置环境变量用 `$env:NAME = 'value'`，POSIX shell 使用 `NAME=value command`，两种写法不能混用。下面是临时交叉编译示例，完成后恢复变量：

```powershell
$env:GOOS = 'linux'
$env:GOARCH = 'amd64'
go build -o hello-linux .
Remove-Item Env:GOOS
Remove-Item Env:GOARCH
```

`-race` 在部分平台需要启用 cgo 和安装兼容 C 编译器；普通 `go test` 通过不代表 race 检测已经执行。仓库 CI 在 Linux 上执行 race 检查。

## 排障顺序

找不到 `go` 时检查 PATH 和新终端；找不到主模块时检查目录与 `go.mod`；下载超时时检查代理；版本不满足时检查 `go version` 与 `go.mod`；权限失败时检查缓存目录的访问权限。不要用删除整个工具链或关闭安全校验作为第一步。

练习：把 `go env` 中上述变量的含义写成自己的话，并确认项目根目录与子目录的 `GOMOD` 是否相同。

延伸阅读：[官方安装说明](https://go.dev/doc/install)、[工具链选择](https://go.dev/doc/toolchain)。

# 34 · 构建、代码生成、容器与 CI

## 理解目标平台

GOOS 和 GOARCH 决定目标操作系统与架构，CGO_ENABLED 影响跨语言链接。纯 Go 程序通常容易交叉编译；使用 cgo 时还需要对应目标的 C 工具链与库。

```sh
go build -trimpath -o bin/wordcount ./cmd/wordcount
go build -trimpath -o bin/taskapi ./cmd/taskapi
go version -m bin/taskapi
```

POSIX shell 临时交叉编译可写 `GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build ...`；PowerShell 使用环境变量形式，见安装章节。文件编译成功不代表能在当前机器运行目标平台二进制。

## 构建约束

`//go:build` 控制文件在特定平台或标签下参与构建，文件名 `_linux.go`、`_amd64.go` 也具有规则。约束注释需要位于正确位置并与 package 留空行。标签不是运行时开关，未参与构建的代码仍需在相关 CI 组合中检查。

## 代码生成

```go
//go:generate go run ./tools/generate
```

go generate 是显式命令，go build 不会自动执行它。生成器版本、输入和命令应可复现；生成文件注明来源，不要直接手工修改后期待下次生成保留。API schema、protobuf、数据库查询生成与 mock 都适合遵循这一流程。

如果源码是权威契约，应先改 schema 再生成派生文件，最后查看 diff 和兼容性。生成成功不等于协议兼容，仍要验证字段号、语义与协调发布。

## 多阶段容器构建

```dockerfile
FROM golang:1.26 AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/taskapi ./cmd/taskapi

FROM scratch
COPY --from=build /out/taskapi /taskapi
USER 65532:65532
EXPOSE 8080
ENTRYPOINT ["/taskapi", "-addr", ":8080"]
```

这是纯 Go API 的最小容器示意。出站 HTTPS 需要 CA，具名时区可能需要时区数据，写文件需要可写目录；scratch 不会自动提供这些。实际发布应固定补丁版本或镜像摘要，并为架构、探针、资源和权限建立部署配置。

`.dockerignore` 排除 Git、虚拟环境、构建站点和本地数据，避免上下文过大或泄露文件。本仓库提供对应 Dockerfile，但容器运行需要本机或 CI 的 Docker 环境验证。

## CI 应验证什么

先执行格式和静态检查，再运行测试与 race，再构建目标。文档站点也需要检查链接、代码块和构建结果。PR 只做验证，部署工作流只在默认分支执行，并把 Pages 写权限限制在部署 job。

本仓库 `.github/workflows/ci.yml` 检查 Go 1.26/1.27，Pages 工作流跟随默认分支。修改默认分支后还需检查环境保护规则是否允许新分支部署，工作流条件本身不能绕过仓库权限。

## 发布与回滚

发布记录应包含源码 commit、工具链、构建参数、依赖和产物校验。数据库迁移应先考虑前后版本兼容，再滚动更新程序。回滚代码不一定能回滚已执行的数据迁移，应提前设计兼容窗口。

静态教程的发布产物是 HTML/CSS/JS，GitHub Pages 不运行 Go 服务。API 项目需要单独的服务器或容器平台，不能把可执行文件上传到 Pages 就获得后台服务。

练习：为一个平台专属实现编写构建约束；在 CI 添加对应构建检查；说明何时应该提交生成代码，以及如何验证生成器版本变化。

延伸阅读：[go generate](https://go.dev/blog/generate)、[GitHub Pages 工作流](https://docs.github.com/en/pages/getting-started-with-github-pages/using-custom-workflows-with-github-pages)。

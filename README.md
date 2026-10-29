# Go 全栈学习笔记

从语言基础、并发和标准库，到测试、性能诊断、服务开发与进阶专题的中文 Go 教程。包含 **41 个编号章节、50 个阅读页面、8 个可运行小示例、3 个综合项目**，每章配有练习，附参考答案与官方文档入口。

**[在线阅读](https://01reportplus.github.io/golang-tutorial/)** · **[学习路线](docs/start/roadmap.md)** · **[完整目录](mkdocs.yml)** · **[发布说明](docs/appendix/publishing.md)**

语言与可执行示例基线为 **Go 1.26**，单独说明 Go 1.27 差异。测试只依赖 Go 标准库；网站使用 MkDocs Material，支持中文搜索、章节导航、代码复制、深浅色切换和移动端阅读。

## 内容范围

| 部分 | 内容 |
| --- | --- |
| 01—04 入门 | 安装、环境、程序结构、工具与调试 |
| 05—17 语言 | 声明、类型、常量、控制流、函数、错误、指针、Unicode、容器、结构体、接口、泛型与迭代器 |
| 18—23 并发 | goroutine、channel、select、锁、原子操作、context、背压、内存模型 |
| 24—29 标准库 | I/O、文件、编码、时间、正则、HTTP、模板、embed、SQL 与事务 |
| 30—34 工程 | Modules、工作区、测试、fuzz、benchmark、profile、GC、配置日志、安全、构建和 CI |
| 35—38 进阶 | 反射、unsafe、cgo、汇编、Wasm、RPC、Protobuf、分布式系统与架构 |
| 39—41 实战 | 文本统计 CLI、任务 HTTP API、有界并发流水线 |
| 附录 | 规范覆盖索引、标准库导航、版本差异、FAQ、练习答案、Pages 发布与参考资料 |

## 运行示例

```sh
go test ./...
go run ./examples/hello
go run ./cmd/wordcount -text "Go go 中文" -json
go run ./examples/pipeline
go run ./cmd/taskapi
```

API 默认地址 `http://127.0.0.1:8080`，接口与调用方式见[任务 API 实战](docs/projects/task-api.md)。数据只保存在内存中，重启清空。[全部示例说明](examples/README.md)。

## 本地阅读网站

Windows PowerShell：

```powershell
python -m venv .venv
.\.venv\Scripts\python.exe -m pip install -r requirements.txt
.\.venv\Scripts\python.exe -m mkdocs serve
```

macOS/Linux：

```sh
python3 -m venv .venv
.venv/bin/python -m pip install -r requirements.txt
.venv/bin/python -m mkdocs serve
```

然后打开 `http://127.0.0.1:8000/`。不启动网站也能直接在 GitHub 阅读 docs/ 中的 Markdown；项目文章里的源码嵌入指令由站点构建展开，源码本身保存在 cmd/ 与 internal/。

## 验证

使用安装了 requirements.txt 的 Python：

```sh
go vet ./...
go test ./...
python scripts/check_docs.py
python scripts/check_snippets.py
python -m mkdocs build --strict
python scripts/check_site.py
```

支持 race 的环境另执行 `go test -race ./...`。CI 在 Linux 上以 Go 1.26/1.27 检查，并运行短时 fuzz。完整检查结果、环境限制与未执行的外部集成见 [VALIDATION.md](VALIDATION.md)。

## GitHub Pages

仓库 Settings → Pages → Source 选择 **GitHub Actions**，推送到默认分支即触发部署。Pages 工作流根据仓库当前默认分支判断，不把 main/master 写死。修改默认分支时还需检查 github-pages 环境允许分支。

质量检查与部署分别由 `.github/workflows/ci.yml` 和 `pages.yml` 管理。网页只包含静态教程；Go API 项目需要单独运行或部署。

## 参考与许可

参考用户指定的菜鸟教程、Go 语言圣经、Go 语言高级编程、Go 官方规范与文档、Golang 中文学习文档、W3Schools 和 GeeksforGeeks；内容重新组织和独立编写，技术结论优先核对官方资料。[参考与版权](docs/appendix/references.md)。

原创文字：CC BY 4.0；原创代码与配置：MIT。见 [LICENSE](LICENSE) 与 [LICENSE-CODE](LICENSE-CODE)。

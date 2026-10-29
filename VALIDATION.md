# 验证记录

验证日期：2026-10-29（Asia/Shanghai）。环境：Windows amd64，Go 1.26.5，Python 3.14.5，MkDocs 1.6.1 / Material 9.8.0。

## 本机执行

| 检查 | 结果 |
| --- | --- |
| `gofmt -l cmd examples internal` | 无未格式化文件 |
| `go vet ./...` | 通过 |
| `go test -count=1 -cover ./...` | 通过；basics 100%、wordcount 100%、pipeline 96.3%、taskapi 95.9%、CLI run 82.2% |
| `go test ./examples/basics -run '^$' -fuzz FuzzReverse -fuzztime 5s` | 通过，约 90 万次执行 |
| `go build ./cmd/...` | 通过 |
| 8 个 examples 命令与 wordcount CLI | 全部运行通过 |
| wordcount 基准 | 运行通过；本机单次观察约 19.7 µs/op，不能用作跨环境性能结论 |
| `scripts/check_docs.py` | 50 个页面全部在导航中，Markdown 本地链接与源码嵌入路径有效 |
| `scripts/check_snippets.py` | 16 个独立完整 Go 代码块编译通过 |
| `mkdocs build --strict` | 通过 |
| `scripts/check_site.py` | 51 个 HTML 页面（含 404），内部资源与锚点有效，中文搜索索引非空 |
| 浏览器抽查 | 首页布局、中文搜索“切片”及搜索结果已检查 |

覆盖率按相应包统计，不是整个教程所有知识的覆盖率。命令 main 与演示入口的执行路径不都由单元测试覆盖；纯演示程序另以运行验证。完整代码块只编译，不自动运行网络服务；局部片段和外部环境示意不在完整代码块编译统计中。

Windows 沙箱的临时文件重命名曾阻止覆盖率输出，改为正常用户权限执行后通过；这属于验证环境权限问题。未改变用户系统安全设置。

## 远端 CI 与 Pages

工作流已配置 Go 1.26/1.27 的 Linux 测试、race、短时 fuzz，以及文档和站点检查。Pages 只从仓库默认分支部署。

本次工程保存在本地并已提交；远端因当前认证身份缺少目标仓库写权限而未能推送，因此远端 CI、Go 1.27 矩阵与 Pages 部署均尚未执行。后续获得仓库写权限，并在 Settings → Pages 将 Source 设为 GitHub Actions 后，再推送、检查 Actions 和实际网页。不能以工作流存在替代成功执行证据。

## 尚需专门环境的验证

- 本机 CGO_ENABLED=0，未在本机执行 race；由 Linux CI 执行并检查结果。
- SQL 章节需要选定数据库驱动和真实数据库，未执行数据库集成测试。
- cgo、汇编、Wasm 与 gRPC 章节包含带前提的讲解和示意，未声称已完成跨平台集成验证。
- Dockerfile 已提供，本机没有可用 Docker 命令，未执行容器构建和运行。
- 教学 API 无持久化、认证、分页或生产压测，不以本地测试宣称可直接用于生产。

## 复验

安装 requirements.txt 后，执行 README 中的检查命令。外部链接可能随时间失效，当前脚本检查本地文档和生成站点，不对所有外链做持续可用性保证。Go 新版本差异请对照官方发布说明和包文档复验。

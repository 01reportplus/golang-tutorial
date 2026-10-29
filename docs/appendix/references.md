# 参考资料、核对原则与版权

## 本教程如何使用参考资料

教程按从入门到工程实践的学习路径重新组织内容，讲解、项目和测试独立编写。参考资料用于选题、查漏和核对，不复制整章，不搬运书中的大段叙述与图片。技术规则优先使用当前 Go 官方规范、包文档和发布说明，第三方教材提供阅读路径和补充视角。

用户提供的八个入口均已查阅；核对日期为 2026-10-29。旧教材可能早于泛型、Modules、循环变量和新标准库，读到版本差异时请结合版本附录。

## 用户提供的参考入口

| 资料 | 适合怎样使用 |
| --- | --- |
| [菜鸟教程：Go](https://www.runoob.com/go/go-tutorial.html) | 按基础知识点快速查阅和复习 |
| [Go 语言圣经（中文版）](https://golang-china.github.io/gopl-zh/index.html) | 系统理解语言、组合、接口和并发 |
| [Go 语言高级编程](https://chai2010.cn/advanced-go-programming-book/index.html) | cgo、汇编、RPC、Web 与分布式专题 |
| [Go 语言规范](https://go.dev/ref/spec) | 语言规则的权威定义 |
| [Golang 中文学习文档](https://golang.xiniushu.com/guide.html) | 按主题发现中文补充资料 |
| [Go 官方文档](https://go.dev/doc/) | 安装、教程、数据库、诊断与工具链 |
| [W3Schools Go Tutorial](https://www.w3schools.com/go/) | 短例子与基础语法练习 |
| [GeeksforGeeks Go Tutorial](https://www.geeksforgeeks.org/go-language/go/) | 基础、数据结构与题目补充 |

## 权威核对入口

| 主题 | 文档 |
| --- | --- |
| 语言 | [Specification](https://go.dev/ref/spec) |
| 并发可见性 | [Memory Model](https://go.dev/ref/mem) |
| 模块版本 | [Modules Reference](https://go.dev/ref/mod) |
| 标准库 | [pkg.go.dev/std](https://pkg.go.dev/std) |
| 代码风格 | [Effective Go](https://go.dev/doc/effective_go)、[Code Review Comments](https://go.dev/wiki/CodeReviewComments) |
| 性能与内存 | [Diagnostics](https://go.dev/doc/diagnostics)、[GC Guide](https://go.dev/doc/gc-guide) |
| 版本 | [1.26](https://go.dev/doc/go1.26)、[1.27](https://go.dev/doc/go1.27)、[Release History](https://go.dev/doc/devel/release) |
| 发布 | [GitHub Pages 自定义工作流](https://docs.github.com/en/pages/getting-started-with-github-pages/using-custom-workflows-with-github-pages) |
| 文档工具 | [Material for MkDocs](https://squidfunk.github.io/mkdocs-material/) |

## 发生冲突时怎么办

先记录工具链和 module 的 go 版本，判断争议是语言规则、库契约还是实现观察。语言查规范，库查当前版本包文档，实现查对应版本源码。博客中的性能数字、内存布局和旧版限制不能直接推广到所有环境。

引用第三方代码时应额外保留其许可证与署名，不能因为“网上公开”就默认可整体复制。本仓库未把外部站点的许可套用到原创内容，也不声称获得原作者背书。

## 本仓库许可

原创教程文字按 CC BY 4.0 提供，原创 Go/Python/CSS 与配置代码按 MIT 提供；详细说明见仓库 LICENSE 与 LICENSE-CODE。引用、改编时保留署名与来源，第三方依赖、工具、链接目标和其商标仍遵守各自许可。

## 如何贡献

先报告明确问题：版本、位置、可复现代码、期望与实际。修正文档时同时更新示例与测试，避免只改一处。新增主题应包含前提、核心机制、边界、练习和官方参考；不能只添加目录或一段没有验证的结论。

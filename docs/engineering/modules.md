# 30 · 包、Modules、版本与工作区

## 三种边界

文件组织代码，包提供封装与编译单位，module 提供依赖版本单位。import 路径指向包，不等于仓库名字，也不等于 package 声明名。一个仓库通常一个 module，也可以包含多个，但多 module 会增加发布、测试与工作区管理成本。

循环导入在 Go 中不允许。出现循环时，通常说明职责耦合：把共享类型提取到合适小包，调整依赖方向，或让使用方声明小接口。不要为了打破循环建立一个装下所有内容的 common 包。

## go.mod 与 go.sum

```text
module github.com/your-account/your-project

go 1.26.0

require example.org/library v1.2.3
```

go.mod 描述模块路径、最低版本和依赖要求；go.sum 保存模块内容校验信息，不是完整的依赖锁文件。Go 通过最小版本选择等规则计算构建列表，而不是永远选择最新版本。

```sh
go mod tidy
go mod verify
go list -m all
go mod graph
go mod why -m example.org/library
```

tidy 依据源码、测试等整理依赖，不应在发布前不看 diff 就直接接受全部变化。verify 检查缓存内容与已有校验记录的一致性，不替代漏洞扫描或对依赖作者的信任判断。

## 增加和安装依赖

`go get module@version` 调整当前 module 依赖；安装独立工具使用 `go install path/to/cmd@version`。固定工具版本有利于 CI 重现，避免每次构建都使用未审查的 latest。Go 1.24+ 可用 go.mod 的 tool 指令管理工具依赖，结合 `go tool` 使用。

Go 1 兼容承诺并不代表所有第三方库都无破坏升级。v2 及以上 module 通常需要 `/v2` 路径后缀。升级前检查 API、行为、最低 Go 版本和传递依赖，运行边界测试。

## replace 的含义

replace 把模块来源替换成本地目录或另一个版本，适合联调和临时修复。它只影响主 module 的构建，依赖库中的 replace 不会替使用者自动生效。发布库前，应避免遗留只在开发者机器存在的绝对路径。

本地替换目录可以有自己的 Git 历史、分支或子模块状态；go.mod 的版本文字与实际参与编译源码可能不同。排障时查看 go list 的模块来源，不要只看 require 后的版本号。

## go work

```sh
go work init ./service ./library
go work use ./another-module
```

工作区让多个 module 在本地协同，不必反复 replace。是否提交 go.work 取决于仓库结构和团队约定；单库发布测试应考虑 `GOWORK=off`，避免工作区隐藏依赖问题。CI 不能依赖开发机上未提交的同级目录。

## vendor、私有模块与可复现性

vendor 可以把选定依赖源码纳入仓库，适合离线或有审计需求的构建，但要维护同步并保留许可证。私有模块用 GOPRIVATE 等配置控制下载与校验路径，认证交给 Git 或凭据管理器。

可复现构建还受工具链、操作系统、cgo、生成器和构建参数影响。记录版本、生成步骤、环境约束和产物元数据，比只提交 go.sum 更完整。

## 文档与 API 设计

导出标识符配套注释，以名称开头并解释行为、错误和并发契约。包名简短具体，避免 `util`、`manager` 承担模糊职责。使用 internal 表达内部边界，公共 API 一旦发布，演进成本会增加。

练习：创建两个独立 module，用 go work 联调；再关闭工作区验证依赖能否独立解析。解释 go.sum 为什么不是 npm 式的完整版本锁。

延伸阅读：[Modules Reference](https://go.dev/ref/mod)、[工作区教程](https://go.dev/doc/tutorial/workspaces)。

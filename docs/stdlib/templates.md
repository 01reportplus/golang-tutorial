# 28 · 模板、静态资源与 embed

## text/template 与 html/template

text/template 用于通用文本，html/template 按 HTML 上下文自动转义，适合页面输出。二者语法相近，但安全语义不同。网页中的不可信数据不要用 text/template 直接插入 HTML。

```go
package main

import (
    "html/template"
    "os"
)

func main() {
    t := template.Must(template.New("page").Parse(`<h1>{{.Title}}</h1>`))
    err := t.Execute(os.Stdout, struct{ Title string }{Title: "<script>alert(1)</script>"})
    if err != nil { panic(err) }
}
```

输出中的危险字符会被转义，而不是作为 script 执行。`template.HTML` 等类型表示开发者承担安全判断，不能直接把用户输入强制转换进去“修复显示问题”。

## 数据访问与控制结构

`.` 表示当前数据，`{{if ...}}`、`{{range ...}}`、`{{with ...}}` 控制输出，`$` 可保留变量引用，管道把结果传给下一函数。模板应负责展示，复杂业务运算应在 Go 中提前整理成视图模型。

FuncMap 增加辅助函数时，应在 Parse 前注册，并让函数行为可预测。不要在模板函数里悄悄执行慢查询、发网络请求或修改状态，否则渲染性能与失败边界难以观察。

## 解析一次，执行多次

模板常在启动时解析，处理请求时执行。初始化阶段可用 template.Must 表达“模板错误属于部署错误”，运行时用户提交模板则应正常返回解析错误。已经解析完毕且不再修改的模板可按文档并发执行，但写入同一个 Writer 仍需协调。

HTTP 输出可先执行到有限大小的缓冲区，成功后再写响应，避免执行中途失败导致半个 HTML 页面与错误状态混合。缓冲也有内存成本，大页面应有明确策略。

## embed 的作用

```go
import "embed"

//go:embed templates/*.html static/*
var assets embed.FS
```

embed 在编译时把文件加入可执行文件，适合模板、迁移脚本、小型前端静态资源。模式相对声明所在包目录，不能随意跨越 `..`。需要至少匹配有效文件，隐藏文件与特殊模式有单独规则。

嵌入让部署更简单，也意味着资源更新要重新构建程序。不要嵌入密钥、开发数据库或巨型文件；生成文件的许可证也要保留。

## fs.Sub 与静态服务

嵌入的路径可能带 static 前缀，可以用 fs.Sub 获得子文件系统，再交给 http.FileServer。路径前缀、目录列表与缓存头要按产品需求处理，不能仅因为资源只读就认为所有内容适合公开。

对单页应用，应区别浏览器页面路由与真实静态文件缺失；对 API，不要把 404 请求统统回退成 index.html。GitHub Pages 上的本教程使用静态生成器构建 HTML，并不运行 Go HTTP 服务。

## 练习

创建一个模板展示用户提供的名称，验证 script 输入不会执行；把模板嵌入程序并在另一个工作目录运行；解释编译期嵌入与运行时读取对部署、更新和调试的影响。

延伸阅读：[html/template](https://pkg.go.dev/html/template)、[text/template](https://pkg.go.dev/text/template)、[embed](https://pkg.go.dev/embed)。

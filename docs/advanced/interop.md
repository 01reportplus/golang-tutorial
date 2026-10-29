# 36 · cgo、汇编、插件与 WebAssembly

## 跨语言边界先考虑成本

Go 调用 C 常用于复用系统库、数据库或图形计算。收益是成熟生态，代价是构建工具链、平台链接、调用成本和两套内存管理。纯 Go 替代是否更合适，要看依赖能力与测量，不能只按“原生一定快”判断。

本章需要额外 C 编译器或目标运行时；仓库默认 Go 测试不编译 cgo 和汇编示意，避免把概念例子当作已验证的跨平台库。

## 最小 cgo 结构

```text
package main

/*
#include <stdlib.h>
*/
import "C"

import "unsafe"

func main() {
    p := C.CString("hello")
    defer C.free(unsafe.Pointer(p))
    // 把 p 传给约定明确的 C 函数
}
```

这里用 text 代码块明确表示需要单独 cgo 环境的示意。C 前导注释必须与 import "C" 的位置满足 cgo 规则。C.CString 在 C 堆分配，不由 Go GC 自动释放；C.GoString 创建 Go 字符串副本。每次转换都要说明所有权与释放者。

## 指针传递规则

不能随意把含 Go 指针的内存交给 C 长期保存。Go 指针、C 指针、固定内存与回调之间有严格限制，应按当前 cgo 文档核对。runtime.Pinner 和 runtime/cgo.Handle 可处理部分特定需求，但它们不是绕过规则的通用开关。

回调要考虑线程、重入、阻塞与生命周期；C 端保存的句柄必须按协议释放。C 分配的内存不一定进入 Go 堆统计，因此 RSS 增长不能只看 heap profile。

## 构建与链接

动态链接需要部署对应库，静态链接也可能受系统库、许可证和目标平台限制。交叉编译通常需要目标 C 编译器。用构建约束提供纯 Go fallback，有助于隔离平台差异，但两条实现都应测试。

cgo 调用开销使“把每个很小操作都交给 C”可能更慢。适当批处理和减少边界转换，往往比局部计算优化更有效。

## Go 汇编

Go 工具链使用自己的汇编语法和 ABI 规则，不能直接把普通系统汇编粘贴进来。常见用途是极小热点、CPU 指令优化或特殊底层接口。寄存器、栈映射、调用约定和 GC 可见性都需要严格遵守。

阅读编译输出可使用 go tool compile 或 go tool objdump；写汇编前先测量普通 Go 代码和编译器优化。runtime 内部 ABI 可能变化，不应把某篇旧文章中的布局视作稳定接口。

## 动态插件

plugin 包在支持的平台上加载 Go 插件，但平台、工具链、构建标签和依赖一致性要求很严格，卸载也不具有普通动态模块那样的通用能力。跨团队扩展常可考虑进程边界、RPC 或明确数据协议，代价是通信开销与部署复杂度。

## WebAssembly

Go 可以面向浏览器 js/wasm 或 WASI 环境构建，目标选择、宿主 API 和运行方式不同。浏览器不提供普通操作系统文件和线程环境，WASI 也受宿主能力限制。syscall/js 属于特定目标，不能在普通服务器包里无约束导入。

Wasm 不是把服务端程序原样变成网页的按钮。需要检查体积、启动时间、宿主交互、异步模型与可用标准库能力。本教程站点本身是静态 HTML，不需要浏览器运行 Go runtime。

练习：为一个 C 字符串调用画出分配与释放时序；列出交叉编译 cgo 程序与纯 Go 程序的不同前提；说明选择 RPC 扩展边界的优缺点。

延伸阅读：[cgo](https://pkg.go.dev/cmd/cgo)、[Go 汇编指南](https://go.dev/doc/asm)、[Go WebAssembly](https://go.dev/wiki/WebAssembly)、[plugin](https://pkg.go.dev/plugin)。

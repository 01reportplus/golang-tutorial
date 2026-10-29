# 12 · 字符串、字节与 Unicode

## string 是不可变字节序列

Go 的 string 可以保存任意字节，不保证一定是合法 UTF-8。字符串字面量通常是 UTF-8 文本，但来自网络或文件的数据可能不是。`len(s)` 计算字节数，`s[i]` 取得 byte；这与用户看见的字符数量不是一回事。

```go
package main

import (
    "fmt"
    "unicode/utf8"
)

func main() {
    s := "Go中文"
    fmt.Println(len(s), utf8.RuneCountInString(s)) // 8 4
    for i, r := range s {
        fmt.Printf("%d %U %c\n", i, r, r)
    }
}
```

range 返回每个 rune 开始处的字节索引与码点值。Unicode 码点也不是完整的“用户感知字符”：组合重音、旗帜和部分 emoji 由多个码点组成。需要按可见字符截断时，应采用合适的字素簇处理方案，不能单纯 `[]rune` 后截取就声称完全正确。

## 字面量与转义

双引号字符串支持 `\n`、`\t`、`\"` 和 Unicode 转义。反引号字符串适合路径、正则和多行文本，其中转义序列不再按双引号规则解释，但不能直接包含反引号。单引号表示 rune 字面量，例如 `'中'`，而不是 string。

原始字符串不会替你解决 SQL 注入或 HTML 安全问题。它只是不同的文本表示形式，输入仍需在使用边界按相应协议处理。

## 修改与转换

```go
b := []byte("hello")
b[0] = 'H'
s := string(b)
r := []rune("中文")
r[0] = '华'
```

string 与 byte/rune 切片转换表达了新的可变性与编码处理需求，不要依赖编译器某次优化来推断是否共享。`s[:n]` 按字节截取，可能把 UTF-8 字符从中间切开；应先判断需求是字节协议还是文本展示。

## 常用文本工具

`strings` 提供 Contains、HasPrefix、TrimSpace、Split、Fields、Join 和 ReplaceAll。`Fields` 按 Unicode 空白分词，适合简单空白分隔文本，但不是中文自然语言分词器，也不是完整英文词法分析器。

大量拼接可使用 `strings.Builder`，二进制缓冲使用 `bytes.Buffer`。循环里反复 `s += part` 可能产生很多分配，是否值得优化要看基准。Builder 使用后不要复制，也不要在未同步情况下并发写。

`fmt` 适合通用格式化，`strconv` 适合明确的数字与布尔转换。输出调试值时 `%q` 能显示转义和边界，`%#v` 常用于观察 Go 形式的结构；日志中仍应控制敏感数据。

## 大小写、规范化与合法性

`strings.EqualFold` 处理 Unicode 大小写折叠，不等于按任意语言规则比较姓名。Unicode 规范化不由字符串相等操作自动完成。用户名、路径和签名等敏感场景，应先规定规范化、允许字符和比较策略。

`utf8.ValidString` 检查合法性；range 遇到非法编码会产生替换码点，若要区分真实替换字符与非法字节，需要结合解码返回的字节宽度。

练习：实现按 rune 数量截断并附加省略号的函数，明确写出它不处理字素簇；为非法 UTF-8、空串和截断边界编写测试。

延伸阅读：[Go 中的字符串、字节与字符](https://go.dev/blog/strings)、[unicode/utf8](https://pkg.go.dev/unicode/utf8)。

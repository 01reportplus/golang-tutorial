# 标准库导航与选型速查

标准库应按问题查阅，不必逐个背 API。下面列出常见包与适合的入口；全部包及当前版本以 [Go 标准库目录](https://pkg.go.dev/std)为准。包数量和 API 会变化，教程不把这张表当成某个版本的完整清单。

## 文本、数据与容器

| 包 | 主要能力 | 阅读位置 |
| --- | --- | --- |
| fmt、strconv | 通用格式化、数字与文本转换 | [类型](../language/types.md)、[字符串](../language/strings.md) |
| strings、bytes | 文本与字节序列操作、缓冲 | [字符串](../language/strings.md) |
| unicode、unicode/utf8、unicode/utf16 | 字符属性与编码 | [字符串](../language/strings.md) |
| regexp、regexp/syntax | 正则匹配与语法分析 | [常用工具](../stdlib/utilities.md) |
| encoding/json、xml、csv、gob | 结构化编码 | [编码](../stdlib/encoding.md) |
| encoding/binary、hex、base32、base64 | 二进制协议与文本表示 | [编码](../stdlib/encoding.md) |
| slices、maps、cmp、iter | 泛型容器算法与迭代 | [泛型](../language/generics.md) |
| sort、container/heap、list、ring | 排序与传统容器 | [常用工具](../stdlib/utilities.md) |
| math、math/bits、math/big、math/cmplx | 数值、位、大数和复数 | [常用工具](../stdlib/utilities.md) |
| math/rand/v2 | 非安全随机模拟 | [常用工具](../stdlib/utilities.md) |

## 文件、压缩与操作系统

| 包 | 主要能力 | 注意 |
| --- | --- | --- |
| io、bufio | 流与缓冲 | 正确处理 n 与 err、Flush、Scanner 限制 |
| os、io/fs、path/filepath、path | 文件系统与路径 | 平台路径和逻辑路径不同 |
| embed | 编译时嵌入资源 | 不嵌入密钥和超大文件 |
| os/exec、os/signal | 子进程与信号 | 参数边界、进程树和优雅停机 |
| archive/tar、archive/zip | 归档 | 解压防路径逃逸与解压炸弹 |
| compress/gzip、zlib、flate、bzip2、lzw | 压缩流 | 关闭/刷新与解压后大小限制 |
| flag | CLI 参数 | 错误输出、互斥参数和退出码 |
| syscall | 低层系统接口 | 可移植代码优先更高层包，扩展接口查 x/sys |

学习入口：[I/O 与文件](../stdlib/io-files.md)、[模板与 embed](../stdlib/templates.md)、[CLI 实战](../projects/cli.md)。对归档输入必须在解压后的真实目标路径和尺寸上实施限制，不能只检查压缩包大小。

## 网络、服务与安全

| 包 | 主要能力 | 注意 |
| --- | --- | --- |
| net、net/netip、net/url | 连接、地址与 URL | TCP 消息分帧、地址策略与 deadline |
| net/http、http/httptest、http/httputil | HTTP、协议测试与代理工具 | Client 复用、body 限制和关闭 |
| net/mail、mime、mime/multipart | 邮件格式、媒体类型与表单 | 格式解析不等于完整业务验证 |
| net/rpc | 传统 Go RPC | 新系统评估协议与生态 |
| crypto/rand | 安全随机数 | token 与密钥不使用 math/rand |
| crypto/tls、x509 | TLS 与证书 | 保留校验、管理证书与轮换 |
| crypto/aes、cipher、hmac、sha256 | 基础密码学构件 | 优先经过审查的高层方案，不自行设计协议 |
| crypto/subtle | 常量时间相关工具 | 不能让整个应用自动避免所有侧信道 |
| hash、hash/crc32、hash/fnv、hash/maphash | 校验与非密码学散列 | 不替代密码学认证和密码存储 |
| database/sql、sql/driver | 数据库抽象与驱动接口 | 连接池、事务、NULL 与真实数据库测试 |

学习入口：[HTTP](../stdlib/http.md)、[SQL](../stdlib/database.md)、[服务安全](../engineering/production.md)。不同 Go 版本的密码学算法与默认 TLS 策略可能变化，升级时查看发布说明。

## 并发、诊断与元编程

| 包 | 主要能力 | 阅读位置 |
| --- | --- | --- |
| context | 生命周期、取消和元数据 | [context](../concurrency/context.md) |
| sync、sync/atomic | 锁、等待、快照与原子操作 | [同步](../concurrency/sync.md) |
| time | 时间点、时长、定时器 | [时间](../stdlib/utilities.md) |
| errors | 错误链匹配与聚合 | [错误](../language/errors.md) |
| log、log/slog | 文本与结构化日志 | [工程实践](../engineering/production.md) |
| testing、quick、fstest、iotest | 测试与辅助工具 | [测试](../engineering/testing.md) |
| testing/synctest | 受支持的并发虚拟时间测试 | [测试](../engineering/testing.md) |
| runtime、runtime/debug、metrics | 运行时信息与诊断 | [性能](../engineering/performance.md) |
| runtime/pprof、trace、net/http/pprof | profile 与执行跟踪 | [性能](../engineering/performance.md) |
| reflect、unsafe | 动态类型操作和底层内存 | [反射](../advanced/reflection.md) |
| go/ast、parser、token、types、format、doc | Go 源码分析与工具开发 | [工具链](../start/tools.md) |
| plugin、syscall/js | 特定环境扩展与宿主交互 | [跨语言](../advanced/interop.md) |

## 标准库外的边界

`golang.org/x/...` 是扩展模块，需要独立版本管理；数据库驱动、gRPC、ORM、Web 框架和 OpenTelemetry 也不是标准库。导入路径看起来常见，不代表不需要 go.mod 或额外许可检查。

遇到陌生包，先看 Overview、类型方法、并发安全说明和 Example，再根据当前工具链用 `go doc` 验证。不要从博客里截取一个函数调用而跳过资源释放和错误契约。

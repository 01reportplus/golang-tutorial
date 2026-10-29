# 27 · 网络、HTTP 服务与客户端

## 从 TCP 到 HTTP

net 提供 TCP、UDP、地址和连接抽象。TCP 是字节流，没有自动消息边界；应用协议必须自行定义长度、分隔符或帧结构。每次 Read/Write 都可能只处理部分数据，连接应有 deadline 与关闭策略。

HTTP 在连接上增加方法、路径、状态、头部与消息体语义。net/http 已处理大量协议细节，普通 Web 服务不必自己实现 HTTP 解析。

## Handler 与路由

Handler 接口只有 ServeHTTP，HandlerFunc 把普通函数适配成 Handler。Go 1.22+ ServeMux 支持方法与路径变量，适合很多基础 API：

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /hello/{name}", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain; charset=utf-8")
        fmt.Fprintf(w, "你好，%s\n", r.PathValue("name"))
    })
    srv := &http.Server{
        Addr: "127.0.0.1:8080", Handler: mux,
        ReadHeaderTimeout: 5*time.Second,
        ReadTimeout: 10*time.Second, WriteTimeout: 10*time.Second,
        IdleTimeout: 60*time.Second,
    }
    log.Fatal(srv.ListenAndServe())
}
```

这段程序是路由与超时的最小示例。项目实战进一步处理优雅停机、输入限制与测试。开发默认监听回环地址，部署时再显式决定暴露边界。

## 响应头与状态码

必须在 WriteHeader 或首次 Write 前设置响应头；首次 Write 若未显式状态通常发送 200。错误处理应在响应尚未开始时决定状态；已经写出部分响应后不能可靠改成另一个状态。

200 表示成功响应，201 常用于创建，204 没有响应体，400 表示请求格式或校验失败，404 表示资源不存在，409 表示冲突，500 表示内部故障。不要为了统一格式把全部错误返回为 200。

## 中间件

中间件常写成 `func(http.Handler) http.Handler`，用于请求 ID、日志、认证、限流和恢复。顺序影响行为，例如日志放在恢复外层才能记录恢复后的状态。包装 ResponseWriter 时还要考虑 Flusher、Hijacker 等能力是否需要保留。

recover 中间件应记录栈并避免泄露内部信息，但不能撤回已经发送的数据，也不能捕获另起 goroutine 的 panic。

## HTTP 客户端

```go
client := &http.Client{Timeout: 5*time.Second}
req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
if err != nil { return err }
resp, err := client.Do(req)
if err != nil { return err }
defer resp.Body.Close()
if resp.StatusCode != http.StatusOK {
    return fmt.Errorf("下游返回 %s", resp.Status)
}
body, err := io.ReadAll(io.LimitReader(resp.Body, 1024*1024+1))
if err != nil { return err }
if len(body) > 1024*1024 { return errors.New("响应过大") }
```

Do 不会因为 404 或 500 自动返回 error，需要自行检查状态。复用 Client 和 Transport；读完并关闭响应体有利于连接复用。不要在每次请求重新建立 Transport，也不要忽略重定向、代理与目标地址策略。

## URL、查询参数、表单和 Cookie

路径参数通过 PathValue 读取，查询参数通过 r.URL.Query 获取，重复查询参数应按接口契约处理，而不是只凭 Get 取第一个值。`ParseForm` 会合并特定来源的参数，认证或签名场景要明确参数来自 query 还是 body，避免边界混淆。

multipart 上传应先限制总 body，再按请求类型解析，并清理可能生成的临时文件。文件名来自用户输入，不能直接当作服务器磁盘路径。Cookie 需要根据使用场景设置 Secure、HttpOnly、SameSite、Path 和生命周期；Cookie 本身不证明身份，需要验证完整会话或签名。

## 流式响应与 WebSocket

SSE 等流式响应需要 Flush、合适的代理缓冲策略和客户端取消检测。普通短请求的统一 WriteTimeout 可能不适合长时间流，需要按产品协议分开配置。每条消息都要有边界，慢客户端要有背压或断开策略。

WebSocket 需要连接升级与消息协议，标准 net/http 不提供完整高层 WebSocket 客户端/服务端 API。选择专用库时检查维护、尺寸限制、ping/pong、并发读写约束和关闭握手；不要只验证握手成功就宣称连接生命周期完整。

## 服务生命周期与边界

使用 signal.NotifyContext 接收终止信号，调用 Server.Shutdown 停止接收新连接并等待活跃请求，设置单独关闭预算。长连接、WebSocket 和后台任务需额外管理。

上传要限制 body，出站 URL 要防 SSRF，TLS 不应通过 InsecureSkipVerify 随意关闭验证；CORS 是浏览器访问策略，不能替代认证。使用 httptest 在本地验证协议行为，避免单元测试依赖公网。

练习：为缺失路径、错误方法、非法 JSON、取消请求、超大 body 和下游 500 写测试；解释 Client.Timeout 与请求 context 截止时间怎样共同约束请求。

延伸阅读：[net/http](https://pkg.go.dev/net/http)、[httptest](https://pkg.go.dev/net/http/httptest)、[net](https://pkg.go.dev/net)。

# 37 · RPC、Protobuf 与分布式系统边界

## RPC 没有消除网络

RPC 让远程调用看起来像函数，但它仍可能超时、重复、部分失败或发生版本不一致。调用成功前，不知道远端是否执行；超时后，也不能简单认定远端没有执行。API 设计必须暴露这些现实，而不是依靠重试掩盖它们。

net/rpc 是标准库中的传统 RPC 实现，适合学习远程方法调用，现代跨语言系统常采用 gRPC、HTTP/JSON 或其他协议。选型要考虑浏览器接入、流式传输、生态、调试与组织经验。

## Protobuf 契约

```proto
syntax = "proto3";
package tutorial.task.v1;
option go_package = "example.com/project/gen/task/v1;taskv1";

message CreateTaskRequest {
  string title = 1;
  string idempotency_key = 2;
}
message Task {
  int64 id = 1;
  string title = 2;
  bool done = 3;
}
service TaskService {
  rpc CreateTask(CreateTaskRequest) returns (Task);
}
```

字段编号是线上协议标识，删除字段后应保留编号和名称，不能随意复用。兼容性不仅是“生成器能编译”，还包括字段存在性、缺省语义、枚举未知值和 JSON 映射。go_package 与 proto package 负责不同命名空间。

## 生成代码的流程

安装匹配版本的 protoc 与 Go 插件，固定版本，执行生成命令，把 schema 与生成步骤纳入版本控制。先修改 proto，再生成代码，禁止只手改生成文件。服务端与客户端升级应允许兼容窗口，不要要求所有实例瞬时同步。

本仓库不默认下载 gRPC 依赖或连接服务。需要实做时，使用官方 Go quickstart 建立独立 module，再把示例扩展到超时、取消与认证。

## 调用控制

给每次调用设置总预算，并把 context 传到服务端处理和数据库。重试只用于明确可重试且满足幂等条件的情况，结合退避、抖动、最大次数和熔断。拦截器适合日志、认证与追踪，但错误映射应保留稳定语义。

流式 RPC 需要额外考虑双方关闭、接收循环、消息大小、背压和半关闭；流建立成功不代表后续所有消息都成功。连接长期复用，也要处理目标更新与负载均衡。

## 幂等与去重

创建任务时让调用者提供幂等键，数据库中以租户与键建立唯一约束，记录请求摘要和结果。相同键不同内容应拒绝，不能把任意重复键当成功。去重记录有保留时间，超过时间的重复请求如何处理也属于契约。

“至少一次”交付意味着可能重复，消费者应把业务变更与消费记录放在可保证原子性的边界。所谓“恰好一次”必须说明范围，例如某个日志系统内的事务，并不自动覆盖外部数据库和邮件。

## Outbox、Saga 与一致性

Outbox 在业务数据库事务中同时写业务数据与待投递事件，后台可靠投递并重试，消费者去重。它解决提交与发布之间的空隙，但还需处理顺序、积压与重复。

Saga 用多个本地事务和补偿动作管理长流程。补偿不是任意物理回滚，已发出的通知或外部支付可能需要单独业务逆操作。设计前先判断是否真的需要分布式事务，单库事务通常更简单。

## 锁、ID 与缓存

分布式锁需要租约、所有权、失效和 fencing 等设计，不能只写 SETNX 就宣称安全。分布式 ID 要考虑唯一性、时钟回退、热点和排序需求。缓存会带来失效窗口、击穿、雪崩和一致性问题，线程安全 map 不是完整缓存系统。

练习：为创建任务设计幂等键表；模拟“服务端提交成功但客户端超时”并解释重试结果；写出 Outbox 事件在投递成功但标记失败时为何会重复。

延伸阅读：[gRPC Go quickstart](https://grpc.io/docs/languages/go/quickstart/)、[Protobuf Go 指南](https://protobuf.dev/getting-started/gotutorial/)、[net/rpc](https://pkg.go.dev/net/rpc)。

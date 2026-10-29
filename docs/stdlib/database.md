# 29 · database/sql、连接池与事务

## 标准库提供抽象，驱动连接具体数据库

database/sql 不自带 PostgreSQL、MySQL 等驱动。sql.Open 使用注册的驱动创建数据库句柄，通常不会立即建立并验证所有连接；用 PingContext 检查可达性。`*sql.DB` 是可并发使用的连接池，应在应用层长期复用。

本章是需要实际驱动和数据库的集成示例，不属于仓库纯标准库测试已覆盖的运行场景。驱动版本、连接串和占位符请以目标数据库与驱动文档为准。

## 连接池不是越大越好

SetMaxOpenConns 控制打开连接总上限，SetMaxIdleConns 控制空闲连接，SetConnMaxLifetime/IdleTime 控制连接回收。设置需要结合数据库容量、实例数量、事务持续时间与请求并发度。所有实例上限之和不能无视数据库总预算。

池满时调用会等待，context 截止时间应覆盖等待与查询。DB.Stats 可以观察等待次数、等待时长和连接状态，帮助区分数据库慢与连接池排队。

## 参数化查询

```go
// PostgreSQL 风格占位符；其他驱动可能使用 ?。
row := db.QueryRowContext(ctx,
    "SELECT id, name FROM users WHERE id = $1", id)
var user User
if err := row.Scan(&user.ID, &user.Name); err != nil {
    if errors.Is(err, sql.ErrNoRows) { return User{}, ErrNotFound }
    return User{}, fmt.Errorf("查询用户: %w", err)
}
```

参数化保护值的边界，表名、列名和排序方向通常不能作为普通参数占位，需要白名单选择。不要用 fmt.Sprintf 把用户值拼进 SQL。

## 多行查询与 NULL

```go
rows, err := db.QueryContext(ctx, "SELECT id, name FROM users")
if err != nil { return err }
defer rows.Close()
for rows.Next() {
    var id int64
    var name sql.NullString
    if err := rows.Scan(&id, &name); err != nil { return err }
    // name.Valid 区分 NULL 与空字符串
}
if err := rows.Err(); err != nil { return err }
```

Next 返回 false 既可能结束，也可能发生错误，所以必须检查 rows.Err。未关闭 rows 可能长时间占用连接。将 NULL 扫描到不接受 NULL 的字段会失败，应采用 Null 类型或明确转换策略。

## 事务边界

```go
tx, err := db.BeginTx(ctx, nil)
if err != nil { return err }
defer tx.Rollback() // Commit 成功后再调用只会报告已结束
if _, err := tx.ExecContext(ctx, statement, args...); err != nil { return err }
if err := tx.Commit(); err != nil { return fmt.Errorf("提交事务: %w", err) }
return nil
```

事务内所有相关 SQL 应使用 tx，误用 db 会跑到另一个连接上。Commit 失败不能当作成功；网络异常下提交结果也可能需要进一步判断。事务应尽量短，不在其中执行不必要的慢网络请求。

## 用条件更新维护库存

下面是 PostgreSQL 风格的局部实现，假设 `items.stock` 为非负整数，调用方已验证 quantity > 0，并已经建立 tx：

```go
result, err := tx.ExecContext(ctx,
    "UPDATE items SET stock = stock - $1 WHERE id = $2 AND stock >= $1",
    quantity, itemID)
if err != nil { return fmt.Errorf("扣减库存: %w", err) }
affected, err := result.RowsAffected()
if err != nil { return err }
if affected != 1 { return ErrInsufficientStock }
```

条件和修改放在一条更新语句里，避免先 SELECT 得到库存、另一个事务修改后再按旧值覆盖。仍需要明确“不存在”和“库存不足”是否使用同一业务错误，以及数据库约束如何兜底。

如果还要插入订单或操作记录，应放在同一个事务里。事务外发起支付或消息投递并不会获得相同原子性。对序列化失败或死锁重试，通常需要重试整个事务，并确认事务内操作可安全重放。

## Prepare 与查询组织

预处理语句是否提升性能依驱动、连接池与数据库而定，不应把每条只执行一次的 SQL 都预处理。Stmt 需要关闭；在 Tx 中创建的语句受事务生命周期约束。长列表 IN 条件、批量写入与返回自增 ID 的语法，也需要按数据库设计。

把 SQL 保存在靠近数据访问层的位置，并为字段顺序、NULL、时区和扫描类型增加集成测试。不要因为 ORM 隐藏了 SQL 就放弃查看执行计划和实际查询数量。

## 隔离级别与业务一致性

事务不自动阻止所有并发业务问题。库存扣减可使用带条件更新、受影响行数、行锁或乐观版本控制，选择依赖数据库隔离语义。读出余额再写回而没有正确条件，即使每个语句都在事务里，也可能丢失更新。

数据库事务不能覆盖任意外部邮件或消息投递。需要协调时可采用 outbox 等模式，由后台投递并保证幂等。必须说明何时提交、何时重试，以及重复事件如何去重。

## 迁移与测试

迁移脚本应版本化，明确前向兼容、索引创建、锁影响和回滚方案。单元测试可替换存储接口，但 SQL 语法、驱动扫描和隔离行为要用真实数据库集成测试验证，mock 不能证明这些契约。

练习：设计唯一约束防止同一幂等键创建两次任务；在事务里扣减库存，并用两个并发连接验证不会超卖；记录真实数据库版本和隔离级别。

延伸阅读：[访问数据库](https://go.dev/doc/database/)、[database/sql](https://pkg.go.dev/database/sql)。

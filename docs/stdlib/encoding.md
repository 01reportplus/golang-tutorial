# 25 · JSON、二进制、CSV 与数据契约

## 编码是边界契约

内存里的结构体和网络数据不是同一件事。编码要定义字段名称、缺省行为、精度、日期格式与未知字段策略；版本升级还要考虑旧客户端。不要让数据库实体直接决定所有对外 API。

```go
type CreateUser struct {
    Name string `json:"name"`
    Age  *int   `json:"age,omitempty"`
}
```

指针字段能区分部分“未提供”和“提供零值”的情况，但普通解码到新对象时，缺失与 null 都可能成为 nil。PATCH 若需区分缺失、null、具体值三态，应使用带 Present 标志的自定义类型或 RawMessage 检查。

## Marshal 与 Decoder

json.Marshal 生成字节切片，Encoder 流式写出，Decoder 流式读入。只导出字段参与常规编码，tag 控制名称和忽略规则。`omitempty` 不适合表达所有业务缺省；Go 1.24+ 的 `omitzero` 有不同规则，应按版本文档选择。

```go
dec := json.NewDecoder(r)
dec.DisallowUnknownFields()
var req CreateUser
if err := dec.Decode(&req); err != nil { return err }
var extra any
if err := dec.Decode(&extra); err != io.EOF {
    return errors.New("请求体必须只包含一个 JSON 值")
}
```

Decode 一次成功不代表没有尾随第二个 JSON 值。HTTP 入口还应先限制请求体大小，检查媒体类型，再进行结构和业务校验。DisallowUnknownFields 适合严格输入，但是否拒绝未来字段是兼容性选择。

## 数字精度与动态数据

解码到 any 时，传统 encoding/json 默认把数字作为 float64，较大的整数可能丢精度。用明确结构体类型最好；确需动态数据时可 UseNumber，再显式解析 json.Number。ID 若用字符串传输，应在接口契约里一致规定。

RawMessage 用于延迟解析某段 JSON，例如按 type 字段选择 payload 结构。自定义 MarshalJSON/UnmarshalJSON 要避免无意递归调用自身，通常建立不带方法的局部别名辅助处理，并验证 nil 与空值。

## 重复字段与不可信数据

JSON 对象重复字段、字段大小写匹配和未知字段行为会影响安全与兼容。DisallowUnknownFields 并不会自动拒绝所有重复 key。跨语言签名或授权判断不能假定所有解析器行为一致，应规定规范化编码或在边界明确拒绝歧义输入。

不同版本新增的编码包或实验功能不应默默替换现有协议。迁移前对 Unicode、重复字段、零值、数值边界和自定义方法建立兼容测试。

## 其他编码格式

| 格式/包 | 适用场景 | 注意 |
| --- | --- | --- |
| encoding/csv | 表格交换 | 用 Reader/Writer 处理引号和逗号，检查 Error |
| encoding/xml | XML 协议 | 命名空间、实体与大小限制 |
| encoding/binary | 定长二进制协议 | 明确大小端与位宽，勿直接依赖结构体内存布局 |
| encoding/base64、hex | 文本化二进制 | 是编码，不是加密 |
| encoding/gob | Go 类型数据交换 | 注意类型演化与不可信输入，不作通用跨语言协议 |

## 字节序与协议设计

二进制协议至少定义 magic、版本、长度、字段位宽和校验策略。先验证长度再分配内存，防止恶意超大尺寸。网络字节序常用大端，但以具体协议为准。结构体 padding 和平台差异使 unsafe 内存拷贝不适合作为可移植序列化。

练习：为 API 写缺失字段、零值、null、未知字段、尾随 JSON 和超大整数测试；用 binary.BigEndian 编码一个 uint32 并验证解码。

延伸阅读：[encoding/json](https://pkg.go.dev/encoding/json)、[encoding/binary](https://pkg.go.dev/encoding/binary)、[encoding/csv](https://pkg.go.dev/encoding/csv)。

# Agent 项目规范

## 技术栈

| 类别 | 选型 |
|------|------|
| 语言 | Go 1.26 |
| API 定义 | Proto3 + `google.api.http` 注解 |
| 构建工具 | Buf |
| RPC | gRPC + grpc-gateway（HTTP 反向代理） |
| ORM | GORM + GORM Gen（代码生成） |
| 数据库 | PostgreSQL |
| 日志 | 标准库 `log/slog` |
| 配置 | YAML |

## 目录职责

```
api/v1/                proto 接口定义，按业务分包
internal/
├── api/v1/            buf 生成的代码（.pb.go / _grpc.pb.go / .pb.gw.go），勿手动编辑
├── config/            配置加载，提供类型安全的配置结构体
├── consts/            全局常量、枚举
├── dto/               数据传输对象（非 proto 生成的 DTO）
├── middleware/        gRPC 拦截器 / HTTP 中间件
├── model/             GORM Gen 生成的数据模型，勿手动编辑
├── mq/                消息队列
├── pkg/               通用工具包
├── repository/
│   ├── db/dao/        GORM Gen 生成的 DAO 查询代码，勿手动编辑
│   ├── db/            手写的数据库查询扩展
│   └── cache/         缓存层
└── service/           业务逻辑实现，直接调用 repository
cmd/
├── main.go            入口，编排启动流程
├── router.go          服务注册，启动 gRPC + HTTP 网关
└── gen/               GORM Gen 代码生成器入口
config/
└── app.yaml           运行时配置
```

## 代码生成工作流

### proto → HTTP + gRPC

```
api/v1/*.proto
  → buf generate               # make pb
  → internal/api/v1/           # .pb.go + _grpc.pb.go + .pb.gw.go
```

### 数据库表 → model + dao

```
PostgreSQL 表结构
  → go run cmd/gen/            # make db
  → internal/model/            # 实体 struct
  → internal/repository/db/dao/ # 类型安全 CRUD 方法
```

## 编程规范

### 日志

使用标准库 `log/slog`，不得使用第三方日志库。

**格式要求：**

```go
slog.InfoContext(ctx, "[servicename] [method]", "key", value, "key2", value2)
```

- `servicename` — 服务名，如 `hello`、`user`、`order`
- `method` — 方法名，如 `SayHello`、`CreateUser`
- 参数使用 slog 原生结构化 key-value 交替传入，不要手动拼接字符串

**示例：**

```go
// 正确 — slog 原生 key-value
slog.InfoContext(ctx, "[hello] [SayHello]", "name", req.Name, "user_id", userID)

// 错误 — 不要手动拼接
// slog.InfoContext(ctx, "[hello] [SayHello] name="+req.Name)

// 错误 — 不要用第三方日志库
// logrus.Info(...)
// zap.L().Info(...)

// 错误 — 不要用 log.Printf
// log.Printf("SayHello: %s", req.Name)

// 错误 — 不要用固定字符串无 key
// slog.Info("say hello")
```

### Context 规范

所有业务函数和外部调用（数据库、RPC、MQ 等）**必须将 `ctx context.Context` 作为第一个参数**。

- 纯内部工具函数或私有辅助函数可以例外
- Service 结构体持有 `ctx` 字段，`NewService()` 中赋值为 `context.Background()`，作为基础上下文

**示例：**

```go
type Service struct {
    ctx context.Context
    pb.UnimplementedHelloServiceServer
}

func NewService() *Service {
    return &Service{
        ctx: context.Background(),
    }
}

// 业务方法必须接收 ctx 参数（gRPC 框架传入的请求上下文）
func (s *Service) SayHello(ctx context.Context, req *pb.HelloReq) (*pb.HelloResp, error) {
    slog.InfoContext(ctx, "[hello] [SayHello]", "name", req.Name)
    // 调用外部服务时传递 ctx
    user, err := repository.GetDB().WithContext(ctx).Find(...)
    if err != nil {
        return nil, err
    }
    ...
}
```

### 错误处理

- 业务错误使用 `status.Errorf(codes.InvalidArgument, ...)` 返回 gRPC 错误码
- 不要在 service 层吞掉错误，无法处理时向上抛
- 数据库错误统一由 repository 层处理

### 命名

- proto 文件：蛇形命名 `hello.proto`
- proto message：驼峰 `HelloReq` / `HelloResp`
- 目录：`service/hello/`、`service/user/`
- 接口方法：`SayHello`、`CreateUser`、`GetUserByID`

### HTTP 路由

所有 HTTP 路径必须以 `v1/service-name` 为前缀，多级路径用 `/` 分隔，多个单词用 `-` 连接。

**格式：**
```
/v1/{service-name}/{action}
```

**示例：**
```
/v1/hello/filter-models
/v1/user/model/list
/v1/user/model/add
/v1/user/model/update
/v1/order/list
/v1/order/get-by-id
```

**proto 中定义：**
```protobuf
service UserService {
  rpc ListUser(ListUserReq) returns (ListUserResp) {
    option (google.api.http) = {
      get: "/v1/user/model/list"
    };
  }
  rpc CreateUser(CreateUserReq) returns (CreateUserResp) {
    option (google.api.http) = {
      post: "/v1/user/model/add"
      body: "*"
    };
  }
}
```

### 新增服务步骤

1. `api/v1/` 下创建 `xxx/xxx.proto`，定义 service + message + HTTP 注解
2. `make generate` 生成代码
3. `internal/service/xxx/` 下实现业务逻辑，嵌入 `UnimplementedXxxServer`
4. `cmd/router.go` 中注册 gRPC + HTTP 网关
5. （可选）`internal/repository/db/` 下编写自定义查询

# go-project-template

基于 proto + buf + grpc-gateway 的 Go HTTP + gRPC 服务骨架。

## 技术栈

- **Go 1.26** — 语言
- **Proto** — API 定义
- **Buf** — protobuf 构建工具
- **grpc-gateway** — HTTP 反向代理，一个 proto 同时生成 HTTP + gRPC
- **gRPC** — RPC 框架
- **OpenAPI** — 自动生成 Swagger 文档

## 快速开始

### 安装依赖

```bash
# protobuf 工具链
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# buf 构建工具
# Windows: scoop install buf
# macOS:   brew install buf
# Linux:   go install github.com/bufbuild/buf/cmd/buf@latest

# gRPC 调试工具（可选）
scoop install grpcurl   # Windows
brew install grpcurl    # macOS
```

### 更改配置
```bash
cp config/app.template.yaml config/app.yaml   # 首次使用
```

### 生成并运行

```bash
# 1. 生成 proto 代码（修改 proto 后需要重新生成）
make pb

# 2. 生成数据库 model + dao（需先配好 config/app.yaml 并确保数据库可连）
make db

# 3. 编译
make build

# 4. 运行
make run

# gRPC 服务监听 :9090
# HTTP 服务监听 :8080
```

## 测试

```bash
# HTTP 测试
curl "http://localhost:8080/v1/hello?name=world"

# gRPC 测试（需安装 grpcurl）
grpcurl -plaintext -d '{"name": "world"}' localhost:9090 hello.HelloService/SayHello
```

## 目录结构

```
├── api/v1/               proto 接口定义
├── cmd/
│   ├── main.go           入口，编排启动流程
│   ├── router.go         服务注册，启动 gRPC + HTTP 网关
│   ├── router_gen.go     自动生成的服务注册代码（勿手动编辑）
│   ├── gormgen/            GORM Gen 代码生成器入口
│   └── register-gen/     服务注册代码生成器
├── config/               运行时配置
├── internal/
│   ├── api/v1/           buf 生成的 proto 代码（勿手动编辑）
│   ├── config/           配置加载
│   ├── consts/           全局常量、枚举（预留）
│   ├── dto/              数据传输对象（预留）
│   ├── middleware/        gRPC/HTTP 中间件（预留）
│   ├── model/            GORM Gen 数据模型 + 自定义查询接口
│   ├── mq/               消息队列（预留）
│   ├── pkg/              通用工具包（预留，含 i18n 等）
│   ├── repository/       数据访问层（DB + Cache）
│   └── service/          业务逻辑层
├── locales/              国际化 / 自定义错误码消息
└── scripts/              脚本（SQL、Shell 等，预留）
```

## 添加新服务

1. 在 `api/v1/` 下创建新的 proto 文件
2. 定义 service 和 message，添加 `google.api.http` 注解
3. 运行 `make pb` 生成代码（proto → pb + 自动注册）
4. 在 `internal/service/` 下实现业务逻辑

> 服务注册由 `cmd/register-gen` 自动完成，`make pb` 已包含此步骤。

## Proto 规范

- 使用 `google.api.http` 注解定义 HTTP 路由
- 路径格式：`/v1/{service-name}/{action}`，多词用 `-` 连接（如 `/v1/user/model/list`）
- 所有 proto 文件放在 `api/v1/` 下，按业务分包
- 生成代码使用 `source_relative` 路径风格

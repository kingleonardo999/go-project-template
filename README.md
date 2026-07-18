# Agent

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

## 添加新服务

1. 在 `api/v1/` 下创建新的 proto 文件
2. 定义 service 和 message，添加 `google.api.http` 注解
3. 运行 `make pb` 生成代码
4. 在 `internal/service/` 下实现业务逻辑
5. 在 `cmd/router.go` 中注册新服务（gRPC + HTTP 网关）

## Proto 规范

- 使用 `google.api.http` 注解定义 HTTP 路由
- 路径格式：`/v1/{service-name}/{action}`，多词用 `-` 连接（如 `/v1/user/model/list`）
- 所有 proto 文件放在 `api/v1/` 下，按业务分包
- 生成代码使用 `source_relative` 路径风格

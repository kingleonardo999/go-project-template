# 项目 Makefile
.PHONY: pb build run db clean fmt lint test mod

# 拉取依赖
mod:
	go mod tidy

# 生成 proto 代码并自动注册服务
pb:
	buf generate
	go run ./cmd/register-gen/

# 生成数据库 DAO 代码（连接数据库后自动生成 model + dao）
db:
	go run ./cmd/gormgen/

# 格式化代码
fmt:
	go fmt ./...

# 静态检查
lint:
	go vet ./...

# 编译项目（先拉取依赖 → 格式化 → 检查）
build: mod fmt lint
	go build -o bin/agent ./cmd/

# 运行测试
test:
	go test ./... -v -count=1

# 运行服务（先拉取依赖）
run: mod
	go run ./cmd/

# 清理
clean:
	rm -rf bin/
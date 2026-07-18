# 项目 Makefile
.PHONY: pb build run db clean fmt lint test

# 生成 proto 代码
pb:
	buf generate

# 生成数据库 DAO 代码（连接数据库后自动生成 model + dao）
db:
	go run ./cmd/gen/

# 格式化代码
fmt:
	go fmt ./...

# 静态检查
lint:
	go vet ./...

# 编译项目（先格式化 + 检查）
build: fmt lint
	go build -o bin/agent ./cmd/

# 运行测试
test:
	go test ./... -v -count=1

# 运行服务
run:
	go run ./cmd/

# 清理
clean:
	rm -rf bin/
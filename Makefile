.PHONY: help proto build test clean docker-up docker-down

# 默认目标：显示帮助信息
help:
	@echo "CheYa v3.0 开发工具集"
	@echo ""
	@echo "可用命令:"
	@echo "  make proto         - 编译所有 .proto 文件"
	@echo "  make build         - 构建所有微服务"
	@echo "  make test          - 运行测试"
	@echo "  make docker-up     - 启动 Docker 基础设施"
	@echo "  make docker-down   - 停止 Docker 基础设施"
	@echo "  make clean         - 清理生成的文件"

# 编译 Protobuf 文件
proto:
	@echo "编译 Protobuf 文件..."
	@protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/**/*.proto
	@echo "✅ Protobuf 编译完成"

# 构建所有微服务
build:
	@echo "构建微服务..."
	@go build -o bin/gateway ./apps/gateway
	@go build -o bin/vehicle ./apps/vehicle
	@go build -o bin/telemetry ./apps/telemetry
	@go build -o bin/ai-agent ./apps/ai-agent
	@go build -o bin/auth ./apps/auth
	@echo "✅ 构建完成"

# 运行测试
test:
	@echo "运行测试..."
	@go test -v ./...
	@echo "✅ 测试完成"

# 启动 Docker 基础设施
docker-up:
	@echo "启动基础设施容器..."
	@docker-compose up -d
	@echo "✅ 容器已启动"
	@echo "访问 Kafka UI: http://localhost:8090"

# 停止 Docker 基础设施
docker-down:
	@echo "停止基础设施容器..."
	@docker-compose down
	@echo "✅ 容器已停止"

# 清理生成的文件
clean:
	@echo "清理生成的文件..."
	@rm -rf bin/
	@echo "✅ 清理完成"

# 安装依赖
deps:
	@echo "安装 Go 依赖..."
	@go mod download
	@go mod tidy
	@echo "✅ 依赖安装完成"


#!/bin/bash

echo "=========================================="
echo "  CheYa v3.0 开发环境验证报告"
echo "=========================================="
echo ""

# 检查 Go
echo "【1/7】检查 Go 环境..."
if command -v go &> /dev/null; then
    GO_VERSION=$(go version)
    echo "✅ $GO_VERSION"
    echo "   GOPATH: $(go env GOPATH)"
    echo "   GOPROXY: $(go env GOPROXY)"
else
    echo "❌ Go 未安装"
fi
echo ""

# 检查 protoc
echo "【2/7】检查 Protocol Buffers 编译器..."
if command -v protoc &> /dev/null; then
    PROTOC_VERSION=$(protoc --version)
    echo "✅ $PROTOC_VERSION"
else
    echo "❌ protoc 未安装"
fi
echo ""

# 检查 protoc Go 插件
echo "【3/7】检查 protoc Go 插件..."
if [ -f "$(go env GOPATH)/bin/protoc-gen-go" ]; then
    echo "✅ protoc-gen-go 已安装"
else
    echo "❌ protoc-gen-go 未安装"
fi
if [ -f "$(go env GOPATH)/bin/protoc-gen-go-grpc" ]; then
    echo "✅ protoc-gen-go-grpc 已安装"
else
    echo "❌ protoc-gen-go-grpc 未安装"
fi
echo ""

# 检查 Ent
echo "【4/7】检查 Ent ORM..."
if [ -f "$(go env GOPATH)/bin/ent" ]; then
    echo "✅ Ent 代码生成器已安装"
else
    echo "❌ Ent 未安装"
fi
echo ""

# 检查 grpcurl
echo "【5/7】检查 grpcurl..."
if [ -f "$(go env GOPATH)/bin/grpcurl" ]; then
    echo "✅ grpcurl 已安装"
else
    echo "❌ grpcurl 未安装"
fi
echo ""

# 检查 Docker
echo "【6/7】检查 Docker..."
if command -v docker &> /dev/null; then
    DOCKER_VERSION=$(docker --version)
    echo "✅ $DOCKER_VERSION"
    if docker ps &> /dev/null; then
        echo "✅ Docker 守护进程正在运行"
        RUNNING_CONTAINERS=$(docker ps --format "{{.Names}}" | grep -c "cheya" || echo "0")
        echo "   运行中的 CheYa 容器数: $RUNNING_CONTAINERS"
    else
        echo "⚠️  Docker 已安装但守护进程未运行"
        echo "   提示: 在 OrbStack 环境中，请确保宿主机 Docker 服务已启动"
    fi
else
    echo "❌ Docker 未安装"
fi
echo ""

# 检查项目结构
echo "【7/7】检查项目结构..."
if [ -f "go.mod" ]; then
    echo "✅ go.mod 已初始化"
fi
if [ -f "docker-compose.yml" ]; then
    echo "✅ docker-compose.yml 已创建"
fi
if [ -f "Makefile" ]; then
    echo "✅ Makefile 已创建"
fi
if [ -d "api" ]; then
    echo "✅ api/ 目录已创建"
fi
if [ -d "apps" ]; then
    echo "✅ apps/ 目录已创建"
fi
if [ -d "pkg" ]; then
    echo "✅ pkg/ 目录已创建"
fi
echo ""

echo "=========================================="
echo "  验证完成！"
echo "=========================================="
echo ""
echo "📋 下一步操作建议："
echo ""
echo "1. 如需启动基础设施（需 Docker 守护进程）："
echo "   make docker-up"
echo ""
echo "2. 开始定义 Protobuf 接口："
echo "   编辑 api/vehicle/v1/vehicle.proto"
echo ""
echo "3. 查看完整命令列表："
echo "   make help"
echo ""
echo "🎉 准备就绪，开始编写 CheYa v3.0！"

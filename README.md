# CheYa v3.0 - 云原生车联网平台

> 基于 Go + gRPC + Kubernetes 的高性能分布式车辆监控系统

## 📋 项目概述

CheYa v3.0 是一个云原生微服务架构的车联网平台，采用以下技术栈：

**后端技术栈**:
- **后端框架**: Go 1.21+
- **RPC 通信**: gRPC + Protocol Buffers
- **数据库**: PostgreSQL 15
- **缓存**: Redis 7
- **消息队列**: Kafka 3.5
- **ORM**: Ent Framework
- **容器编排**: Kubernetes + Docker

**前端技术栈**:
- **框架**: React 18 + TypeScript
- **路由**: React Router v6
- **样式**: Tailwind CSS
- **构建工具**: Vite
- **代码规范**: ESLint + Airbnb Style Guide

## 🏗️ 架构设计

### 微服务列表

1. **API Gateway** (`apps/gateway`) - 统一入口，路由转发
2. **Vehicle Service** (`apps/vehicle`) - 车辆主数据管理
3. **Telemetry Service** (`apps/telemetry`) - 实时遥测数据处理
4. **AI Agent Service** (`apps/ai-agent`) - 智能分析与预警
5. **Auth Service** (`apps/auth`) - 认证与鉴权

### 目录结构

```
cheya/
├── api/                    # Protobuf 定义
│   ├── vehicle/v1/
│   ├── telemetry/v1/
│   ├── ai/v1/
│   └── auth/v1/
├── apps/                   # 微服务应用
│   ├── gateway/
│   ├── vehicle/
│   ├── telemetry/
│   ├── ai-agent/
│   └── auth/
├── pkg/                    # 共享工具库
│   ├── logger/
│   ├── config/
│   ├── errors/
│   └── middleware/
├── front-dashboard/        # 前端管理后台 (React + TypeScript)
│   ├── src/
│   ├── public/
│   ├── package.json
│   └── vite.config.ts
├── deploy/                 # 部署配置
│   ├── k8s/
│   └── docker/
├── scripts/                # 工具脚本
├── docs/                   # 文档
├── docker-compose.yml      # 本地开发环境
├── Makefile               # 构建工具
└── go.mod                 # Go 依赖管理
```

## 🚀 快速开始

### 1. 环境要求

**后端**:
- Go 1.21+
- Protocol Buffers Compiler (protoc)
- Docker & Docker Compose
- Make

**前端**:
- Node.js 18+
- npm 或 pnpm

### 2. 启动基础设施

```bash
# 启动 PostgreSQL, Redis, Kafka
make docker-up

# 查看容器状态
docker ps
```

### 3. 编译 Protobuf

```bash
make proto
```

### 4. 构建微服务

```bash
make build
```

### 5. 启动前端开发服务器

```bash
cd front-dashboard
npm install
npm run dev
# 访问 http://localhost:5173
```

### 6. 运行测试

```bash
make test
```

## 🔧 开发工具

### Makefile 命令

| 命令 | 说明 |
|------|------|
| `make help` | 显示所有可用命令 |
| `make proto` | 编译 Protobuf 文件 |
| `make build` | 构建所有微服务 |
| `make test` | 运行单元测试 |
| `make docker-up` | 启动基础设施容器 |
| `make docker-down` | 停止基础设施容器 |
| `make clean` | 清理生成的文件 |
| `make deps` | 安装 Go 依赖 |

### gRPC 测试工具

使用 `grpcurl` 测试 gRPC 接口：

```bash
# 列出服务
grpcurl -plaintext localhost:50051 list

# 调用方法
grpcurl -plaintext -d '{"vehicle_id": "V001"}' \
  localhost:50051 cheya.vehicle.v1.VehicleService/GetVehicle
```

## 📡 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| 前端开发服务器 | 5173 | Vite Dev Server |
| API Gateway | 8080 | HTTP/REST API |
| gRPC Services | 50051+ | gRPC 通信 |
| PostgreSQL | 5432 | 数据库 |
| Redis | 6379 | 缓存 |
| Kafka | 9092 | 消息队列 |
| Kafka UI | 8090 | Kafka 管理界面 |

## 📚 技术文档

详细的技术设计请参考：
- [技术白皮书](../truck-monitor/TECHNICAL_WHITEPAPER_V3.md)
- [API 文档](docs/api/)
- [部署指南](docs/deployment.md)

## 🤝 开发规范

1. **代码风格**: 遵循 Go 官方代码规范
2. **提交信息**: 使用约定式提交 (Conventional Commits)
3. **测试覆盖**: 单元测试覆盖率 > 80%
4. **文档**: 所有公共 API 必须有注释

## 📝 下一步计划

### Phase 1: RPC 契约设计
- [ ] 定义 `vehicle.proto`
- [ ] 定义 `telemetry.proto`
- [ ] 定义 `ai.proto`
- [ ] 定义 `auth.proto`

### Phase 2: 核心服务开发
- [ ] 实现 Vehicle Service
- [ ] 实现 Telemetry Service
- [ ] 实现 AI Agent Service
- [ ] 实现 Auth Service

### Phase 3: 集成与测试
- [ ] 服务间通信测试
- [ ] 性能压测
- [ ] 安全测试

## 📄 许可证

Copyright © 2025 CheYa Team

---

**架构师**: 薛文涛  
**版本**: v3.0  
**更新时间**: 2025-11-18


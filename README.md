# 通讯录管理系统

这是一个基于 Go 语言的即时通讯（IM）应用雏形。本项目并非完整聊天软件，而是一个教学与实践项目。它重点模拟了通讯应用中“查看最常访问好友”的场景，并深入实践了 **雪花算法 (Snowflake)** 和 **LRU (Least Recently Used) 缓存策略**。

## 📚 1. 项目简介 (Introduction)

本系统主要解决了以下核心问题：

1.  **全局唯一 ID 生成**：在分布式场景下，为每一条消息、每一个好友关系生成趋势递增且不重复的 64 位 ID，能够承受高并发请求（> 1000 QPS）。
2.  **高效缓存管理**：在前后端分别实现 LRU 缓存策略。
    *   **前端**：使用手写哈希表+双向链表缓存 `uid` 到内部 `id` 的映射，减少频繁查找开销。
    *   **后端**：优化“最常访问好友”的查询速度，减少对数据库的击穿。
3.  **完整架构实践**：从数据库设计、后端 API 实现到前端交互展现，提供了一个完整的全栈开发案例。

## 🏗️ 2. 系统架构 (Architecture)

### 2.1 技术栈

| 模块 | 技术选型 | 说明 |
| :--- | :--- | :--- |
| **Backend** | **Go 1.25+** | 高性能后端语言 |
| **Framework** | **Gin、gorm** | 轻量级 Web 框架与数据库ORM框架 |
| **Database** | **PostgreSQL** | 关系型数据库 |
| **Testing** | **Python / Node.js** | 压测脚本 (`id_bench`) 与算法基准测试 (`lru_bench`) |
| **Frontend** | **Vanilla JS** | 原生 JavaScript (无框架)，手写 LRU 算法 |

### 2.2 目录结构说明

```
vvechat/
├── backend/            # Go 后端代码
│   ├── cmd/            # 程序入口 (main.go)
│   ├── conf/           # 配置文件 (config.yaml)
│   ├── internal/       # 业务逻辑 (handler, service, model)
│   └── pkg/            # 公共组件 (snowflake, jwt, response)
├── frontend/           # 前端代码 (纯静态文件)
│   └── lru_cache.js    # 核心：前端手写 LRU 算法实现
├── test/               # 测试套件
│   ├── id_bench/       # Snowflake ID 生成器压测 (Python)
│   ├── lru_bench/      # LRU 算法基准测试 (Node.js/Python)
│   └── api_test/       # API 接口集成测试
└── docs/               # 项目文档
```

## 🛠️ 3. 部署与运行 (Deployment)

### 3.1 环境要求

*   **Go**: 1.25.5 或更高版本
*   **PostgreSQL**: 运行在 `localhost:5432`
*   **Python**: 3.8+ (用于运行测试脚本)
*   **Node.js**: (可选，用于运行前端基准测试)

### 3.2 启动步骤

#### 第一步：配置数据库

确保 PostgreSQL 和 Redis 正在运行。
默认配置位于 `backend/conf/config.yaml`：

```yaml
postgres:
  dsn: "host=localhost user=lojes dbname=wechat port=5432 sslmode=disable TimeZone=Asia/Shanghai"
redis:
  addr: 127.0.0.1:6379
```

请根据你的本地环境修改 `dsn` 中的用户名 (`user`) 和数据库名 (`dbname`)。

#### 第二步：启动后端

```bash
cd backend
# 下载依赖
go mod tidy
# 运行服务
go run cmd/main.go
```
服务默认监听在 `8080` 端口。

#### 第三步：运行前端

该项目前端为纯静态文件，可以使用任意 HTTP Server 启动，或者直接在 VS Code 中使用 "Live Server" 插件打开 `frontend/index.html`。

如果您安装了 Python，也可以快速启动一个静态服务器：
```bash
cd frontend
python -m http.server 3000
```
然后在浏览器访问 `http://localhost:3000`。

## 🧪 4. 测试套件 (Testing)

项目包含一套完善的测试工具，用于验证算法性能。

### 4.1 雪花算法压测 (ID Benchmark)

验证分布式 ID 生成器的性能与唯一性。

```bash
cd test/id_bench
# 安装依赖 (建议在 venv 中)
pip install matplotlib
# 运行压测
python runner.py --duration 10 --workers 4
```

### 4.2 LRU 性能测试

验证前端手写 LRU 缓存的性能。

```bash
cd test/lru_bench
# 运行 Node.js 基准测试
node benchmark_lru.js
```

## 📝 5. 关于项目 (About)

本项目是 **数据结构课程设计** 的产出成果。

*   **核心贡献**:
    *   将理论层面的“双向链表+哈希表”应用于实际的 LRU 缓存实现。
    *   通过位运算实现高效的雪花算法，解决了分布式 ID 生成难题。
    *   构建了可视化的基准测试工具，直观展示算法性能。


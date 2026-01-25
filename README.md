# AI Skeleton - AI 全栈脚手架

面向 AI 辅助开发的全栈项目脚手架，专为中文开发者设计。

## 一键安装

```bash
curl -fsSL https://raw.githubusercontent.com/richer421/ai_skeleton/main/install.sh | bash
```

安装完成后：

```bash
# 创建新项目
ai-skeleton init my_project

# 启动项目
cd my_project
make backend-dev   # 后端：http://localhost:8080
make frontend-dev  # 前端：http://localhost:5173
```

## 快速开始（手动安装）

如果你想手动安装 CLI 工具：

```bash
# 1. 克隆脚手架
git clone https://github.com/richer421/ai_skeleton.git
cd ai_skeleton

# 2. 构建 CLI 工具
cd cli && go build -o ai-skeleton main.go && cd ..

# 3. 初始化新项目（自动安装依赖）
./cli/ai-skeleton init my_project

# 4. 启动项目
cd my_project
make backend-dev   # 后端：http://localhost:8080
make frontend-dev  # 前端：http://localhost:5173
```

CLI 会自动：
- ✓ 检查环境（Go、npm）
- ✓ 复制模板并替换项目信息
- ✓ 安装依赖（Air、Swagger、npm packages）
- ✓ 确保前后端能正常启动

## 技术栈

### 后端
- Go 1.21
- Gin - HTTP 框架
- GORM + Gen - ORM 和代码生成
- Swagger - API 文档
- MCP 协议支持

### 前端
- React 18
- Vite - 构建工具
- Ant Design v5 - UI 组件库
- TypeScript
- Axios - HTTP 客户端

## 详细使用说明

### 方式一：使用 CLI 工具初始化（推荐）

**1. 构建 CLI 工具**

```bash
cd cli
go build -o ai-skeleton main.go
```

**2. 初始化新项目**

```bash
# 在脚手架根目录运行
./cli/ai-skeleton init

# 或指定项目名称和配置
./cli/ai-skeleton init my_project --desc "我的项目" --module "github.com/myname/my_project"
```

CLI 工具会自动：
- ✓ 检查环境依赖（Go、npm）
- ✓ 收集项目信息（交互式输入）
- ✓ 复制模板文件并替换占位符
- ✓ 安装依赖（Air、Swagger、npm packages）
- ✓ 确保前后端都能正常启动

**3. 启动项目**

```bash
cd my_project

# 启动后端
make backend-dev

# 启动前端（另开终端）
make frontend-dev
```

### 方式二：手动安装

### 0. 安装依赖

**后端依赖：**
```bash
# 安装 Air（热重载工具）
go install github.com/air-verse/air@latest
```

### 1. 克隆项目

```bash
git clone <repo>
cd ai_skeleton
```

### 2. 启动后端

```bash
cd backend
cp .env.example .env
make backend-dev
```

后端服务将启动在 `http://localhost:8080`，支持热重载（修改代码自动重启）

### 3. 启动前端

```bash
cd frontend
npm install
make frontend-dev
```

前端服务将启动在 `http://localhost:5173`

### 4. 验证

访问 `http://localhost:5173` 查看前端页面，页面会自动调用后端健康检查接口。

**测试 MCP 接口：**

```bash
# 列出所有 MCP 工具
curl http://localhost:8080/api/v1/mcp/tools

# 执行健康检查工具
curl -X POST http://localhost:8080/api/v1/mcp/execute \
  -H "Content-Type: application/json" \
  -d '{"tool":"health_check","params":{}}'
```

## CLI 工具命令

```bash
# 查看帮助
ai-skeleton --help

# 初始化项目
ai-skeleton init [项目名] [选项]

# 生成代码
ai-skeleton generate service [服务名]

# 配置管理
ai-skeleton config generate
ai-skeleton config validate
```

详见 [CLI 工具使用文档](./requirements/20260125-cli-init-tool.md)

## 可用命令

```bash
make help          # 显示所有可用命令
make backend-dev   # 启动后端开发服务器
make frontend-dev  # 启动前端开发服务器
make gen-swagger   # 生成 Swagger 文档
make gen-sql       # 生成 Gen-GORM 代码
```

## 项目结构

详见 [CLAUDE.md](./CLAUDE.md)

## MCP 协议支持

项目内置 MCP (Model Context Protocol) 协议支持，可以将后端功能暴露给 AI 使用。

**MCP API：**
- `GET /api/v1/mcp/tools` - 列出所有工具
- `POST /api/v1/mcp/execute` - 执行工具

**已注册工具：**
- `health_check` - 系统健康检查

**统一注册层：**

所有 MCP 工具通过 `mcp.RegisterAllTools()` 统一注册，只需：
1. 在 `internal/mcp/tools.go` 的 `Services` 结构体中添加服务
2. 在 `RegisterAllTools()` 中添加注册调用
3. 实现具体的注册函数

详见 [CLAUDE.md](./CLAUDE.md) 中的 MCP 协议支持章节。

## 开发规范

详见 [CLAUDE.md](./CLAUDE.md)

## License

MIT

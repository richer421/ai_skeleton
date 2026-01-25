# AI Skeleton - AI Coding 全栈脚手架

## 项目概述

这是一个面向 AI 辅助开发的全栈项目脚手架，专为中文开发者设计。项目采用前后端分离架构，提供开箱即用的开发环境和最佳实践。

## 技术栈

### 后端 (Backend)
- **语言**: Go
- **HTTP 框架**: Gin - 高性能 Web 框架
- **ORM**: GORM + Gen - 类型安全的 ORM 代码生成器
- **API 文档**: Swagger - 自动生成 API 文档
- **协议支持**: MCP (Model Context Protocol) - AI 模型上下文协议
- **开发工具**: Air (热重载)
- **测试**: Go 原生测试框架

### 前端 (Frontend)
- **框架**: React
- **构建工具**: Vite - 快速的前端构建工具
- **UI 组件库**: Ant Design v5 (antd)
- **包管理**: npm

### 部署
- **容器化**: Dockerfile
- **部署方式**: 前端和后端分别打包为独立的 Docker 镜像

## 项目结构

```
ai_skeleton/
├── backend/              # Go 后端服务（简单三层架构）
│   ├── internal/        # 内部代码
│   │   ├── http/        # HTTP 相关
│   │   │   ├── api/     # API 处理器（Gin API 函数）
│   │   │   ├── router/  # 路由配置
│   │   │   └── middleware/ # 中间件
│   │   ├── service/     # 业务逻辑层
│   │   │   └── health/  # 健康检查服务（每个服务一个目录）
│   │   ├── repository/  # 数据访问层（Gen-GORM）
│   │   ├── model/       # 数据模型
│   │   ├── mcp/         # MCP 协议适配
│   │   ├── config/      # 配置管理
│   │   ├── common/      # 公共代码（错误定义等）
│   │   └── testutil/    # 测试工具（mock、factory、assert）
│   ├── cmd/             # 应用入口
│   │   ├── server/      # 服务启动
│   │   └── gen/         # 代码生成工具
│   ├── tmp/             # Air 热重载临时文件
│   └── bin/             # 编译产物
├── frontend/            # 前端应用
│   ├── src/             # 源代码
│   ├── dist/            # 构建产物
│   └── node_modules/    # 依赖包
├── requirements/        # 项目需求文档
├── scripts/             # 脚本工具
├── Makefile             # 项目命令集合
└── Dockerfile           # Docker 构建文件
```

## 架构设计

### 简单三层架构

项目采用简单清晰的三层架构，便于 AI 理解和开发：

**目录组织原则：**
- HTTP 相关代码（api、router、middleware）统一放在 `internal/http/` 目录
- 每个服务（service）使用独立目录，如 `internal/service/health/`
- 所有内部代码放在 `internal/` 目录，遵循 Go 标准项目布局

#### 1. API 层（HTTP 处理层）
- 位置：`internal/http/api/`
- 接收 HTTP 请求，参数验证
- 调用 Service 层处理业务逻辑
- 返回 HTTP 响应
- 添加 Swagger 注释

#### 2. Service 层（业务逻辑层）
- 位置：`internal/service/{服务名}/`
- 每个服务使用独立目录
- 实现核心业务逻辑
- 调用 Repository 层进行数据操作
- 处理业务规则和验证
- 事务管理

#### 3. Repository 层（数据访问层）
- 位置：`internal/repository/`
- 封装数据库操作
- 使用 Gen-GORM 生成的代码
- 提供 CRUD 接口
- 只负责数据持久化，不包含业务逻辑
- 数据持久化

### 其他模块

- **model/**: 数据模型定义（数据库表结构）
- **mcp/**: MCP 协议适配模块
- **middleware/**: 中间件（认证、日志、CORS 等）
- **router/**: 路由配置和注册
- **config/**: 配置文件加载和管理

### MCP 协议支持

后端接口天然兼容 MCP (Model Context Protocol) 协议，MCP 适配模块位于 `internal/mcp/`：

**核心功能：**
- **MCP 适配层**: 将 Service 层功能转换为 MCP 协议格式
- **工具注册**: 支持将后端功能注册为 MCP 工具
- **上下文管理**: 管理 AI 模型的上下文信息
- **协议转换**: 处理 MCP 请求和响应的序列化/反序列化

**MCP API 接口：**
- `GET /api/v1/mcp/tools` - 列出所有已注册的 MCP 工具
- `POST /api/v1/mcp/execute` - 执行指定的 MCP 工具

**已注册的工具：**
- `health_check` - 检查系统健康状态

**统一注册层设计：**

所有 MCP 工具通过 `mcp.RegisterAllTools()` 统一注册，避免在路由层一个一个添加。

```go
// internal/mcp/tools.go
type Services struct {
    HealthService health.HealthService
    // 添加更多服务
}

func RegisterAllTools(adapter MCPAdapter, services *Services) error {
    // 统一注册所有工具
}
```

**如何添加新的 MCP 工具：**

1. 在 `internal/mcp/tools.go` 的 `Services` 结构体中添加新服务
2. 在 `RegisterAllTools()` 函数中调用新的注册函数
3. 在 `internal/mcp/tools.go` 中实现具体的注册函数（如 `registerXxxTool()`）
4. 在 `internal/http/router/router.go` 中初始化服务并添加到 `Services` 结构体
5. 工具会自动暴露给 AI 使用

**示例：**
```go
// 1. 添加服务到 Services 结构体
type Services struct {
    HealthService health.HealthService
    UserService   user.UserService  // 新增
}

// 2. 在 RegisterAllTools 中注册
func RegisterAllTools(adapter MCPAdapter, services *Services) error {
    if err := registerHealthTool(adapter, services.HealthService); err != nil {
        return err
    }
    if err := registerUserTool(adapter, services.UserService); err != nil {
        return err
    }
    return nil
}

// 3. 实现注册函数
func registerUserTool(adapter MCPAdapter, userService user.UserService) error {
    // 定义工具 schema 和 handler
}
```

### 调用链路

```
HTTP 请求 → API → Service → Repository → 数据库
                     ↓
                  MCP 适配
                     ↓
                  AI 调用
```

## 开发指南

### 产品需求开发流程

**⚠️ 重要：这是强制性的开发流程，必须严格遵守**

当用户提出产品需求时，必须按照以下步骤进行：

#### 第一阶段：需求文档编写（必须完成并确认）

1. **生成开发文档**
   - 在 `requirements/` 目录下创建需求文档（Markdown 格式）
   - 文档命名规则：`YYYYMMDD-功能名称.md`

2. **文档必须包含的内容**
   - **需求概述**：简要描述功能目标
   - **产品交互流程**：详细描述用户操作流程和系统响应
   - **页面结构**：列出涉及的页面和组件
   - **API 接口设计**：定义后端接口（路径、方法、参数、响应）
   - **数据模型**：定义数据库表结构
   - **MCP 工具定义**（如适用）：定义 MCP 协议的工具描述

3. **等待用户确认**
   - 文档生成后，**必须停止**，等待用户确认
   - 使用 AskUserQuestion 工具询问用户是否确认文档
   - **只有用户明确确认后，才能进入开发阶段**

#### 第二阶段：代码开发（仅在文档确认后）

4. **开始开发**
   - 严格按照确认的文档进行开发
   - 前端遵循 Ant Design 规范
   - 后端遵循 Go 和 MCP 规范

5. **开发完成后**
   - 运行相关测试
   - 生成 Swagger 文档
   - 等待用户验收

**禁止行为**：
- ❌ 禁止在文档未确认前开始编写代码
- ❌ 禁止跳过需求文档直接开发
- ❌ 禁止在用户未确认的情况下自行修改需求

### 快速开始

**使用 CLI 工具（推荐）：**

```bash
# 一键安装 CLI 工具
curl -fsSL https://raw.githubusercontent.com/richer421/ai_skeleton/main/install.sh | bash

# 初始化新项目
ai-skeleton init my_project

# 进入项目并启动
cd my_project
make backend-dev   # 启动后端
make frontend-dev  # 启动前端
```

**手动安装依赖：**
```bash
# 安装 Air（Go 热重载工具）
go install github.com/air-verse/air@latest

# 安装 Swagger
go install github.com/swaggo/swag/cmd/swag@latest
```

**启动服务：**

1. **启动前端开发**
   ```bash
   make frontend-dev
   ```

2. **启动后端开发（支持热重载）**
   ```bash
   make backend-dev
   ```
   修改代码后会自动重启服务

3. **生成代码**
   ```bash
   make gen-swagger  # 生成 Swagger 文档
   make gen-sql      # 生成 Gen-GORM 代码
   ```

### 常用命令

| 命令 | 说明 |
|------|------|
| `make help` | 显示所有可用命令 |
| `make frontend-dev` | 启动前端开发服务器 |
| `make backend-dev` | 启动后端开发服务器 (Air 热重载) |
| `make gen-swagger` | 生成 Swagger API 文档 |
| `make gen-sql` | 使用 Gen-GORM 生成数据库操作代码 |

## 代码规范

### Go 后端规范（简单三层架构）

#### 通用规范
- 使用 `go fmt` 格式化代码
- 使用 `go vet` 进行静态检查
- 遵循 Go 官方代码规范
- 测试文件命名: `*_test.go`

#### 分层规范

**API 层（HTTP 处理层）**
- 位置：`internal/http/api/`
- 负责接收 HTTP 请求和返回响应
- 进行参数验证和绑定
- 调用 Service 层处理业务逻辑
- 所有 API 接口必须添加 Swagger 注释
- 不包含业务逻辑，只做请求转发
- API 函数直接定义，不使用 Handler 结构体

**Service 层（业务逻辑层）**
- 位置：`internal/service/{服务名}/`
- 实现核心业务逻辑和规则
- 调用 Repository 层进行数据操作
- 处理事务管理
- 可以调用多个 Repository
- 不直接处理 HTTP 请求
- 每次请求创建新的 Service 实例（无状态）

**Repository 层（数据访问层）**
- 位置：`internal/repository/`
- 封装数据库操作
- 使用 Gen-GORM 生成的代码
- 提供 CRUD 接口
- 只负责数据持久化，不包含业务逻辑

**MCP 模块**
- 位置：`internal/mcp/`
- 调用 Service 层功能，转换为 MCP 协议格式
- MCP 工具需要提供 JSON Schema 描述
- 每个工具的 handler 中创建 Service 实例

#### 可测试性规范（AI 自动化测试）

**⚠️ 重要：所有代码必须便于 AI 生成和执行自动化测试**

**注意**：当前项目使用简单的无参数构造函数（如 `health.NewHealthService()`），每次请求创建新实例。以下规范是针对需要依赖注入的复杂服务（如需要 Repository 层）的最佳实践。

**1. 接口抽象 + 依赖注入（针对复杂服务）**
```go
// Service 层必须定义接口
type UserService interface {
    CreateUser(ctx context.Context, req *CreateUserRequest) (*User, error)
    GetUser(ctx context.Context, id int64) (*User, error)
}

// 实现通过构造函数注入依赖
type userService struct {
    repo repository.UserRepository  // 依赖接口，不是具体实现
}

func NewUserService(repo repository.UserRepository) UserService {
    return &userService{repo: repo}
}
```

**2. 表驱动测试（Table-Driven Tests）**
```go
// 所有测试必须使用表驱动模式
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name    string           // 测试用例名称
        input   *CreateUserRequest
        mock    func(*MockRepo)  // mock 设置
        want    *User
        wantErr bool
    }{
        {
            name:  "创建成功",
            input: &CreateUserRequest{Name: "test"},
            mock:  func(m *MockRepo) { m.On("Create").Return(nil) },
            want:  &User{ID: 1, Name: "test"},
        },
        {
            name:    "名称为空",
            input:   &CreateUserRequest{Name: ""},
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 测试逻辑
        })
    }
}
```

**3. 统一错误处理**
```go
// common/errors.go - 定义所有业务错误
var (
    ErrUserNotFound    = errors.New("user not found")
    ErrInvalidInput    = errors.New("invalid input")
    ErrDuplicateUser   = errors.New("user already exists")
)

// 或使用结构化错误
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}
```

**4. 测试辅助工具**
```go
// testutil/factory.go - 测试数据工厂
func NewMockUserRepository() *MockUserRepository {
    return &MockUserRepository{}
}

func CreateTestUser(opts ...func(*User)) *User {
    user := &User{ID: 1, Name: "test"}
    for _, opt := range opts {
        opt(user)
    }
    return user
}

// testutil/assert.go - 断言辅助函数
func AssertEqual(t *testing.T, got, want interface{}) {
    if !reflect.DeepEqual(got, want) {
        t.Errorf("got %v, want %v", got, want)
    }
}
```

**5. 测试覆盖要求**
- 每个 Service 方法必须有对应的测试
- 至少包含：成功场景、失败场景、边界条件
- API 层测试使用 `httptest` 模拟 HTTP 请求
- Repository 层可以使用内存数据库或 mock

**6. Mock 规范**
- 使用 `testify/mock` 或手写 mock
- Mock 对象放在 `testutil/mock/` 目录
- Mock 方法命名与接口一致

**为什么这样设计？**
- ✅ AI 可以轻松识别接口和依赖关系
- ✅ 表驱动测试结构化，AI 容易生成测试用例
- ✅ 统一的错误定义，AI 可以针对每种错误生成测试
- ✅ 测试辅助函数可复用，减少重复代码
- ✅ 依赖注入使得 mock 替换简单

### 前端规范 - Ant Design 开发要求

**⚠️ 重要：前端开发必须严格遵守以下规则**

1. **组件使用原则**
   - 必须优先使用 Ant Design 组件库提供的组件
   - 必须优先使用 antd 组件自身的属性和配置来实现功能
   - 必须遵循 Ant Design 的设计规范和最佳实践

2. **样式开发限制**
   - **禁止自行编写样式来实现用户功能**
   - 如果确实需要自定义样式，必须先向用户确认
   - 未经用户确认，不得进行任何自定义样式开发
   - 优先使用 antd 组件的 `style`、`className`、`size`、`type` 等内置属性

3. **状态管理原则**
   - **优先使用 antd 组件的内置状态管理能力**，避免手动使用 useState
   - 表单状态：使用 `Form.useForm()` 和 `Form.Item`，不要为每个字段单独写 useState
   - 表格状态：使用 `Table` 组件的 `dataSource`、`pagination`、`onChange` 等内置能力
   - 弹窗状态：使用 `Modal` 的受控模式或 `Modal.confirm()` 等静态方法
   - 下拉选择：通过 `Form` 统一管理，而非单独的 useState
   - 只有在 antd 组件无法提供状态管理能力时，才使用 React 的 useState/useReducer

4. **开发流程**
   - 收到需求后，首先查阅 Ant Design 官方文档
   - 确认是否可以通过 antd 组件的属性配置实现
   - 如果 antd 无法满足需求，必须询问用户是否允许自定义样式
   - 只有在用户明确同意后，才能编写自定义样式

5. **代码质量**
   - 使用项目配置的 ESLint 规则
   - 遵循 React Hooks 最佳实践
   - 组件和文件命名保持一致性

### 通用规范
- 提交信息使用中文，清晰描述改动内容
- 重要功能必须包含测试
- 避免提交敏感信息（配置文件中使用占位符或环境变量）

## AI 协作建议

### 给 AI 助手的指引

**⚠️ 最重要：产品需求开发流程**

当用户提出产品需求时：
1. **第一步：生成需求文档**
   - 在 `requirements/` 目录创建详细的需求文档
   - 包含产品交互流程、页面结构、API 设计、数据模型
   - 文档完成后**必须停止**，使用 AskUserQuestion 工具请求用户确认

2. **第二步：等待用户确认**
   - **禁止在文档未确认前编写任何代码**
   - 只有用户明确确认文档后，才能进入开发阶段

3. **第三步：按文档开发**
   - 严格按照确认的文档进行开发
   - 开发完成后等待用户验收

**这个流程是强制性的，不可跳过或简化**

---

1. **理解项目结构**
   - 这是一个前后端分离的全栈项目
   - 后端采用简单三层架构：API → Service → Repository
   - 前端使用 React + Vite + Ant Design v5
   - 所有操作优先使用 Makefile 中定义的命令
   - 后端接口天然兼容 MCP 协议，MCP 适配代码放在 `internal/mcp/` 模块

2. **后端开发流程（三层架构）**
   - **API 层**：接收 HTTP 请求，参数验证，添加 Swagger 注释，直接定义函数无需结构体
   - **Service 层**：实现业务逻辑，调用 Repository，每次请求创建新实例
   - **Repository 层**：使用 Gen-GORM 生成数据库操作代码
   - **MCP 模块**：调用 Service 层，实现 MCP 协议适配
   - 使用 `make gen-swagger` 生成文档
   - 使用 `make gen-sql` 生成数据库代码
   - 调用链路：HTTP → API → Service → Repository → 数据库

3. **前端开发流程 - 关键要求**
   - **第一步：查阅 Ant Design v5 文档**，确认组件能力
   - **第二步：使用 antd 组件属性**实现功能，不要自己写样式
   - **第三步：优先使用 antd 组件的内置状态管理**（如 Form.useForm()）
   - **第四步：如果必须自定义样式**，先用 AskUserQuestion 工具询问用户
   - **第五步：只有用户同意后**，才能编写自定义样式代码
   - 示例：需要调整间距时，优先使用 `<Space>`、`style={{ margin: '...' }}` 等 antd 方式

4. **文件操作原则**
   - 优先编辑现有文件，避免创建不必要的新文件
   - 需求文档必须放在 `requirements/` 目录
   - 不要修改 `.gitignore` 中列出的文件
   - 配置文件使用 `config.yaml` (Viper + YAML)

5. **开发节奏**
   - 收到产品需求 → 生成需求文档 → 等待确认 → 开始开发
   - 每个阶段完成后都要等待用户确认，不要一次性完成所有工作
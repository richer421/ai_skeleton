# CLI 初始化工具需求文档

**创建日期**: 2026-01-25
**需求类型**: 新功能开发

## 需求概述

为 AI Skeleton 脚手架开发一个命令行初始化工具（CLI），支持项目快速创建、差异化渲染（项目元信息替换）、依赖安装和代码生成功能，确保初始化后前后端都能正常启动。

## 产品交互流程

### 1. 项目初始化流程

```bash
# 用户执行初始化命令
./ai-skeleton init [项目名称]

# 或交互式输入
./ai-skeleton init
```

**交互步骤**：

1. **环境检查**
   - 检查 Go 是否安装（`go version`）
   - 检查 npm 是否安装（`npm -v`）
   - 如果任一工具未安装，提示用户安装后退出
   - 显示当前环境版本信息

2. **收集项目信息**（交互式提问）
   - 项目名称（默认：当前目录名）
   - 项目描述（默认：空）
   - 项目版本（默认：1.0.0）
   - Go 模块路径（默认：`github.com/[用户名]/[项目名]`）

3. **生成项目文件**
   - 复制脚手架模板文件到目标目录
   - 替换以下占位符：
     - `ai_skeleton` → 用户输入的项目名称
     - `github.com/richer/ai_skeleton` → 用户输入的 Go 模块路径
     - `config.yaml` 中的 project 部分（name、version、description）
     - `package.json` 中的 name、version、description
     - `README.md` 中的项目名称和描述
     - `CLAUDE.md` 中的项目名称

4. **安装依赖**
   - 后端依赖：
     - 安装 Air（`go install github.com/air-verse/air@latest`）
     - 安装 Swagger（`go install github.com/swaggo/swag/cmd/swag@latest`）
     - 执行 `go mod tidy`
   - 前端依赖：
     - 进入 `frontend/` 目录
     - 执行 `npm install`

5. **生成初始代码**
   - 生成 Swagger 文档（`make gen-swagger`）
   - 显示后续操作提示

6. **完成提示**
   ```
   ✅ 项目初始化完成！

   项目信息：
   - 名称：[项目名称]
   - 描述：[项目描述]
   - Go 模块：[模块路径]

   下一步操作：
   1. cd [项目目录]
   2. 启动后端：make backend-dev
   3. 启动前端：make frontend-dev
   4. 访问：http://localhost:5173
   ```

### 2. 代码生成流程

```bash
# 生成完整的 CRUD 模块
./ai-skeleton generate service [服务名]
```

**交互步骤**：

1. **收集服务信息**
   - 服务名称（如：user、order）
   - 是否生成 API 层（默认：是）
   - 是否注册 MCP 工具（默认：否）

2. **生成代码文件**
   - `internal/service/[服务名]/[服务名]_service.go` - Service 接口和实现
   - `internal/service/[服务名]/[服务名]_service_test.go` - 测试文件
   - `internal/http/api/[服务名].go` - API Handler（如果选择生成）
   - 在 `router.go` 中添加路由注册（如果选择生成 API）
   - 在 `internal/mcp/tools.go` 中添加 MCP 工具注册代码（如果选择 MCP）

3. **完成提示**
   ```
   ✅ 服务代码生成完成！

   生成的文件：
   - internal/service/[服务名]/[服务名]_service.go
   - internal/service/[服务名]/[服务名]_service_test.go
   - internal/http/api/[服务名].go

   下一步操作：
   1. 实现 Service 层的业务逻辑
   2. 添加 Swagger 注释到 API 层
   3. 运行 make gen-swagger 更新文档
   ```

### 3. 配置管理流程

```bash
# 生成配置文件模板
./ai-skeleton config generate

# 验证配置文件
./ai-skeleton config validate
```

## 页面结构

本功能为纯命令行工具，无 UI 页面。

## CLI 命令设计

### 主命令结构

```
ai-skeleton
├── init [项目名]           # 初始化新项目
├── generate (gen)           # 代码生成
│   ├── service [名称]      # 生成服务代码
│   └── api [名称]          # 生成 API 代码
├── config                   # 配置管理
│   ├── generate            # 生成配置模板
│   └── validate            # 验证配置文件
├── version                  # 显示版本信息
└── help                     # 显示帮助信息
```

### 命令参数

#### `init` 命令
```bash
ai-skeleton init [项目名] [选项]

选项：
  --name, -n         项目名称
  --desc, -d         项目描述
  --version, -v      项目版本（默认：1.0.0）
  --module, -m       Go 模块路径
  --skip-deps        跳过依赖安装
  --skip-npm         跳过 npm 依赖安装
  --skip-go          跳过 Go 工具安装
```

#### `generate service` 命令
```bash
ai-skeleton generate service [服务名] [选项]

选项：
  --with-api         生成 API 层（默认）
  --no-api           不生成 API 层
  --with-mcp         注册 MCP 工具
  --no-test          不生成测试文件
```

## 数据模型

### 项目配置模型

```go
// ProjectMeta 项目元信息
type ProjectMeta struct {
    Name        string // 项目名称
    Description string // 项目描述
    Version     string // 项目版本
    Module      string // Go 模块路径
    Author      string // 作者（可选）
}

// InitOptions 初始化选项
type InitOptions struct {
    ProjectMeta
    SkipDeps   bool // 跳过依赖安装
    SkipNpm    bool // 跳过 npm 安装
    SkipGo     bool // 跳过 Go 工具安装
    TargetDir  string // 目标目录
}
```

### 代码生成模型

```go
// ServiceGenOptions 服务生成选项
type ServiceGenOptions struct {
    Name      string // 服务名称（小写）
    WithAPI   bool   // 是否生成 API 层
    WithMCP   bool   // 是否注册 MCP 工具
    WithTest  bool   // 是否生成测试文件
}
```

## 技术实现

### CLI 框架

使用 Go 的 [cobra](https://github.com/spf13/cobra) 库构建 CLI：
- 命令管理和路由
- 参数解析和验证
- 帮助信息生成

### 模板渲染

使用 Go 的 `text/template` 包：
- 定义代码模板文件
- 支持变量替换和条件渲染
- 模板文件存放在 `cli/templates/` 目录

### 文件操作

- 使用 `embed` 包嵌入模板文件到二进制
- 使用 `os` 和 `io` 包进行文件复制和修改
- 使用正则表达式替换占位符

### 依赖安装

使用 `os/exec` 包执行外部命令：
```go
exec.Command("go", "install", "github.com/air-verse/air@latest").Run()
exec.Command("npm", "install").Run()
```

## 目录结构

```
ai_skeleton/
├── cli/                          # CLI 工具源码
│   ├── cmd/                      # 命令定义
│   │   ├── root.go              # 根命令
│   │   ├── init.go              # init 命令
│   │   ├── generate.go          # generate 命令
│   │   └── config.go            # config 命令
│   ├── templates/                # 代码模板
│   │   ├── service.go.tmpl      # Service 模板
│   │   ├── api.go.tmpl          # API 模板
│   │   └── test.go.tmpl         # 测试模板
│   ├── internal/                 # 内部实现
│   │   ├── renderer/            # 模板渲染器
│   │   ├── installer/           # 依赖安装器
│   │   └── validator/           # 验证器
│   └── main.go                  # CLI 入口
```

## 占位符定义

需要在模板中替换的占位符：

| 占位符 | 说明 | 示例值 |
|--------|------|--------|
| `{{.ProjectName}}` | 项目名称（小写） | `my_project` |
| `{{.ProjectNameTitle}}` | 项目名称（标题格式） | `My Project` |
| `{{.ProjectDesc}}` | 项目描述 | `一个示例项目` |
| `{{.ProjectVersion}}` | 项目版本 | `1.0.0` |
| `{{.ModulePath}}` | Go 模块路径 | `github.com/user/my_project` |
| `{{.ServiceName}}` | 服务名称（小写） | `user` |
| `{{.ServiceNameTitle}}` | 服务名称（标题格式） | `User` |

## 文件替换规则

### 需要替换内容的文件

1. **backend/config.yaml**
   - `project.name`: 替换为新项目名称
   - `project.version`: 替换为新项目版本
   - `project.description`: 替换为新项目描述

2. **backend/go.mod**
   - 第一行 `module` 路径：替换为新模块路径
   - 所有 `github.com/richer/ai_skeleton` 引用：替换为新模块路径

3. **backend/所有 .go 文件**
   - import 中的 `github.com/richer/ai_skeleton`: 替换为新模块路径

4. **frontend/package.json**
   - `name`: 替换为新项目名称（kebab-case）
   - `version`: 替换为新项目版本
   - `description`: 替换为新项目描述

5. **README.md**
   - 第一行标题：替换为新项目名称
   - 项目描述段落：替换为新项目描述

6. **CLAUDE.md**
   - 标题中的项目名称：替换
   - 项目概述中的描述：替换

### 不需要修改的文件

- `.gitignore`
- `Makefile`
- `Dockerfile`
- `frontend/public/` 下的静态资源
- `backend/tmp/` 和 `backend/bin/`（临时目录）

## 环境检查逻辑

```go
// CheckEnvironment 检查环境依赖
func CheckEnvironment() error {
    // 检查 Go
    if _, err := exec.LookPath("go"); err != nil {
        return fmt.Errorf("Go 未安装，请先安装 Go: https://golang.org/dl/")
    }

    // 检查 npm
    if _, err := exec.LookPath("npm"); err != nil {
        return fmt.Errorf("npm 未安装，请先安装 Node.js: https://nodejs.org/")
    }

    // 显示版本信息
    goVersion := exec.Command("go", "version").Output()
    npmVersion := exec.Command("npm", "-v").Output()

    fmt.Printf("✓ Go: %s\n", goVersion)
    fmt.Printf("✓ npm: %s\n", npmVersion)

    return nil
}
```

## 依赖安装步骤

### Go 工具安装

```bash
# Air（热重载）
go install github.com/air-verse/air@latest

# Swagger CLI
go install github.com/swaggo/swag/cmd/swag@latest

# Gen-GORM（数据库代码生成）
go install gorm.io/gen/tools/gentool@latest
```

### Go 依赖安装

```bash
cd backend
go mod tidy
go mod download
```

### 前端依赖安装

```bash
cd frontend
npm install
```

### 验证安装

```bash
# 验证 Air
air -v

# 验证 Swagger
swag -v

# 验证前端依赖
cd frontend && npm list --depth=0
```

## 错误处理

### 环境检查失败

```
❌ 错误：Go 未安装

请访问 https://golang.org/dl/ 安装 Go 后重试。
```

### 依赖安装失败

```
❌ 错误：Air 安装失败

请手动执行以下命令：
  go install github.com/air-verse/air@latest

并确保 $GOPATH/bin 已添加到 PATH 环境变量。
```

### 目录已存在

```
❌ 错误：目录 "my_project" 已存在

请选择其他项目名称或删除现有目录后重试。
```

## 成功标准

初始化完成后，用户应该能够：

1. ✅ 执行 `make backend-dev`，后端成功启动在 http://localhost:8080
2. ✅ 执行 `make frontend-dev`，前端成功启动在 http://localhost:5173
3. ✅ 访问 http://localhost:5173，页面正常显示
4. ✅ 访问 http://localhost:8080/api/v1/health，返回健康检查结果
5. ✅ 访问 http://localhost:8080/api/v1/mcp/tools，返回 MCP 工具列表
6. ✅ 所有 Go import 路径正确，无编译错误
7. ✅ 前端 package.json 中的项目信息已更新
8. ✅ config.yaml 中的项目信息已更新

## 后续扩展（可选）

- 支持选择是否包含 MCP 协议模块
- 支持选择数据库类型（MySQL、PostgreSQL、SQLite）
- 支持选择 UI 组件库（Ant Design、Material-UI）
- 交互式选择技术栈组件
- 生成数据库迁移脚本
- 集成 Docker Compose 配置生成

## 开发计划

### 第一阶段：基础 CLI 框架
1. 搭建 Cobra CLI 基础结构
2. 实现 `init` 命令框架
3. 实现环境检查功能

### 第二阶段：项目初始化
1. 实现项目信息收集（交互式输入）
2. 实现文件复制和占位符替换
3. 实现依赖安装功能

### 第三阶段：代码生成
1. 实现 Service 代码生成
2. 实现 API 代码生成
3. 实现测试代码生成

### 第四阶段：配置管理
1. 实现配置模板生成
2. 实现配置验证功能

### 第五阶段：测试和文档
1. 编写单元测试
2. 编写集成测试（完整初始化流程）
3. 更新 README 和 CLAUDE.md

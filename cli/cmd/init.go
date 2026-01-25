package cmd

import (
	"fmt"

	"github.com/richer/ai_skeleton/cli/internal/installer"
	"github.com/richer/ai_skeleton/cli/internal/renderer"
	"github.com/spf13/cobra"
)

var (
	projectName    string
	projectDesc    string
	projectVersion string
	modulePath     string
	skipDeps       bool
	skipNpm        bool
	skipGo         bool
)

var initCmd = &cobra.Command{
	Use:   "init [项目名称]",
	Short: "初始化新项目",
	Long: `初始化一个新的 AI Skeleton 项目。

此命令会：
1. 检查环境依赖（Go、npm）
2. 收集项目信息
3. 复制模板文件并替换占位符
4. 安装依赖（Air、Swagger、npm packages）
5. 生成初始代码`,
	Args: cobra.MaximumNArgs(1),
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&projectName, "name", "n", "", "项目名称")
	initCmd.Flags().StringVarP(&projectDesc, "desc", "d", "", "项目描述")
	initCmd.Flags().StringVarP(&projectVersion, "version", "v", "1.0.0", "项目版本")
	initCmd.Flags().StringVarP(&modulePath, "module", "m", "", "Go 模块路径")
	initCmd.Flags().BoolVar(&skipDeps, "skipdeps", false, "跳过依赖安装")
	initCmd.Flags().BoolVar(&skipNpm, "skipnpm", false, "跳过 npm 依赖安装")
	initCmd.Flags().BoolVar(&skipGo, "skipgo", false, "跳过 Go 工具安装")
}

func runInit(cmd *cobra.Command, args []string) error {
	// 环境检查
	fmt.Println("🔍 检查环境依赖...")
	if err := installer.CheckEnvironment(); err != nil {
		return err
	}
	fmt.Println()

	// 收集项目信息
	if len(args) > 0 && projectName == "" {
		projectName = args[0]
	}

	meta, err := collectProjectInfo()
	if err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("📝 项目信息：")
	fmt.Printf("  名称：%s\n", meta.Name)
	fmt.Printf("  描述：%s\n", meta.Description)
	fmt.Printf("  版本：%s\n", meta.Version)
	fmt.Printf("  模块：%s\n", meta.Module)
	fmt.Println()

	// 生成项目文件
	fmt.Println("📦 生成项目文件...")
	if err := renderer.RenderProject(meta); err != nil {
		return fmt.Errorf("生成项目失败: %w", err)
	}
	fmt.Println("  ✓ 项目文件生成完成")
	fmt.Println()

	// 安装依赖
	if !skipDeps {
		fmt.Println("📥 安装依赖...")

		if !skipGo {
			if err := installer.InstallGoTools(meta.Name); err != nil {
				fmt.Printf("  ⚠️  Go 工具安装失败: %v\n", err)
			} else {
				fmt.Println("  ✓ Go 工具安装完成")
			}
		}

		if !skipNpm {
			if err := installer.InstallNpmDeps(meta.Name); err != nil {
				fmt.Printf("  ⚠️  npm 依赖安装失败: %v\n", err)
			} else {
				fmt.Println("  ✓ npm 依赖安装完成")
			}
		}
		fmt.Println()
	}

	// 完成提示
	fmt.Println("✅ 项目初始化完成！")
	fmt.Println()
	fmt.Println("下一步操作：")
	fmt.Printf("  1. cd %s\n", meta.Name)
	fmt.Println("  2. 启动后端：make backend-dev")
	fmt.Println("  3. 启动前端：make frontend-dev")
	fmt.Println("  4. 访问：http://localhost:5173")
	fmt.Println()

	return nil
}

func collectProjectInfo() (*renderer.ProjectMeta, error) {
	meta := &renderer.ProjectMeta{
		Name:        projectName,
		Description: projectDesc,
		Version:     projectVersion,
		Module:      modulePath,
	}

	// 使用交互式输入收集信息
	if err := renderer.PromptProjectInfo(meta); err != nil {
		return nil, err
	}

	return meta, nil
}

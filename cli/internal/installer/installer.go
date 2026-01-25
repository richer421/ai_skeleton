package installer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// CheckEnvironment 检查环境依赖
func CheckEnvironment() error {
	// 检查 Go
	goPath, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("❌ Go 未安装\n\n请访问 https://golang.org/dl/ 安装 Go 后重试")
	}

	// 检查 npm
	npmPath, err := exec.LookPath("npm")
	if err != nil {
		return fmt.Errorf("❌ npm 未安装\n\n请访问 https://nodejs.org/ 安装 Node.js 后重试")
	}

	// 获取版本信息
	goVersion := getCommandOutput("go", "version")
	npmVersion := getCommandOutput("npm", "-v")

	fmt.Printf("  ✓ Go: %s (路径: %s)\n", strings.TrimSpace(goVersion), goPath)
	fmt.Printf("  ✓ npm: v%s (路径: %s)\n", strings.TrimSpace(npmVersion), npmPath)

	return nil
}

// InstallGoTools 安装 Go 工具
func InstallGoTools(projectDir string) error {
	tools := []struct {
		name    string
		pkg     string
		checkCmd string
	}{
		{"Air", "github.com/air-verse/air@latest", "air"},
		{"Swagger", "github.com/swaggo/swag/cmd/swag@latest", "swag"},
	}

	for _, tool := range tools {
		fmt.Printf("  📦 安装 %s...\n", tool.name)
		cmd := exec.Command("go", "install", tool.pkg)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("安装 %s 失败: %w", tool.name, err)
		}
	}

	// 执行 go mod tidy
	fmt.Println("  📦 整理 Go 依赖...")
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectDir + "/backend"
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("执行 go mod tidy 失败: %w", err)
	}

	return nil
}

// InstallNpmDeps 安装 npm 依赖
func InstallNpmDeps(projectDir string) error {
	fmt.Println("  📦 安装前端依赖...")
	cmd := exec.Command("npm", "install")
	cmd.Dir = projectDir + "/frontend"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("npm install 失败: %w", err)
	}

	return nil
}

// getCommandOutput 获取命令输出
func getCommandOutput(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(output)
}

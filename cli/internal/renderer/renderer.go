package renderer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
)

// ProjectMeta 项目元信息
type ProjectMeta struct {
	Name        string // 项目名称
	Description string // 项目描述
	Version     string // 项目版本
	Module      string // Go 模块路径
}

// PromptProjectInfo 交互式收集项目信息
func PromptProjectInfo(meta *ProjectMeta) error {
	// 项目名称
	if meta.Name == "" {
		prompt := promptui.Prompt{
			Label:   "项目名称",
			Default: filepath.Base(getCurrentDir()),
		}
		name, err := prompt.Run()
		if err != nil {
			return err
		}
		meta.Name = strings.TrimSpace(name)
	}

	// 项目描述
	if meta.Description == "" {
		prompt := promptui.Prompt{
			Label:   "项目描述",
			Default: "",
		}
		desc, err := prompt.Run()
		if err != nil {
			return err
		}
		meta.Description = strings.TrimSpace(desc)
	}

	// 项目版本
	if meta.Version == "" {
		meta.Version = "1.0.0"
	}

	// Go 模块路径
	if meta.Module == "" {
		defaultModule := fmt.Sprintf("github.com/user/%s", meta.Name)
		prompt := promptui.Prompt{
			Label:   "Go 模块路径",
			Default: defaultModule,
		}
		module, err := prompt.Run()
		if err != nil {
			return err
		}
		meta.Module = strings.TrimSpace(module)
	}

	return nil
}

// RenderProject 渲染项目文件
func RenderProject(meta *ProjectMeta) error {
	// 检查目标目录是否存在
	if _, err := os.Stat(meta.Name); err == nil {
		return fmt.Errorf("目录 %s 已存在，请选择其他项目名称", meta.Name)
	}

	// 获取当前脚手架根目录
	scaffoldRoot, err := getScaffoldRoot()
	if err != nil {
		return err
	}

	// 复制文件
	if err := copyDir(scaffoldRoot, meta.Name, meta); err != nil {
		return err
	}

	return nil
}

// getCurrentDir 获取当前目录
func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "my_project"
	}
	return filepath.Base(dir)
}

// getScaffoldRoot 获取脚手架根目录
func getScaffoldRoot() (string, error) {
	// 方案1：从当前工作目录向上查找（开发模式）
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// 检查当前目录是否是脚手架根目录
	if isScaffoldRoot(cwd) {
		return cwd, nil
	}

	// 检查父目录
	parent := filepath.Dir(cwd)
	if isScaffoldRoot(parent) {
		return parent, nil
	}

	return "", fmt.Errorf("未找到脚手架根目录，请确保在 ai_skeleton 目录或其子目录中运行")
}

// isScaffoldRoot 判断是否是脚手架根目录
func isScaffoldRoot(dir string) bool {
	// 检查关键文件是否存在
	markers := []string{
		"backend/go.mod",
		"frontend/package.json",
		"Makefile",
		"CLAUDE.md",
	}

	for _, marker := range markers {
		if _, err := os.Stat(filepath.Join(dir, marker)); err != nil {
			return false
		}
	}

	return true
}

// copyDir 复制目录并替换占位符
func copyDir(src, dst string, meta *ProjectMeta) error {
	// 创建目标目录
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	// 遍历源目录
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过 CLI 目录、临时目录、构建产物
		relPath, _ := filepath.Rel(src, path)
		if shouldSkip(relPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// 计算目标路径
		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		// 复制文件并替换内容
		return copyFileWithReplace(path, dstPath, meta)
	})
}

// shouldSkip 判断是否跳过文件/目录
func shouldSkip(path string) bool {
	skipDirs := []string{
		"cli",
		".git",
		"backend/tmp",
		"backend/bin",
		"frontend/node_modules",
		"frontend/dist",
		"requirements",
	}

	for _, skip := range skipDirs {
		if strings.HasPrefix(path, skip) {
			return true
		}
	}

	return false
}

// copyFileWithReplace 复制文件并替换占位符
func copyFileWithReplace(src, dst string, meta *ProjectMeta) error {
	// 读取源文件
	content, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// 替换占位符
	newContent := replaceContent(string(content), meta)

	// 写入目标文件
	return os.WriteFile(dst, []byte(newContent), 0644)
}

// replaceContent 替换文件内容中的占位符
func replaceContent(content string, meta *ProjectMeta) string {
	result := content

	// 1. 替换 Go 模块路径（必须先替换，避免被其他规则影响）
	result = strings.ReplaceAll(result, "github.com/richer/ai_skeleton", meta.Module)

	// 2. 替换项目名称
	result = strings.ReplaceAll(result, "ai_skeleton", meta.Name)
	result = strings.ReplaceAll(result, "AI Skeleton", toTitle(meta.Name))
	result = strings.ReplaceAll(result, "ai-skeleton", toKebabCase(meta.Name))

	// 3. 特殊处理 config.yaml 中的 project 部分
	if strings.Contains(result, "project:") && strings.Contains(result, "name:") {
		result = strings.ReplaceAll(result, "name: \""+meta.Name+"\"", fmt.Sprintf("name: \"%s\"", meta.Name))
		result = strings.ReplaceAll(result, "version: \"1.0.0\"", fmt.Sprintf("version: \"%s\"", meta.Version))
		if meta.Description != "" {
			// 查找并替换 description 行
			lines := strings.Split(result, "\n")
			for i, line := range lines {
				if strings.Contains(line, "description:") {
					lines[i] = fmt.Sprintf("  description: \"%s\"", meta.Description)
					break
				}
			}
			result = strings.Join(lines, "\n")
		}
	}

	// 4. 特殊处理 package.json
	if strings.Contains(result, "\"name\":") && strings.Contains(result, "\"version\":") {
		lines := strings.Split(result, "\n")
		for i, line := range lines {
			if strings.Contains(line, "\"name\":") {
				lines[i] = fmt.Sprintf("  \"name\": \"%s-frontend\",", toKebabCase(meta.Name))
			} else if strings.Contains(line, "\"version\":") {
				lines[i] = fmt.Sprintf("  \"version\": \"%s\",", meta.Version)
			} else if meta.Description != "" && strings.Contains(line, "\"description\":") {
				lines[i] = fmt.Sprintf("  \"description\": \"%s\",", meta.Description)
			}
		}
		result = strings.Join(lines, "\n")
	}

	return result
}

// toTitle 转换为标题格式
func toTitle(s string) string {
	words := strings.Split(s, "_")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, " ")
}

// toKebabCase 转换为 kebab-case
func toKebabCase(s string) string {
	return strings.ReplaceAll(strings.ToLower(s), "_", "-")
}

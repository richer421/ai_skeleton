package renderer

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
)

// DefaultTemplateURL 默认模板仓库地址
const DefaultTemplateURL = "https://github.com/richer421/ai_skeleton/archive/main.zip"

// ProjectMeta 项目元信息
type ProjectMeta struct {
	Name         string // 项目名称
	Description  string // 项目描述
	Version      string // 项目版本
	Module       string // Go 模块路径
	TemplateURL  string // 自定义模板仓库地址（可选，用于私有仓库）
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

	// 确定模板URL
	templateURL := DefaultTemplateURL
	if meta.TemplateURL != "" {
		templateURL = meta.TemplateURL
		fmt.Println("  🌐 正在从私有仓库获取模板...")
	} else {
		fmt.Println("  🌐 正在从官方仓库获取最新模板...")
	}

	// 从远程下载模板
	if err := downloadAndExtractTemplate(meta.Name, templateURL, meta); err != nil {
		return fmt.Errorf("下载模板失败: %w", err)
	}

	fmt.Println("  ✓ 项目文件生成完成")
	return nil
}

// downloadAndExtractTemplate 从远程下载并提取模板
func downloadAndExtractTemplate(dst, templateURL string, meta *ProjectMeta) error {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "ai_skeleton_template_*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// 下载模板
	zipPath := filepath.Join(tempDir, "template.zip")
	if err := downloadFile(templateURL, zipPath); err != nil {
		return err
	}

	// 解压模板
	extractPath := filepath.Join(tempDir, "extracted")
	if err := unzip(zipPath, extractPath); err != nil {
		return err
	}

	// 查找实际的模板目录（通常是 ai_skeleton-main）
	templateDir := ""
	entries, err := os.ReadDir(extractPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			templateDir = filepath.Join(extractPath, entry.Name())
			break
		}
	}

	if templateDir == "" {
		return fmt.Errorf("无法找到模板目录")
	}

	// 复制并处理模板
	return copyDir(templateDir, dst, meta)
}

// downloadFile 下载文件
func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败: 状态码 %d", resp.StatusCode)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// unzip 解压ZIP文件
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("非法文件路径: %s", fpath)
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
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

		// 跳过 CLI 目录、临时目录、构建产物和其他无关目录
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
		".github",
		".vscode",
		".idea",
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

// getCurrentDir 获取当前目录
func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return "my_project"
	}
	return filepath.Base(dir)
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

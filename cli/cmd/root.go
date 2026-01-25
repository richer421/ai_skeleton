package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const version = "1.0.0"

var rootCmd = &cobra.Command{
	Use:   "ai-skeleton",
	Short: "AI Skeleton - AI 全栈脚手架 CLI 工具",
	Long: `AI Skeleton 是一个面向 AI 辅助开发的全栈项目脚手架。

此 CLI 工具帮助你快速创建新项目、生成代码模板、管理配置文件。`,
	Version: version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

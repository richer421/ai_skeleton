package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:     "generate",
	Aliases: []string{"gen"},
	Short:   "代码生成工具",
	Long:    `生成 Service、API、测试等代码模板。`,
}

var (
	withAPI  bool
	withMCP  bool
	withTest bool
)

var generateServiceCmd = &cobra.Command{
	Use:   "service [服务名]",
	Short: "生成服务代码",
	Long: `生成完整的服务代码，包括 Service 层、API 层、测试文件。

示例：
  ai-skeleton generate service user
  ai-skeleton gen service order --with-mcp`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serviceName := args[0]
		fmt.Printf("生成服务代码: %s\n", serviceName)
		fmt.Printf("  - 生成 API: %v\n", withAPI)
		fmt.Printf("  - 注册 MCP: %v\n", withMCP)
		fmt.Printf("  - 生成测试: %v\n", withTest)
		// TODO: 实现代码生成逻辑
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.AddCommand(generateServiceCmd)

	generateServiceCmd.Flags().BoolVar(&withAPI, "with-api", true, "生成 API 层")
	generateServiceCmd.Flags().BoolVar(&withMCP, "with-mcp", false, "注册 MCP 工具")
	generateServiceCmd.Flags().BoolVar(&withTest, "with-test", true, "生成测试文件")
}

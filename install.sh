#!/bin/bash

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检测操作系统和架构
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# 转换架构名称
case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}不支持的架构: $ARCH${NC}"
        exit 1
        ;;
esac

# 转换操作系统名称
case $OS in
    darwin)
        OS="darwin"
        ;;
    linux)
        OS="linux"
        ;;
    *)
        echo -e "${RED}不支持的操作系统: $OS${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}AIS 安装脚本${NC}"
echo "操作系统: $OS"
echo "架构: $ARCH"
echo ""

# 检查依赖
echo -e "${YELLOW}检查依赖...${NC}"

if ! command -v go &> /dev/null; then
    echo -e "${RED}错误: Go 未安装${NC}"
    echo "请访问 https://golang.org/dl/ 安装 Go"
    exit 1
fi

if ! command -v npm &> /dev/null; then
    echo -e "${RED}错误: npm 未安装${NC}"
    echo "请访问 https://nodejs.org/ 安装 Node.js"
    exit 1
fi

echo -e "${GREEN}✓ Go: $(go version)${NC}"
echo -e "${GREEN}✓ npm: v$(npm -v)${NC}"
echo ""

# 克隆仓库
REPO_URL="https://github.com/richer421/ai_skeleton.git"
INSTALL_DIR="$HOME/.ai-skeleton"

echo -e "${YELLOW}下载 AI Skeleton...${NC}"

if [ -d "$INSTALL_DIR" ]; then
    echo "目录已存在，更新中..."
    cd "$INSTALL_DIR"
    git pull
else
    git clone "$REPO_URL" "$INSTALL_DIR"
    cd "$INSTALL_DIR"
fi

echo -e "${GREEN}✓ 下载完成${NC}"
echo ""

# 编译 CLI 工具
echo -e "${YELLOW}编译 CLI 工具...${NC}"
cd cli
go build -o ai-skeleton main.go
echo -e "${GREEN}✓ 编译完成${NC}"
echo ""

# 创建符号链接或复制到 PATH
BIN_DIR="/usr/local/bin"
CLI_PATH="$INSTALL_DIR/cli/ai-skeleton"

if [ -w "$BIN_DIR" ]; then
    ln -sf "$CLI_PATH" "$BIN_DIR/ai-skeleton"
    echo -e "${GREEN}✓ 已安装到 $BIN_DIR/ai-skeleton${NC}"
else
    echo -e "${YELLOW}提示: 无法写入 $BIN_DIR，需要 sudo 权限${NC}"
    if sudo ln -sf "$CLI_PATH" "$BIN_DIR/ai-skeleton"; then
        echo -e "${GREEN}✓ 已安装到 $BIN_DIR/ai-skeleton${NC}"
    else
        echo -e "${RED}安装失败，请手动添加到 PATH:${NC}"
        echo "export PATH=\"\$PATH:$INSTALL_DIR/cli\""
        exit 1
    fi
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}  AI Skeleton 安装成功！${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "使用方法："
echo "  1. 创建新项目:"
echo -e "     ${YELLOW}ai-skeleton init my_project${NC}"
echo ""
echo "  2. 进入项目目录:"
echo -e "     ${YELLOW}cd my_project${NC}"
echo ""
echo "  3. 启动后端:"
echo -e "     ${YELLOW}make backend-dev${NC}"
echo ""
echo "  4. 启动前端 (另开终端):"
echo -e "     ${YELLOW}make frontend-dev${NC}"
echo ""
echo "更多命令:"
echo -e "  ${YELLOW}ai-skeleton --help${NC}"
echo ""

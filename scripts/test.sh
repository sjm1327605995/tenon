#!/bin/bash
# 无头测试脚本 — 用 Xvfb 提供虚拟显示器
# 用法: ./scripts/test.sh [go test 参数]
# 示例: ./scripts/test.sh ./pkg/v2/ui/ -v -run TestBuilder
#        ./scripts/test.sh ./... -count=1

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

# 检查 Xvfb 是否可用
if ! command -v Xvfb &>/dev/null; then
    echo "ERROR: Xvfb not found. Install with: apt-get install xvfb" >&2
    exit 1
fi

# 查找可用的 display 编号
DISPLAY_NUM=99
for i in $(seq 99 199); do
    if ! ls /tmp/.X${i}-lock &>/dev/null 2>&1; then
        DISPLAY_NUM=$i
        break
    fi
done

export DISPLAY=":${DISPLAY_NUM}"

# 启动 Xvfb
Xvfb "$DISPLAY" -screen 0 1024x768x24 -ac +extension GLX +render -noreset &>/dev/null &
XVFB_PID=$!

# 确保 Xvfb 退出时清理
cleanup() {
    kill "$XVFB_PID" 2>/dev/null || true
    rm -f "/tmp/.X${DISPLAY_NUM}-lock" 2>/dev/null || true
}
trap cleanup EXIT

# 等待 Xvfb 就绪
sleep 0.5

# 运行测试
cd "$PROJECT_DIR"
go test "$@"

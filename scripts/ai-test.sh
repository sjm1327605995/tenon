#!/bin/bash
# AI 友好的测试脚本 — 输出结构化结果，并始终保留 go test 的退出状态
# 用法: ./scripts/ai-test.sh [包路径]
# 示例: ./scripts/ai-test.sh
#       ./scripts/ai-test.sh ./pkg/ui/
#       ./scripts/ai-test.sh ./pkg/shadcn/

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_DIR"

PKG="${1:-./...}"
PASS=0
FAIL=0
SKIP=0

run_with_xvfb() {
    DISPLAY_NUM=99
    for i in $(seq 99 199); do
        if ! ls /tmp/.X${i}-lock &>/dev/null 2>&1; then
            DISPLAY_NUM=$i
            break
        fi
    done
    export DISPLAY=":${DISPLAY_NUM}"
    Xvfb "$DISPLAY" -screen 0 1024x768x24 -ac +extension GLX +render -noreset &>/dev/null &
    XVFB_PID=$!
    cleanup() {
        kill "$XVFB_PID" 2>/dev/null || true
        rm -f "/tmp/.X${DISPLAY_NUM}-lock" 2>/dev/null || true
    }
    trap cleanup EXIT
    sleep 0.3
}

# Linux 无显示环境且装有 Xvfb 时主动提供虚拟显示器；否则直接运行测试，
# 让真实错误和退出码暴露出来，不把编译或断言失败误报为 SKIP。
if [[ "$(uname -s)" == Linux* && -z "${DISPLAY:-}" ]] && command -v Xvfb &>/dev/null; then
    run_with_xvfb
    echo "Xvfb started on DISPLAY=$DISPLAY"
fi

echo "=== Running tests: $PKG ==="
echo ""

set +e
TEST_OUTPUT=$(go test "$PKG" -count=1 -timeout 60s -v 2>&1)
TEST_STATUS=$?
set -e

while IFS= read -r line; do
    if [[ "$line" =~ ^[[:space:]]*---[[:space:]]PASS: ]]; then
        name="${line#*PASS: }"
        name="${name%% (*}"
        PASS=$((PASS + 1))
        echo "  ✓ $name"
    elif [[ "$line" =~ ^[[:space:]]*---[[:space:]]FAIL: ]]; then
        name="${line#*FAIL: }"
        name="${name%% (*}"
        FAIL=$((FAIL + 1))
        echo "  ✗ $name"
    elif [[ "$line" =~ ^[[:space:]]*---[[:space:]]SKIP: ]]; then
        name="${line#*SKIP: }"
        name="${name%% (*}"
        SKIP=$((SKIP + 1))
        echo "  ○ $name"
    elif [[ "$line" == ok* ]]; then
        pkg=$(echo "$line" | awk '{print $2}')
        echo "  Package: $pkg — OK"
    fi
done <<< "$TEST_OUTPUT"

echo ""
echo "=== Summary ==="
echo "  PASS: $PASS"
echo "  FAIL: $FAIL"
echo "  SKIP: $SKIP"

if [ "$TEST_STATUS" -ne 0 ]; then
    echo ""
    echo "=== Failed test output ==="
    echo "$TEST_OUTPUT"
    exit "$TEST_STATUS"
fi

echo ""
echo "All tests passed."

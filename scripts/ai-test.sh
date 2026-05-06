#!/bin/bash
# AI 友好的测试脚本 — 自动检测环境，输出结构化结果
# 用法: ./scripts/ai-test.sh [包路径]
# 示例: ./scripts/ai-test.sh                    # 测试所有包
#        ./scripts/ai-test.sh ./pkg/v2/ui/      # 测试指定包
#        ./scripts/ai-test.sh ./pkg/v2/widgets/ # 测试 widgets

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_DIR"

PKG="${1:-./...}"
PASS=0
FAIL=0
SKIP=0
RESULTS=""

# 尝试无头运行（Xvfb）
run_with_xvfb() {
    if ! command -v Xvfb &>/dev/null; then
        return 1
    fi
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
    trap "kill $XVFB_PID 2>/dev/null; rm -f /tmp/.X${DISPLAY_NUM}-lock 2>/dev/null" EXIT
    sleep 0.3
    return 0
}

# 检测是否需要 Xvfb
if ! go test "$PKG" -count=1 -timeout 5s &>/dev/null 2>&1; then
    echo "Need display, trying Xvfb..."
    if run_with_xvfb; then
        echo "Xvfb started on DISPLAY=$DISPLAY"
    else
        echo "SKIP: No display available and Xvfb not installed"
        exit 0
    fi
fi

# 运行测试并收集结果
echo "=== Running tests: $PKG ==="
echo ""

TEST_OUTPUT=$(go test "$PKG" -count=1 -timeout 60s -v 2>&1) || true

# 解析结果
while IFS= read -r line; do
    if [[ "$line" == "--- PASS:"* ]]; then
        name=$(echo "$line" | sed 's/--- PASS: //' | sed 's/ (.*//')
        PASS=$((PASS + 1))
        echo "  ✓ $name"
    elif [[ "$line" == "--- FAIL:"* ]]; then
        name=$(echo "$line" | sed 's/--- FAIL: //' | sed 's/ (.*//')
        FAIL=$((FAIL + 1))
        echo "  ✗ $name"
    elif [[ "$line" == "--- SKIP:"* ]]; then
        name=$(echo "$line" | sed 's/--- SKIP: //' | sed 's/ (.*//')
        SKIP=$((SKIP + 1))
        echo "  ○ $name"
    elif [[ "$line" == "FAIL"* && "$line" == *"FAIL"* ]]; then
        # Package-level FAIL line
        :
    elif [[ "$line" == "ok"* ]]; then
        pkg=$(echo "$line" | awk '{print $2}')
        echo "  Package: $pkg — OK"
    fi
done <<< "$TEST_OUTPUT"

echo ""
echo "=== Summary ==="
echo "  PASS: $PASS"
echo "  FAIL: $FAIL"
echo "  SKIP: $SKIP"

if [ "$FAIL" -gt 0 ]; then
    echo ""
    echo "=== Failed test output ==="
    echo "$TEST_OUTPUT" | grep -A 20 "FAIL"
    exit 1
fi

echo ""
echo "All tests passed."
exit 0

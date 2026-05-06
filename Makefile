.PHONY: test test-ui test-widgets test-yoga test-all bench lint vet

# 快速测试（跳过需要 GUI 的测试）
test-yoga:
	go test ./yoga/... -v -count=1

# 完整测试（需要 Xvfb 或真实显示器）
test-ui:
	./scripts/test.sh ./pkg/v2/ui/ -v -count=1

test-widgets:
	./scripts/test.sh ./pkg/v2/widgets/ -v -count=1

test-render:
	./scripts/test.sh ./pkg/v2/render/ -v -count=1

# 全量测试
test-all:
	./scripts/test.sh ./... -count=1

# AI 友好测试（自动检测环境）
test:
	./scripts/ai-test.sh ./...

# 代码检查
vet:
	go vet ./...

lint: vet
	@echo "All checks passed."

# 构建验证
build:
	go build ./...

# benchmark
bench:
	./scripts/test.sh ./yoga/... -bench=. -benchmem

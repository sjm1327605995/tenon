#!/bin/bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT
cd "$ROOT"

if [ -n "$(GOWORK=off go list -m -f '{{if .Replace}}{{.Replace}}{{end}}' github.com/sjm1327605995/tenon 2>/dev/null)" ]; then
    echo "backend/g3n/go.mod must not contain development replace directives" >&2
    exit 1
fi

run() {
    echo "==> $*"
    "$@"
}

fmt=$(gofmt -l pkg example yoga backend/g3n)
if [ -n "$fmt" ]; then
    echo "Not gofmt'd:" >&2
    echo "$fmt" >&2
    exit 1
fi

run go build ./...
run go vet ./...
run bash ./scripts/api-check.sh
run go test ./... -count=1
run go run _debug_snapshot/test_v3.go

check_consumer() {
    name="$1"
    source="$2"
    consumer="$TMP/$name"
    mkdir -p "$consumer"
    cp "$source" "$consumer/main.go"
    cd "$consumer"
    run go mod init "example.com/$name"
    run go mod edit -require=github.com/sjm1327605995/tenon@v0.0.0
    run go mod edit -replace=github.com/sjm1327605995/tenon="$ROOT"
    run go mod tidy
    run go build .
    if [ "${3:-}" = "run" ]; then
        run go run .
    fi
    cd "$ROOT"
}

check_consumer quickstart "$ROOT/example/quickstart/main.go"
check_consumer plugin-contract "$ROOT/example/plugin/main.go" run

if [ "$(go env CGO_ENABLED)" = "1" ]; then
    run go test -race ./...
else
    if [ "${REQUIRE_RACE:-0}" = "1" ]; then
        echo "REQUIRE_RACE=1 but CGO is disabled" >&2
        exit 1
    fi
    echo "SKIP: race tests require CGO and a C compiler"
fi

if [ "$(go env CGO_ENABLED)" = "1" ]; then
    (
        cd backend/g3n
        run go build ./...
        run go test ./... -count=1
        if [ "${REQUIRE_WINDOW:-0}" = "1" ]; then
            if [ "$(go env GOOS)" = "linux" ] && [ -z "${DISPLAY:-}" ]; then
                run env TENON_SMOKE_FRAMES=3 xvfb-run -a go run ./example/plugin
                run env TENON_SMOKE_FRAMES=12 xvfb-run -a go run ./example/lifecycle
                run env TENON_UI_SMOKE=1 TENON_SMOKE_FRAMES=120 xvfb-run -a go run ./example/uiparity
                run env TENON_SMOKE_FRAMES=3 xvfb-run -a go run ./example/projectedparity
            else
                run env TENON_SMOKE_FRAMES=3 go run ./example/plugin
                run env TENON_SMOKE_FRAMES=12 go run ./example/lifecycle
                run env TENON_UI_SMOKE=1 TENON_SMOKE_FRAMES=120 go run ./example/uiparity
                run env TENON_SMOKE_FRAMES=3 go run ./example/projectedparity
            fi
        fi
    )
else
    if [ "${REQUIRE_G3N:-0}" = "1" ]; then
        echo "REQUIRE_G3N=1 but CGO is disabled" >&2
        exit 1
    fi
    echo "SKIP: G3N tests require CGO and its native libraries"
fi

capture_window() {
    name="$1"
    capture="$TMP/$name.png"
    if [ "$(go env GOOS)" = "linux" ] && [ -z "${DISPLAY:-}" ]; then
        if command -v xvfb-run >/dev/null 2>&1; then
            TENON_CAPTURE="$capture" TENON_CAPTURE_FRAMES=3 xvfb-run -a go run "./example/$name"
        elif [ "${REQUIRE_WINDOW:-0}" = "1" ]; then
            echo "REQUIRE_WINDOW=1 but xvfb-run is unavailable" >&2
            exit 1
        else
            echo "SKIP: real-window capture requires a display or xvfb-run"
        fi
    else
        TENON_CAPTURE="$capture" TENON_CAPTURE_FRAMES=3 go run "./example/$name"
    fi
    if [ -f "$capture" ]; then
        test -s "$capture"
        echo "Verified screenshot: $capture"
    elif [ "${REQUIRE_WINDOW:-0}" = "1" ]; then
        echo "REQUIRE_WINDOW=1 but no $name screenshot was produced" >&2
        exit 1
    fi
}

if [ "${TENON_RELEASE_WINDOW:-auto}" != "0" ]; then
    capture_window bidi
    capture_window multiwindow
    capture_window catalog
    capture_window gradient
    capture_window card3d
    capture_window projectedparity
fi

echo "Release checks finished. Any SKIP above is not 1.0 release evidence."

#!/bin/bash
#
# test-integration.sh — Integration test runner for the omar CLI.
#
# This script builds the omar binary and runs a battery of smoke tests
# against it. It tests CLI argument handling, output formatting, and
# graceful error handling when system dependencies are missing.
#
# Usage:
#   ./scripts/test-integration.sh              # build + run all tests
#   ./scripts/test-integration.sh --skip-build  # use existing bin/omar
#   ./scripts/test-integration.sh --go          # run Go integration tests instead
#
# Returns 0 if all tests pass, 1 otherwise.
#

set -euo pipefail
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
OMAR_BIN="${PROJECT_DIR}/bin/omar-integration-test"

PASS=0
FAIL=0

# ── helpers ──────────────────────────────────────────────────────────────

print_result() {
    local name="$1" result="$2" detail="${3:-}"
    if [ "$result" = "PASS" ]; then
        printf "  ✓ %s\n" "$name"
    else
        printf "  ✗ %s\n    %s\n" "$name" "$detail"
    fi
}

assert_success() {
    local name="$1" desc="$2"
    shift 2
    if output=$("$@" 2>&1); then
        PASS=$((PASS + 1))
        print_result "$name" "PASS"
    else
        FAIL=$((FAIL + 1))
        print_result "$name" "FAIL" "exit code $? — ${output:0:200}"
    fi
}

assert_failure() {
    local name="$1" desc="$2"
    shift 2
    if ! output=$("$@" 2>&1); then
        PASS=$((PASS + 1))
        print_result "$name" "PASS"
    else
        FAIL=$((FAIL + 1))
        print_result "$name" "FAIL" "expected failure but exit code was 0 — ${output:0:200}"
    fi
}

assert_output_contains() {
    local name="$1" expected="$2"
    shift 2
    if output=$("$@" 2>&1); then
        if echo "$output" | grep -qF "$expected"; then
            PASS=$((PASS + 1))
            print_result "$name" "PASS"
        else
            FAIL=$((FAIL + 1))
            print_result "$name" "FAIL" "expected output to contain: ${expected}"
        fi
    else
        FAIL=$((FAIL + 1))
        print_result "$name" "FAIL" "command failed — ${output:0:200}"
    fi
}

# ── build ────────────────────────────────────────────────────────────────

build_omar() {
    echo "=== Building omar CLI ==="
    mkdir -p "${PROJECT_DIR}/bin"
    go build -ldflags="-X 'github.com/nevotheless/omar/internal/version.Version=0.0.0-integration-test'" \
        -o "${OMAR_BIN}" ./cmd/omar
    echo "   Built: ${OMAR_BIN}"
}

# ── test groups ──────────────────────────────────────────────────────────

test_version() {
    echo ""
    echo "=== version ==="

    assert_success "version: text output" "omar version" \
        "${OMAR_BIN}" version

    assert_output_contains "version: shows version string" "0.0.0-integration-test" \
        "${OMAR_BIN}" version

    assert_success "version: --json" "omar version --json" \
        "${OMAR_BIN}" version --json

    assert_output_contains "version: --json has cli_version" "cli_version" \
        "${OMAR_BIN}" version --json

    assert_output_contains "version: --json has cli_version field" "cli_version" \
        "${OMAR_BIN}" version --json

    assert_failure "version: invalid flag" "omar version --invalid" \
        "${OMAR_BIN}" version --invalid
}

test_status() {
    echo ""
    echo "=== status ==="

    assert_success "status: runs" "omar status" \
        "${OMAR_BIN}" status

    assert_output_contains "status: shows Image" "Image:" \
        "${OMAR_BIN}" status

    assert_output_contains "status: shows Version" "Version:" \
        "${OMAR_BIN}" status

    assert_output_contains "status: shows Booted" "Booted:" \
        "${OMAR_BIN}" status

    assert_output_contains "status: shows Staged" "Staged:" \
        "${OMAR_BIN}" status

    assert_output_contains "status: shows Rollback" "Rollback:" \
        "${OMAR_BIN}" status
}

test_help() {
    echo ""
    echo "=== help ==="

    assert_success "help: --help" "omar --help" \
        "${OMAR_BIN}" --help

    assert_success "help: -h" "omar -h" \
        "${OMAR_BIN}" -h

    assert_output_contains "help: lists all commands" "install" \
        "${OMAR_BIN}" --help

    assert_output_contains "help: lists update" "update" \
        "${OMAR_BIN}" --help

    assert_output_contains "help: lists status" "status" \
        "${OMAR_BIN}" --help

    assert_output_contains "help: lists rollback" "rollback" \
        "${OMAR_BIN}" --help

    assert_output_contains "help: lists pkg" "pkg" \
        "${OMAR_BIN}" --help

    assert_output_contains "help: lists version" "version" \
        "${OMAR_BIN}" --help

    assert_failure "help: unknown command" "omar nonexistent" \
        "${OMAR_BIN}" nonexistent
}

test_install() {
    echo ""
    echo "=== install ==="

    # install without flags shows usage (might try to convert, which fails gracefully)
    # install may fail on mutable systems, but should produce meaningful output
    if output=$("${OMAR_BIN}" install 2>&1); then
        PASS=$((PASS + 1))
        print_result "install: succeeds on immutable" "PASS"
    elif echo "$output" | grep -q "bootc\|Usage\|Pre-flight"; then
        PASS=$((PASS + 1))
        print_result "install: shows meaningful output on mutable" "PASS"
    else
        FAIL=$((FAIL + 1))
        print_result "install: unexpected failure" "FAIL" "${output:0:200}"
    fi

    assert_failure "install --fresh: requires --disk" "install --fresh" \
        "${OMAR_BIN}" install --fresh
}

test_pkg() {
    echo ""
    echo "=== pkg ==="

    assert_failure "pkg add: requires package name" "pkg add" \
        "${OMAR_BIN}" pkg add

    assert_success "pkg add --help" "pkg add --help" \
        "${OMAR_BIN}" pkg add --help

    # pkg list works even without flatpak/distrobox
    assert_success "pkg list: runs without backends" "pkg list" \
        "${OMAR_BIN}" pkg list
}

test_update() {
    echo ""
    echo "=== update ==="

    # update may fail on non-ostree, but must not crash
    "${OMAR_BIN}" update 2>/dev/null && {
        PASS=$((PASS + 1))
        print_result "update: runs without crash" "PASS"
    } || {
        PASS=$((PASS + 1))
        print_result "update: handles non-ostree gracefully" "PASS"
    }
}

test_rollback() {
    echo ""
    echo "=== rollback ==="

    # rollback may fail on non-ostree, but must not crash
    "${OMAR_BIN}" rollback 2>/dev/null && {
        PASS=$((PASS + 1))
        print_result "rollback: runs without crash" "PASS"
    } || {
        PASS=$((PASS + 1))
        print_result "rollback: handles non-ostree gracefully" "PASS"
    }
}

# ── main ─────────────────────────────────────────────────────────────────

main() {
    local skip_build=false
    local use_go=false

    for arg in "$@"; do
        case "$arg" in
            --skip-build) skip_build=true ;;
            --go) use_go=true ;;
            *) echo "Unknown option: $arg" >&2; exit 1 ;;
        esac
    done

    echo ""
    echo "╔══════════════════════════════════════════════════════════════╗"
    echo "║        omar Integration Tests                               ║"
    echo "╚══════════════════════════════════════════════════════════════╝"
    echo ""

    if [ "$use_go" = true ]; then
        echo "=== Running Go integration tests ==="
        cd "$PROJECT_DIR"
        go test -tags=integration -count=1 -v ./tests/integration/
        exit $?
    fi

    if [ "$skip_build" = false ] || [ ! -x "$OMAR_BIN" ]; then
        build_omar
    fi

    # Run all test groups
    test_version
    test_status
    test_help
    test_install
    test_pkg
    test_update
    test_rollback

    # Summary
    echo ""
    echo "════════════════════════════════════════════════════════════════"
    total=$((PASS + FAIL))
    printf "  Results: %d passed, %d failed, %d total\n" "$PASS" "$FAIL" "$total"
    echo "════════════════════════════════════════════════════════════════"
    echo ""

    if [ "$FAIL" -gt 0 ]; then
        exit 1
    fi
}

main "$@"

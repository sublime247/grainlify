#!/usr/bin/env bash
set -euo pipefail

DEPLOY_SCRIPT="./contracts/scripts/deploy.sh"

fail() { echo "✘ FAIL: $1"; exit 1; }
pass() { echo "✔ PASS: $1"; }

run_expect_fail() {
    desc="$1"
    expected_msg="$2"
    shift 2

    set +e
    output=$($DEPLOY_SCRIPT "$@" 2>&1)
    exit_code=$?
    set -e

    if [[ $exit_code -eq 0 ]]; then
        echo "$output"
        fail "$desc (expected failure, got exit 0)"
    fi

    if ! echo "$output" | grep -q "$expected_msg"; then
        echo "$output"
        fail "$desc (expected message '$expected_msg')"
    fi

    pass "$desc"
}

echo "=== Deployment Script Failure Tests ==="

# Missing WASM file
run_expect_fail "Missing WASM file" "No WASM file specified"

# Invalid WASM path
run_expect_fail "Invalid WASM file path" "WASM file does not exist" "nonexistent.wasm"

# Missing config file
run_expect_fail "Missing config file" "Config file not found" "./test_data/bad.wasm" --config "./test_data/missing.env"

# Missing identity
run_expect_fail "Invalid identity" "Identity not found" "./test_data/bad.wasm" --identity "ghost_id"

# Missing dependency
PATH="" run_expect_fail "Missing soroban CLI" "Dependency missing" "./test_data/bad.wasm"

# Failed WASM install (simulate empty output)
SUDO_FAKE_INSTALL_FAIL=1 run_expect_fail "WASM install failure" "Failed to install WASM" "./test_data/bad.wasm"

echo "All deployment failure tests passed!"

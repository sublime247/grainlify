#!/usr/bin/env bash
set -euo pipefail

UPGRADE_SCRIPT="./contracts/scripts/upgrade.sh"

fail() { echo "✘ FAIL: $1"; exit 1; }
pass() { echo "✔ PASS: $1"; }

run_expect_fail() {
    desc="$1"
    expected_msg="$2"
    shift 2

    set +e
    output=$($UPGRADE_SCRIPT "$@" 2>&1)
    exit_code=$?
    set -e

    if [[ $exit_code -eq 0 ]]; then
        fail "$desc (expected failure, got exit 0)"
    fi

    if ! echo "$output" | grep -q "$expected_msg"; then
        echo "$output"
        fail "$desc (expected message '$expected_msg')"
    fi

    pass "$desc"
}

echo "=== Upgrade Script Failure Tests ==="

# Missing contract ID
run_expect_fail "Missing contract ID" "No contract ID specified"

# Invalid contract ID
run_expect_fail "Invalid format" "Contract ID format may be invalid" "BAD_ID" "./test_data/bad.wasm"

# Missing WASM file
run_expect_fail "Missing WASM file" "No WASM file specified" "C123456789..." 

# Nonexistent WASM
run_expect_fail "Invalid WASM file path" "WASM file does not exist" "C123456789..." "missing.wasm"

# Invalid identity
run_expect_fail "Missing identity" "Identity not found" "C123456789..." "./test_data/bad.wasm" --source ghost_id

# Fake install fail
SUDO_FAKE_INSTALL_FAIL=1 run_expect_fail "Install failure" "Failed to install WASM" "C123456789..." "./test_data/bad.wasm"

# Fake upgrade fail
SUDO_FAKE_UPGRADE_FAIL=1 run_expect_fail "Upgrade invocation fails" "Upgrade invocation failed" "C123456789..." "./test_data/bad.wasm"

echo "All upgrade failure tests passed!"

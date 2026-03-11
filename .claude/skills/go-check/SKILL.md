---
name: go-check
description: Run Go tests and linter together. Use when the user wants to verify code quality, run checks before committing, or validate changes after implementation.
allowed-tools: Bash, Read, Grep
---

Run the full Go quality check pipeline: tests + lint.

**Input**: Optionally specify a package path (e.g., `./users/...`). Defaults to all packages.

**Steps**

1. **Run tests and lint in parallel**

   If a specific package is given via `$ARGUMENTS`:
   ```bash
   go test -v -race $ARGUMENTS
   golangci-lint run $ARGUMENTS
   ```

   Otherwise run all:
   ```bash
   go test -v -race ./...
   golangci-lint run
   ```

   Run both commands in parallel using separate Bash tool calls.

2. **Analyze results**

   Parse output for:
   - Test failures (FAIL lines)
   - Lint violations (file:line format)
   - Race conditions detected
   - Build errors

3. **Report summary**

   ```
   ## Go Check Results

   **Tests**: ✓ PASS (or ✗ FAIL)
   - Packages tested: N
   - Tests run: N passed, N failed

   **Lint**: ✓ CLEAN (or ✗ N issues)
   - [list any violations with file:line references]

   **Action needed**: [if any failures, suggest specific fixes]
   ```

4. **On failure, offer to fix**

   If there are lint issues or test failures:
   - For lint: show the specific violations and offer to fix
   - For tests: show the failing test output and analyze the root cause
   - Ask before making any changes

**Guardrails**
- Always run both tests and lint — never skip one
- Use `-race` flag for test runs to catch race conditions
- Report results concisely — don't dump raw output unless failures exist
- Never auto-fix without asking the user first

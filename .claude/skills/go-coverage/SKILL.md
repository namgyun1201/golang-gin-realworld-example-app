---
name: go-coverage
description: Analyze test coverage and identify gaps. Use when the user asks about coverage, wants to improve test coverage, or needs a coverage report.
allowed-tools: Bash, Read, Grep, Glob
---

Generate a test coverage report, identify low-coverage areas, and suggest improvements.

**Input**: Optionally specify a package (e.g., `./users/...`). Defaults to all packages.

**Steps**

1. **Generate coverage profile**

   ```bash
   go test -coverprofile=coverage.out -covermode=atomic ./...
   ```

   If `$ARGUMENTS` specifies a package, scope to that package.

2. **Get per-function coverage**

   ```bash
   go tool cover -func=coverage.out
   ```

3. **Analyze results**

   Parse the output to identify:
   - Overall coverage percentage
   - Per-package coverage breakdown
   - Functions with 0% coverage (untested)
   - Functions below 80% coverage (undertested)

4. **Compare against project targets**

   Read [coverage-targets.md](coverage-targets.md) for package-level targets, test infrastructure details, and high-value test targets.

   Flag any package that dropped below its target.

5. **Report**

   ```
   ## Coverage Report

   **Overall**: XX.X% [↑/↓ vs target 90%]

   | Package  | Coverage | Target | Status |
   |----------|----------|--------|--------|
   | articles | XX.X%    | 93%+   | ✓/✗    |
   | users    | XX.X%    | 99%+   | ✓/✗    |
   | common   | XX.X%    | 85%+   | ✓/✗    |

   ### Gaps (functions below 80%)
   - `package.FunctionName` — XX.X% (file:line)

   ### Untested Functions
   - `package.FunctionName` (file:line)

   ### Suggestions
   - [specific suggestions for improving coverage]
   ```

6. **Clean up**

   ```bash
   rm -f coverage.out
   ```

**Guardrails**
- Always clean up coverage.out after analysis
- Compare against project targets, not arbitrary thresholds
- Suggest specific test cases for gaps, not generic advice
- Don't write tests automatically — suggest and let the user decide

---
name: review-code
description: Review code for quality, GORM v2 patterns, RealWorld spec compliance, and security. Use when the user asks for a code review, wants to check code quality, or before merging changes.
allowed-tools: Bash, Read, Grep, Glob
---

Perform a comprehensive code review focused on this project's specific patterns and requirements.

**Input**: Optionally specify files or a package to review (e.g., `./articles/...`). If omitted, review recently changed files (`git diff`).

**Steps**

1. **Determine scope**

   If `$ARGUMENTS` specifies files/packages, use those.
   Otherwise, get changed files:
   ```bash
   git diff --name-only HEAD~1
   ```
   Or if there are uncommitted changes:
   ```bash
   git diff --name-only
   ```

   Only review `.go` files.

2. **Read the files to review**

   Read each file in the review scope using the Read tool.

3. **Check GORM v2 patterns**

   Read [gorm-v2-patterns.md](gorm-v2-patterns.md) for the full anti-pattern reference.
   Scan for each violation listed (Related→Preload, Update→Updates, Delete pointer, Count type, Association usage).

4. **Check RealWorld API spec compliance**

   Read [realworld-api-spec.md](realworld-api-spec.md) for the complete status code and response format reference.
   For router/handler files, verify status codes, response wrappers, error format, and pagination parameters match the spec.

5. **Check security patterns**

   Read [security-checklist.md](security-checklist.md) for the full checklist.
   Verify authentication, authorization, input validation, password security, and database security patterns.

6. **Check Go best practices**

   - [ ] Errors handled explicitly (no `_` for error returns)
   - [ ] Proper use of `required` validator tag (not deprecated `exists`)
   - [ ] Consistent error response format
   - [ ] No unused imports or variables
   - [ ] Functions are focused and not overly long

7. **Report**

   ```
   ## Code Review: [scope]

   ### Critical Issues
   - [blocking issues that must be fixed]

   ### Warnings
   - [non-blocking but should be addressed]

   ### GORM v2 Patterns
   - ✓/✗ [pattern check results]

   ### RealWorld Spec Compliance
   - ✓/✗ [compliance check results]

   ### Security
   - ✓/✗ [security check results]

   ### Suggestions
   - [optional improvements, not required]

   **Verdict**: ✓ Approved / ⚠ Approved with suggestions / ✗ Changes requested
   ```

**Guardrails**
- Read all files before making any judgments
- Focus on project-specific patterns (GORM v2, RealWorld spec), not generic Go style
- Distinguish between critical issues (must fix) and suggestions (nice to have)
- Reference specific lines with `file:line` format
- Don't suggest changes that aren't related to the review scope
- Never auto-fix — report findings and let the user decide

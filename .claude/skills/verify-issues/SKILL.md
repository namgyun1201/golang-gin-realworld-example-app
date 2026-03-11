---
name: verify-issues
description: Verify acceptance criteria for completed GitHub issues. Use after implementation work to confirm all issue requirements are met before closing.
allowed-tools: Bash, Read, Grep, Glob, Agent
---

Verify that completed work satisfies the acceptance criteria of GitHub issues.

**Input**: Issue numbers to verify (e.g., `1 2 3` or `1-5`). If omitted, checks all open issues assigned to the current branch's work.

**Steps**

1. **Fetch issue details**

   For each issue number in `$ARGUMENTS`:
   ```bash
   gh issue view <number> --json title,body,labels,state
   ```

   Parse the issue body to extract:
   - Problem description
   - Acceptance criteria / solution requirements
   - Affected files

2. **Verify each acceptance criterion**

   For each requirement found in the issue:
   - **Code verification**: Use `Grep` to confirm the fix exists in the codebase
   - **Negative verification**: Confirm the problematic pattern is removed (e.g., hardcoded secrets, unsafe code)
   - **Build check**: Run `go build ./...` to verify compilation
   - **Test check**: Run `go test ./...` to verify tests pass
   - **Static analysis**: Run `go vet ./...` for correctness

   Run build, test, and vet in parallel.

3. **Report per-issue verification**

   For each issue, produce a checklist:
   ```
   ## Issue #N: <title>

   | Criterion | Evidence | Result |
   |-----------|----------|--------|
   | <requirement> | <file:line or grep result> | PASS/FAIL |

   **Verdict**: PASS / FAIL (with details)
   ```

4. **Summary**

   ```
   ## Verification Summary

   | Issue | Title | Verdict |
   |-------|-------|---------|
   | #N | <title> | PASS/FAIL |

   **Build**: PASS/FAIL
   **Tests**: PASS/FAIL
   **Vet**: PASS/FAIL

   Ready to close: #N, #N
   Needs attention: #N (reason)
   ```

5. **Close verified issues (with user confirmation)**

   If all criteria pass for an issue, ask the user before closing:
   ```bash
   gh issue close <number> --comment "Verified and resolved in <commit/branch>."
   ```

**Guardrails**
- Never close issues without explicit user confirmation
- Always verify with code evidence (grep/read), not assumptions
- Run full build + test + vet before declaring any issue verified
- If an issue has no clear acceptance criteria in the body, extract them from the title and problem description
- Report FAIL immediately if any criterion is not met, with details on what's missing

#!/bin/bash
# check-readme.sh — PostToolUse hook for Edit/Write
# Detects source file changes and reminds Claude to update README.md

input=$(cat)
tool_name=$(echo "$input" | jq -r '.tool_name // empty')

# Extract file path based on tool type
file_path=""
if [ "$tool_name" = "Write" ] || [ "$tool_name" = "Edit" ]; then
  file_path=$(echo "$input" | jq -r '.tool_input.file_path // empty')
fi

# No file path — skip
if [ -z "$file_path" ]; then
  exit 0
fi

# Skip if editing README itself or non-project files
case "$file_path" in
  *README.md*|*CLAUDE.md*|*AGENTS.md*|*.claude/*|*.omc/*|*openspec/*|*coverage.out*)
    exit 0
    ;;
esac

# Only care about source files that affect README content
should_remind=false
case "$file_path" in
  *.go|*go.mod|*go.sum|*.env.example|*Dockerfile|*docker-compose*|*.github/workflows/*)
    should_remind=true
    ;;
esac

if [ "$should_remind" = true ]; then
  # Get the basename for a concise message
  basename=$(basename "$file_path")
  cat <<MSG
[README-SYNC] Source file modified: $basename
If this change affects any of the following, update README.md accordingly:
- Project structure or directory layout
- Dependencies (go.mod)
- API endpoints or routes
- Environment variables or configuration
- Build/run/test commands
- Test coverage numbers
Do NOT update README for minor internal refactors or bug fixes that don't change the public interface.
MSG
fi

exit 0

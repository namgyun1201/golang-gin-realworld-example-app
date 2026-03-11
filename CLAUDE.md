# Vibe Coding Workflow

This project uses **openspec** for plan-centered development. All feature work follows the explore → propose → review → implement → archive cycle.

## Development Cycle

1. **Explore** (`/opsx:explore`): Investigate the codebase, understand the problem, think through ideas
2. **Propose** (`/opsx:propose <change-name>`): Generate a change with proposal.md, design.md, tasks.md artifacts
3. **Review**: Human reviews artifacts in `openspec/changes/<change-name>/`
4. **Implement** (`/opsx:apply <change-name>`): Execute tasks from the approved change
5. **Archive** (`/opsx:archive <change-name>`): Archive completed changes and update specs

**Never implement features directly** — always go through the openspec propose/apply cycle.

## Project Commands

```bash
# Development
go run hello.go              # Start server (localhost:3000)
go test ./...                # Run all tests
go test -coverprofile=coverage.out ./...  # Tests with coverage
golangci-lint run            # Lint check
go fmt ./...                 # Format code

# OpenSpec
openspec list                # List active changes
openspec status --change <name>  # Check change status
openspec validate            # Validate config and changes
```

## Key Rules

- Read `AGENTS.md` for coding conventions and GORM v2 patterns before making changes
- Run `go test ./...` after every code change
- Run `golangci-lint run` before committing
- Keep changes scoped to one package (users/, articles/, common/) when possible
- Follow the RealWorld API spec for all endpoint changes
- Use conventional commits for git messages

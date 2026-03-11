.PHONY: test test-unit test-integration test-e2e test-all test-coverage lint fmt vet clean clean-db

# Go binary (auto-detect)
GO ?= go

# Default: run unit tests only
test: test-unit

# Clean test DB files (needed between test stages to avoid cross-contamination)
clean-db:
	@find . -name "*.db" -path "*/data/*" -delete 2>/dev/null || true

# Unit tests (no build tags, fast)
# -p 1: run packages sequentially (shared SQLite DB file)
test-unit: clean-db
	$(GO) test -p 1 ./... -count=1 -v

# Integration tests (DB-dependent, moderate speed)
test-integration: clean-db
	$(GO) test -p 1 -tags integration ./... -count=1 -v

# E2E tests (full API flow, slower)
# E2E tests live in root package only; run unit tests separately to avoid cross-tag issues
test-e2e: clean-db
	$(GO) test -p 1 -tags e2e -count=1 -v .
	$(GO) test -p 1 ./... -count=1

# All tests combined
test-all: test-unit test-integration test-e2e

# Coverage report (unit tests)
test-coverage:
	$(GO) test -p 1 -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out
	@echo ""
	@echo "To view HTML report: go tool cover -html=coverage.out"

# Coverage with all tags
test-coverage-all:
	$(GO) test -p 1 -tags "integration" -coverprofile=coverage.out ./...
	$(GO) test -p 1 -tags "e2e" -coverprofile=coverage_e2e.out .
	$(GO) tool cover -func=coverage.out
	@echo ""
	@echo "To view HTML report: go tool cover -html=coverage.out"

# Lint
lint:
	golangci-lint run

# Lint with auto-fix
lint-fix:
	golangci-lint run --fix

# Format code
fmt:
	$(GO) fmt ./...

# Vet
vet:
	$(GO) vet ./...
	$(GO) vet -tags integration ./...
	$(GO) vet -tags e2e ./...

# Run all checks (CI-friendly)
ci: fmt vet lint test-all

# Clean test artifacts
clean: clean-db
	rm -f coverage.out coverage_e2e.out
	$(GO) clean -testcache

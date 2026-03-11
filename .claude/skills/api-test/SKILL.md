---
name: api-test
description: Run RealWorld API end-to-end tests using Newman/Postman. Use when the user wants to run E2E tests, validate API endpoints, or check RealWorld spec compliance.
disable-model-invocation: true
allowed-tools: Bash, Read
---

Run the RealWorld API E2E test suite against a running server.

**Input**: Optionally specify the API URL (default: `http://localhost:3000/api`).

**Steps**

1. **Check prerequisites**

   Verify Newman is installed:
   ```bash
   which newman || npm list -g newman
   ```

   If not installed, inform the user:
   ```
   Newman is required for API tests. Install with:
   npm install -g newman
   ```

2. **Check if server is running**

   ```bash
   curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/api/tags || true
   ```

   If server is not running, offer two options:
   - Start the server in background: `go run hello.go &`
   - Let the user start it manually

3. **Run the API test suite**

   Use the project's test runner script if available:
   ```bash
   ./scripts/run-api-tests.sh
   ```

   Or run Newman directly with the Postman collection:
   ```bash
   newman run api/Conduit.postman_collection.json \
     --environment api/Conduit.postman_environment.json \
     --reporters cli
   ```

   If `$ARGUMENTS` provides a custom URL, pass it as the `APIURL` environment variable.

4. **Analyze results**

   Parse Newman output for:
   - Total requests executed
   - Assertions passed/failed
   - Specific endpoint failures
   - Response time anomalies

5. **Report**

   ```
   ## API Test Results

   **Status**: ✓ PASS / ✗ FAIL
   **Requests**: N executed
   **Assertions**: N passed, N failed

   ### Failed Endpoints (if any)
   - `METHOD /api/path` — assertion failure description
     Expected: X, Got: Y

   ### Summary
   [brief assessment of API spec compliance]
   ```

6. **On failure, diagnose**

   Read [realworld-endpoints.md](realworld-endpoints.md) for the complete endpoint reference with expected status codes and response formats.

   For each failing endpoint:
   - Check the corresponding router handler in `users/routers.go` or `articles/routers.go`
   - Check the serializer output format
   - Compare against the endpoint reference expectations
   - Suggest specific fixes

**Guardrails**
- This skill requires a running server — never skip this check
- Don't modify the Postman collection or environment files
- Report failures with enough context to diagnose (request, expected, actual)
- If the server needs to be started, always offer the user the choice

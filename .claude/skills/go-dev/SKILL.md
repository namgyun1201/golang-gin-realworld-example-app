---
name: go-dev
description: Start the development server. Use when the user wants to run the app, start the server, or test endpoints locally.
disable-model-invocation: true
allowed-tools: Bash, Read
---

Start the Go development server with health verification.

**Input**: Optionally specify a port (default: from `PORT` env var or 3000).

**Steps**

1. **Check for existing server**

   ```bash
   lsof -i :3000 -t 2>/dev/null || true
   ```

   If a process is already running on the port:
   - Inform the user: "Server already running on port 3000 (PID: XXXX)"
   - Ask if they want to restart (kill existing + start new)

2. **Build check**

   ```bash
   go build ./...
   ```

   If build fails, report errors and stop. Don't start a broken server.

3. **Start the server**

   ```bash
   go run hello.go &
   ```

   Run in background so the user can continue working.

4. **Verify startup**

   Wait briefly, then check health:
   ```bash
   sleep 2 && curl -s http://localhost:3000/api/tags
   ```

   If health check fails, check server logs for errors.

5. **Report**

   ```
   ## Dev Server Started

   **URL**: http://localhost:3000
   **API Base**: http://localhost:3000/api
   **PID**: XXXX

   Quick test endpoints:
   - GET  http://localhost:3000/api/tags
   - POST http://localhost:3000/api/users (register)
   - POST http://localhost:3000/api/users/login (login)

   To stop: kill XXXX
   ```

**Guardrails**
- Always check for existing processes on the port first
- Always verify build succeeds before starting
- Always confirm the server is healthy after starting
- Never kill an existing process without asking the user

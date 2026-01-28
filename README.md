# OPA Test Project

This repository contains experiments and implementations using Open Policy Agent (OPA) with Go.

## Folder Structure

### `01_latency`
Contains a benchmark suite for measuring OPA policy evaluation latency.
- **Goal:** Assess performance with large datasets and complex rules.
- **Key File:** `main.go` runs the loop of policy evaluations.

### `02_go_server` (Simple Go Server)
A standalone Go HTTP server integrated with OPA for request authorization.
- **Goal:** Demonstrate how to embed OPA in a Go web server.
- **Key Features:**
    - `POST /protected`: Protected endpoint requiring specific user headers.
    - `GET /public`: Publicly accessible endpoint.
    - Uses `policy/policy.rego` for authorization logic.

## Usage
Navigate to each folder and follow the instructions in their respective `README.md` files (if available) or run:
```bash
go run main.go
```
(Ensure you run `go mod tidy` if needed).

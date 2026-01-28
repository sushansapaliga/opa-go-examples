# Simple Go Server with OPA Policies

This project demonstrates a simple Go HTTP server integrated with Open Policy Agent (OPA) for request authorization.

## stricture
- `main.go`: The HTTP server and OPA integration logic.
- `policy/policy.rego`: The OPA policy definitions.
- `policy/policy_test.rego`: Tests for the OPA policy.

## Policy Loading

The server is configured to load all Rego files from the `policy/` directory **except** files ending in `_test.rego`. This ensures that test cases are not evaluated as part of the live authorization logic.

## Running the Server

1. Initialize the module (if not already done):
   ```bash
   go mod tidy
   ```
2. Run the server:
   ```bash
   go run main.go
   ```

## API Endpoints

- `GET /public`: Publicly accessible endpoint.
- `POST /protected`: Protected endpoint. Requires `X-User: admin` header.

## content Testing

To run the OPA policy tests:

```bash
cd policy
opa test -v .
```
(Requires `opa` CLI installed)

## Curl Examples

**Public Access:**
```bash
curl -v http://localhost:8080/public
```

**Admin Access:**
```bash
curl -v -H "X-User: admin" -X POST http://localhost:8080/protected
```

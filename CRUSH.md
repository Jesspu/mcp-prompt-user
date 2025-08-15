# CRUSH Test

This file outlines the coding standards and conventions for the project.

---

### Commands

- **Build:** `go build ./...`
- **Test:** `go test ./...`
- **Test a single function:** `go test -run ^TestMyFunction$`
- **Lint:** `golangci-lint run`

---

### Code Style

- **Formatting:** Use `gofmt` or `goimports` to format code.
- **Imports:** Group imports into three blocks: standard library, third-party packages, and internal packages.
- **Naming:** Follow Go's idiomatic naming conventions (e.g., `camelCase` for variables, `PascalCase` for exported identifiers).
- **Error Handling:** Use `if err != nil` for error handling. Avoid panicking in library code.
- **Types:** Use structs to represent complex data types.
- **Comments:** Add comments to explain complex logic and document exported functions.
- **Testing:** Write unit tests for all new features and bug fixes.

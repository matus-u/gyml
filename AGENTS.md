# AGENT GUIDELINES FOR `gyml`

This document provides essential guidelines for agentic coding agents operating within the `gyml` repository. Adhering to these conventions ensures consistency, maintainability, and high quality of the codebase.

## 1. Project Overview

`gyml` is a Go library designed for generic YAML file manipulation using path-based lookups.

## 2. Build, Lint, and Test Commands

Agents should always verify changes using the project's standard build, lint, and test commands.

### Build
To compile the entire project:
```bash
go build ./...
```
To compile the current directory:
```bash
go build .
```

### Linting
The primary linting tool for Go projects is `go vet`. It helps identify suspicious constructs.
```bash
go vet ./...
```
For more comprehensive static analysis, consider `golangci-lint` if it were configured. As of now, `go vet` is the default.

### Testing
To run all tests in the project:
```bash
go test ./...
```

### Running a Single Test
To execute a specific test function (e.g., `TestParseYAMLFile` in `gyml_test.go`), use the `-run` flag:
```bash
go test -run TestParseYAMLFile ./...
```
Replace `TestParseYAMLFile` with the exact name of the test function you intend to run. The `./...` ensures all packages are considered for test discovery.

## 3. Code Style Guidelines

Go has strong, opinionated style guidelines, many of which are enforced by official tools.

### Imports
- **Grouping:** Imports should be grouped with standard library packages first, followed by third-party packages, and then internal project packages. Each group should be separated by a blank line.
- **Formatting:** The `goimports` tool (which includes `go fmt`) automatically handles import organization and formatting.

Example:
```go
package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"gyml/internal/config"
)
```

### Formatting
- **Tool:** Use `go fmt` to automatically format Go source code. Agents must run `go fmt ./...` after any code modifications.
```bash
go fmt ./...
```
- **Line Length:** While `go fmt` doesn't strictly enforce line length, aim for readability and avoid excessively long lines (typically under 120 characters).

### Types
- **Static Typing:** Go is a statically typed language. Always use appropriate types and ensure type compatibility.
- **Clarity:** Prefer clear, explicit type declarations.
- **Structs:** Use structs for composite data types.

### Naming Conventions
- **Exported Identifiers:** Identifiers that are exported (visible outside their package) must start with an uppercase letter (e.g., `FunctionName`, `StructName`).
- **Unexported Identifiers:** Identifiers that are unexported (private to their package) must start with a lowercase letter (e.g., `functionName`, `structName`).
- **Acronyms:** Acronyms (like `URL`, `HTTP`, `API`) should be all uppercase when exported and all lowercase when unexported, without mixing cases (e.g., `ServeHTTP`, `parseURL`).
- **Variables:** Use concise but descriptive names. For loop variables, single letters (e.g., `i`, `j`, `k`) are common. Error variables are typically `err`.
- **Packages:** Package names should be lowercase, single-word, and reflect the package's purpose.

### Error Handling
- **Explicit Returns:** Go handles errors by returning an `error` as the last return value. Always check for errors immediately.
- **Propagating Errors:** Errors should generally be propagated up the call stack until they can be handled meaningfully.
- **Error Values:** Avoid discarding errors. If an error is truly ignorable, explicitly assign it to `_`.
- **Context:** When returning an error, provide sufficient context, often by wrapping errors using `fmt.Errorf` with `%w`.

Example:
```go
func ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", path, err)
	}
	return data, nil
}
```

## 4. Cursor/Copilot Rules

No specific Cursor rules (`.cursor/rules/` or `.cursorrules`) or Copilot instructions (`.github/copilot-instructions.md`) were found in this repository. Agents should default to the Go community best practices and the guidelines outlined above.

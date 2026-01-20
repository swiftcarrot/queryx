# AGENTS.md - Queryx Codebase Guide

## Project Overview

Queryx is a schema-first, type-safe ORM for Go and TypeScript. It parses HCL schema files, generates type-safe ORM code, and manages database migrations using Atlas. Supports PostgreSQL, MySQL, and SQLite.

**Key Directories:**

- `cmd/queryx/` - CLI application (uses Cobra)
- `schema/` - HCL schema parsing and model building
- `generator/client/{golang,typescript}/` - Template-based code generation
- `adapter/` - Database adapter implementations
- `inflect/` - String inflection utilities (camelCase, snake_case, etc.)
- `internal/integration/` - Integration tests

---

## Build/Test/Lint Commands

### Build & Install

```bash
make fmt              # Format all Go files with gofmt
make build            # Build queryx binary to bin/
make install          # Install to /usr/local/bin (requires sudo)
make clean            # Remove bin/queryx

# Manual build
go build -ldflags "-X github.com/swiftcarrot/queryx/cmd/queryx/action.Version=`git rev-parse HEAD`" \
  -o bin/queryx cmd/queryx/main.go
```

### Testing

```bash
# Run unit tests for specific package (single test)
go test ./schema/...                    # Test schema package
go test ./inflect/...                   # Test inflect package
go test -v ./adapter/...                # Verbose output
go test -run TestSnake ./inflect/...    # Run specific test

# Run all unit tests (excluding generated code)
go test $(go list ./... | grep -Ev "generator|internal") -race -coverprofile=coverage.txt

# Integration tests (requires database setup)
make test-sqlite          # SQLite integration tests
make test-postgresql      # PostgreSQL integration tests
make test-mysql           # MySQL integration tests
make test                 # All integration tests

# TypeScript tests (from internal/integration/)
cd internal/integration && yarn test
cd internal/integration && yarn vitest run --dir client
```

### Linting & Formatting

```bash
# Go formatting
gofmt -w $(find . -name '*.go')
make fmt

# Go linting (used in CI)
golangci-lint run --timeout=3m

# Format schema files
queryx format --schema schema.hcl
queryx fmt --schema schema.hcl
```

---

## Code Style Guidelines

### Go Import Organization

Use 3 groups separated by blank lines:

```go
import (
    // 1. Standard library
    "context"
    "fmt"

    // 2. External dependencies
    "github.com/hashicorp/hcl/v2"
    "github.com/spf13/cobra"

    // 3. Internal packages
    "github.com/swiftcarrot/queryx/adapter"
    "github.com/swiftcarrot/queryx/inflect"
)
```

### Naming Conventions

**Files:**

- Snake_case: `adapter_test.go`, `string_array.go`
- Database-specific: `adapter.postgresql.go`, `adapter.mysql.go`
- Templates: `*.gotmpl` (Go), `*.tstmpl` (TypeScript)
- Tests: `*_test.go` (Go), `*.test.ts` (TypeScript)

**Go Types & Functions:**

- Exported types: PascalCase (`SelectStatement`, `PostgreSQLAdapter`)
- Unexported types: camelCase (`selectStatement`, `envKey`)
- Exported functions: PascalCase (`NewAdapter()`, `CreateDatabase()`)
- Unexported functions: camelCase (`newSchema()`, `createFile()`)
- Constructors: `New<Type>()` or `new<Type>()`
- Interfaces: No "I" prefix (`Adapter`, `DB`, `Migrator`)
- Receiver names: Single letter (`s *SelectStatement`, `c *Config`)

**Variables:**

- camelCase: `schemaFile`, `database`, `rawURL`
- Exported acronyms: `URL`, `ID`, `UUID`, `HTTP`
- Unexported acronyms: `url`, `id`, `uuid`

**TypeScript:**

- Files: kebab-case (`select.ts`, `clause.test.ts`)
- Classes: PascalCase (`SelectStatement`, `QXClient`)
- Functions: camelCase (`newSelect()`, `queryUser()`)
- Private fields: underscore prefix (`_selection`, `_from`)

### Error Handling Pattern

Check immediately, return early:

```go
// Pattern 1: Return error directly
func (a *PostgreSQLAdapter) Open() error {
    db, err := sql.Open("postgres", a.Config.URL)
    if err != nil {
        return err
    }
    a.DB = db
    return nil
}

// Pattern 2: Wrap with context
func LoadTemplates(src embed.FS) error {
    if err := fs.WalkDir(src, "templates", walkFunc); err != nil {
        return fmt.Errorf("error walking directory: %w", err)
    }
    return nil
}

// Pattern 3: Fatal in init/main
func init() {
    client, err := db.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    c = client
}
```

### Comment Style

```go
// Package comments at top of file
// Package adapter provides database adapter implementations
package adapter

// Function comments describe what, not how
// NewAdapter creates a database adapter based on the configuration
func NewAdapter(cfg *schema.Config) (Adapter, error) { ... }

// Inline comments explain why
// TODO: support default values
Default interface{}
```

### File Organization

```go
// 1. Package declaration
package queryx

// 2. Import block (3 groups)
import (...)

// 3. Type definitions
type SelectStatement struct { ... }

// 4. Constructor functions
func NewSelect() *SelectStatement { ... }

// 5. Methods grouped by receiver
func (s *SelectStatement) Select(...) *SelectStatement { ... }
func (s *SelectStatement) Where(...) *SelectStatement { ... }

// 6. Helper functions
func rebind(query string) (string, []interface{}) { ... }
```

---

## Common Patterns

### Template-Based Code Generation

Uses `embed.FS` + Go templates with custom functions from `inflect.TemplateFunctions`:

```go
//go:embed templates
var templatesFS embed.FS

tmpl := template.New(name).Funcs(inflect.TemplateFunctions)
```

### Adapter Pattern

Factory pattern for database abstraction (`adapter/adapter.go`):

```go
func NewAdapter(cfg *schema.Config) (Adapter, error) {
    switch cfg.Adapter {
    case "postgresql": return NewPostgreSQLAdapter(config), nil
    case "mysql": return NewMySQLAdapter(config), nil
    case "sqlite": return NewSQLiteAdapter(config), nil
    }
}
```

### Query Builder (Fluent Interface)

Method chaining pattern - always return `this` or `*self`:

```go
func (s *SelectStatement) Where(clauses ...*Clause) *SelectStatement {
    s.where = clauses[0].And(clauses[1:]...)
    return s
}
```

### Testing Patterns

```go
// Table-driven tests
func TestSnake(t *testing.T) {
    require.Equal(t, "user", Snake("User"))
    require.Equal(t, "user_post", Snake("UserPost"))
}

// Integration tests with init
var c *db.QXClient

func init() {
    client, err := db.NewClient()
    if err != nil {
        log.Fatal(err)
    }
    c = client
}
```

---

**Generated:** 2026-01-20
**Go Version:** 1.16+
**Main Dependencies:** Cobra (CLI), HCL (schema), Atlas (migrations), Testify (testing)

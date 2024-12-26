package adapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/swiftcarrot/queryx/schema"
)

type Adapter interface {
	Open() error
	Close() error
	CreateDatabase() error
	DropDatabase() error
	DumpSchema() (string, error)
	CreateMigrationsTable(ctx context.Context) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryVersion(ctx context.Context, version string) (*sql.Rows, error)
}

func NewAdapter(cfg *schema.Config) (Adapter, error) {
	config := NewConfig(cfg)
	if config.Adapter == "postgresql" {
		return NewPostgreSQLAdapter(config), nil
	} else if config.Adapter == "mysql" {
		return NewMySQLAdapter(config), nil
	} else if config.Adapter == "sqlite" {
		return NewSQLiteAdapter(config), nil
	}

	// TODO: list supported adapters
	return nil, fmt.Errorf("unsupported adapter: %q", config.Adapter)
}

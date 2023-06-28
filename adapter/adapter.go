package adapter

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/swiftcarrot/queryx/schema"
)

type Adapter interface {
	Open() error
	CreateDatabase() error
	DropDatabase() error
	CreateMigrationsTable(ctx context.Context) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	GetAdapter() string
	GetDBName() string
}

func NewAdapter(config *schema.Config) (Adapter, error) {
	if config.Adapter == "postgresql" {
		return NewPostgreSQLAdapter(config), nil
	} else if config.Adapter == "mysql" {
		return NewMysqlAdapter(config), nil
	} else if config.Adapter == "sqlite" {
		return NewSQLiteAdapter(config), nil
	}
	return nil, fmt.Errorf("unsupported adapter: %q", config.Adapter)
}

package adapter

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	sqlschema "ariga.io/atlas/sql/schema"
	"ariga.io/atlas/sql/sqlite"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteAdapter struct {
	*sql.DB
	Config *Config
}

func NewSQLiteAdapter(config *Config) *SQLiteAdapter {
	return &SQLiteAdapter{
		Config: config,
	}
}

func (a *SQLiteAdapter) Open() error {
	db, err := sql.Open("sqlite3", a.Config.URL)
	if err != nil {
		return err
	}
	a.DB = db
	return nil
}

func (a *SQLiteAdapter) CreateDatabase() error {
	_, err := os.Create(a.Config.Database)
	if err != nil {
		return err
	}

	fmt.Println("Created database", a.Config.Database)

	return nil
}

func (a *SQLiteAdapter) DropDatabase() error {
	err := os.Remove(a.Config.Database)
	if err != nil {
		return err
	}
	fmt.Println("Dropped database", a.Config.Database)
	return nil
}

// create migrations table with atlas
func (a *SQLiteAdapter) CreateMigrationsTable(ctx context.Context) error {
	driver, err := sqlite.Open(a)
	if err != nil {
		return err
	}

	from, err := driver.InspectSchema(ctx, "main", &sqlschema.InspectOptions{
		Tables: []string{"schema_migrations"},
	})
	if err != nil {
		return err
	}

	version := sqlschema.NewStringColumn("version", "varchar")
	to := sqlschema.New("main").AddTables(
		sqlschema.NewTable("schema_migrations").AddColumns(
			sqlschema.NewStringColumn("version", "varchar"),
		).SetPrimaryKey(sqlschema.NewPrimaryKey(version)))

	changes, err := driver.SchemaDiff(from, to)
	if err != nil {
		return err
	}

	if err := driver.ApplyChanges(ctx, changes); err != nil {
		return err
	}

	return nil
}

func (a *SQLiteAdapter) QueryVersion(ctx context.Context, version string) (*sql.Rows, error) {
	return a.DB.QueryContext(ctx, "select version from schema_migrations where version = $1", version)
}

package adapter

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	sqlschema "ariga.io/atlas/sql/schema"
	"ariga.io/atlas/sql/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"github.com/swiftcarrot/queryx/schema"
)

type SQLiteAdapter struct {
	*sql.DB
	Config *schema.Config
}

func NewSQLiteAdapter(config *schema.Config) *SQLiteAdapter {
	return &SQLiteAdapter{
		Config: config,
	}
}

func (a *SQLiteAdapter) Open() error {
	db, err := sql.Open("sqlite3", a.Config.ConnectionURL(true))
	if err != nil {
		return err
	}
	a.DB = db
	return nil
}

func (a *SQLiteAdapter) CreateDatabase() error {
	// fmt.Println("env...", os.Getenv(a.Config.URL.EnvKey))
	// _, err := os.Create(os.Getenv(a.Config.URL.EnvKey))
	// if err != nil {
	// 	return err
	// }

	// fmt.Println("Created database", a.Config.Database)

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
func (a *SQLiteAdapter) GetAdapter() string {
	return a.Config.Adapter
}

func (a *SQLiteAdapter) GetDBName() string {
	return ""
}

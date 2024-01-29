package adapter

import (
	"context"
	"database/sql"
	"fmt"

	"ariga.io/atlas/sql/mysql"
	sqlschema "ariga.io/atlas/sql/schema"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLAdapter struct {
	*sql.DB
	Config *Config
}

func NewMySQLAdapter(config *Config) *MySQLAdapter {
	return &MySQLAdapter{
		Config: config,
	}
}

func (a *MySQLAdapter) Open() error {
	db, err := sql.Open("mysql", a.Config.URL)
	if err != nil {
		return err
	}
	a.DB = db
	return nil
}

func (a *MySQLAdapter) CreateDatabase() error {
	db, err := sql.Open("mysql", a.Config.URL2)
	if err != nil {
		return err
	}
	defer db.Close()

	a.DB = db
	sql := fmt.Sprintf("CREATE DATABASE %s", a.Config.Database)
	_, err = a.ExecContext(context.Background(), sql)
	if err != nil {
		return err
	}

	fmt.Println("Created database", a.Config.Database)

	return nil
}

func (a *MySQLAdapter) DropDatabase() error {
	db, err := sql.Open("mysql", a.Config.URL2)
	if err != nil {
		return err
	}
	defer db.Close()

	a.DB = db
	sql := fmt.Sprintf("DROP DATABASE IF EXISTS %s", a.Config.Database)
	_, err = a.ExecContext(context.Background(), sql)
	if err != nil {
		return err
	}

	fmt.Println("Dropped database", a.Config.Database)

	return nil
}

// create migrations table with atlas
func (a *MySQLAdapter) CreateMigrationsTable(ctx context.Context) error {
	driver, err := mysql.Open(a)
	if err != nil {
		return err
	}

	from, err := driver.InspectSchema(ctx, a.Config.Database, &sqlschema.InspectOptions{
		Tables: []string{"schema_migrations"},
	})
	if err != nil {
		return err
	}
	version := sqlschema.NewStringColumn("version", "varchar(256)")
	to := sqlschema.New(a.Config.Database).AddTables(
		sqlschema.NewTable("schema_migrations").AddColumns(
			sqlschema.NewStringColumn("version", "varchar(256)"),
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

func (a *MySQLAdapter) QueryVersion(ctx context.Context, version string) (*sql.Rows, error) {
	return a.DB.QueryContext(ctx, "select version from schema_migrations where version = ?", version)
}

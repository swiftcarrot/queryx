package adapter

import (
	"context"
	"database/sql"
	"fmt"

	"ariga.io/atlas/sql/postgres"
	sqlschema "ariga.io/atlas/sql/schema"
	_ "github.com/lib/pq"
	"github.com/swiftcarrot/queryx/schema"
)

type PostgreSQLAdapter struct {
	*sql.DB
	Config *schema.Config
}

func NewPostgreSQLAdapter(config *schema.Config) *PostgreSQLAdapter {
	return &PostgreSQLAdapter{
		Config: config,
	}
}

func (a *PostgreSQLAdapter) Open() error {
	db, err := sql.Open("postgres", a.Config.ConnectionURL(true))
	if err != nil {
		return err
	}
	a.DB = db
	return nil
}

func (a *PostgreSQLAdapter) CreateDatabase() error {
	db, err := sql.Open("postgres", a.Config.ConnectionURL(false))
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

func (a *PostgreSQLAdapter) DropDatabase() error {
	db, err := sql.Open("postgres", a.Config.ConnectionURL(false))
	if err != nil {
		return err
	}
	defer db.Close()

	a.DB = db
	sql := fmt.Sprintf("DROP DATABASE %s", a.Config.Database)
	_, err = a.ExecContext(context.Background(), sql)
	if err != nil {
		return err
	}

	fmt.Println("Dropped database", a.Config.Database)

	return nil
}

// create migrations table with atlas
func (a *PostgreSQLAdapter) CreateMigrationsTable(ctx context.Context) error {
	driver, err := postgres.Open(a)
	if err != nil {
		return err
	}

	from, err := driver.InspectSchema(ctx, "public", &sqlschema.InspectOptions{
		Tables: []string{"schema_migrations"},
	})
	if err != nil {
		return err
	}

	version := sqlschema.NewStringColumn("version", "varchar")
	to := sqlschema.New("public").AddTables(
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

func (a *PostgreSQLAdapter) GetDBName() string {
	a.Config.ConnectionURL(false)
	return a.Config.Database
}

func (a *PostgreSQLAdapter) GetAdapter() string {
	return a.Config.Adapter
}

package adapter

import (
	"context"
	"database/sql"
	"fmt"

	my "ariga.io/atlas/sql/mysql"
	sqlschema "ariga.io/atlas/sql/schema"
	_ "github.com/go-sql-driver/mysql"
	"github.com/swiftcarrot/queryx/schema"
)

type MysqlAdapter struct {
	*sql.DB
	Config *schema.Config
}

func NewMysqlAdapter(config *schema.Config) *MysqlAdapter {
	return &MysqlAdapter{
		Config: config,
	}
}

func (a *MysqlAdapter) Open() error {
	db, err := sql.Open("mysql", a.Config.ConnectionURL(true))
	if err != nil {
		return err
	}
	a.DB = db
	return nil
}

func (a *MysqlAdapter) CreateDatabase() error {
	db, err := sql.Open("mysql", a.Config.ConnectionURL(false))
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

func (a *MysqlAdapter) DropDatabase() error {
	db, err := sql.Open("mysql", a.Config.ConnectionURL(false))
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
func (a *MysqlAdapter) CreateMigrationsTable(ctx context.Context) error {
	driver, err := my.Open(a)
	if err != nil {
		return err
	}

	dbName := a.GetDBName()
	if dbName == "" {
		return fmt.Errorf("the ad can not be nil")
	}
	from, err := driver.InspectSchema(ctx, dbName, &sqlschema.InspectOptions{
		Tables: []string{"schema_migrations"},
	})
	if err != nil {
		return err
	}
	_ = a.Config.ConnectionURL(false)
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

func (a *MysqlAdapter) GetDBName() string {
	a.Config.ConnectionURL(false)
	return a.Config.Database
}

func (a *MysqlAdapter) GetAdapter() string {
	return a.Config.Adapter
}
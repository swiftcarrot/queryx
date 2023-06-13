package action

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"ariga.io/atlas/sql/migrate"
	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/postgres"
	sqlschema "ariga.io/atlas/sql/schema"
	"ariga.io/atlas/sql/sqlite"
	"github.com/spf13/cobra"
	"github.com/swiftcarrot/queryx/adapter"
	"github.com/swiftcarrot/queryx/schema"
)

func newAdapter() (adapter.Adapter, error) {
	sch, err := newSchema()
	if err != nil {
		return nil, err
	}

	environment := os.Getenv("QUERYX_ENV")
	if environment == "" {
		environment = "development"
	}

	return adapter.NewAdapter(sch.Databases[0].LoadConfig(environment))
}

// TODO: should check unapplied migrations
func findSchemaChanges(adapter adapter.Adapter, database *schema.Database) ([]*migrate.Change, error) {
	ctx := context.Background()

	// TODO: support other drivers in atlas
	driver, err := postgres.Open(adapter)
	if err != nil {
		return nil, err
	}

	// TODO: multiple schema support
	from, err := driver.InspectSchema(ctx, "public", &sqlschema.InspectOptions{
		Tables: database.Tables(),
	})
	if err != nil {
		return nil, err
	}
	desired := database.CreateSQLSchema()
	changes, err := driver.SchemaDiff(from, desired)
	if err != nil {
		return nil, err
	}

	for _, c := range changes {
		if d, ok := c.(*sqlschema.DropTable); ok {
			d.Extra = append(d.Extra, &sqlschema.IfExists{})
		}
	}

	plan, err := driver.PlanChanges(ctx, "plan", changes)
	if err != nil {
		return nil, err
	}

	// TODO: warning for falsy plan.Reversible?

	return plan.Changes, nil
}

// TODO: should check unapplied migrations
func findMYSQLSchemaChanges(adapter adapter.Adapter, database *schema.Database) ([]*migrate.Change, error) {
	ctx := context.Background()

	// TODO: support other drivers in atlas
	driver, err := mysql.Open(adapter)
	if err != nil {
		return nil, err
	}
	environment := os.Getenv("QUERYX_ENV")
	if environment == "" {
		environment = "development"
	}
	config := database.LoadConfig(environment)
	dbName := config.GetDatabaseName()
	if dbName == "" {
		return nil, fmt.Errorf("the database no exists")
	}
	from, err := driver.InspectSchema(ctx, dbName, &sqlschema.InspectOptions{
		Tables: database.Tables(),
	})
	if err != nil {
		return nil, err
	}
	desired := database.CreateMYSQLSchema(dbName)
	changes, err := driver.SchemaDiff(from, desired)
	if err != nil {
		return nil, err
	}

	for _, c := range changes {
		if d, ok := c.(*sqlschema.DropTable); ok {
			d.Extra = append(d.Extra, &sqlschema.IfExists{})
		}
	}

	plan, err := driver.PlanChanges(ctx, "plan", changes)
	if err != nil {
		return nil, err
	}

	// TODO: warning for falsy plan.Reversible?

	return plan.Changes, nil
}

func findSQLiteSchemaChanges(adapter adapter.Adapter, database *schema.Database) ([]*migrate.Change, error) {
	ctx := context.Background()

	// TODO: support other drivers in atlas
	driver, err := sqlite.Open(adapter)
	if err != nil {
		return nil, err
	}
	from, err := driver.InspectSchema(ctx, "main", &sqlschema.InspectOptions{
		Tables: database.Tables(),
	})
	if err != nil {
		return nil, err
	}
	desired := database.CreateSQLiteSchema("main")
	changes, err := driver.SchemaDiff(from, desired)
	if err != nil {
		return nil, err
	}

	for _, c := range changes {
		if d, ok := c.(*sqlschema.DropTable); ok {
			d.Extra = append(d.Extra, &sqlschema.IfExists{})
		}
	}

	plan, err := driver.PlanChanges(ctx, "plan", changes)
	if err != nil {
		return nil, err
	}

	// TODO: warning for falsy plan.Reversible?

	return plan.Changes, nil
}

var dbCreateCmd = &cobra.Command{
	Use:   "db:create",
	Short: "Creates database",
	RunE: func(cmd *cobra.Command, args []string) error {
		adapter, err := newAdapter()
		if err != nil {
			return err
		}
		return adapter.CreateDatabase()
	},
}

var dbDropCmd = &cobra.Command{
	Use:   "db:drop",
	Short: "Drops database",
	RunE: func(cmd *cobra.Command, args []string) error {
		adapter, err := newAdapter()
		if err != nil {
			return err
		}
		return adapter.DropDatabase()
	},
}

var dbMigrateCmd = &cobra.Command{
	Use:   "db:migrate",
	Short: "Migrates database according to schema",
	// TODO: support --version
	RunE: func(cmd *cobra.Command, args []string) error {
		a, err := newAdapter()
		if err != nil {
			return err
		}
		if err := a.Open(); err != nil {
			return err
		}

		migrator, err := adapter.NewMigrator(a)
		if err != nil {
			return err
		}

		if err := migrator.Up(); err != nil {
			return err
		}

		if err := dbMigrateGenerate(); err != nil {
			return err
		}
		if err := migrator.FindMigrations(); err != nil {
			return err
		}
		return migrator.Up()
	},
}

func createFile(f string, content string) error {
	// TODO: add logger
	// TODO: mkdirp
	fmt.Println("Created", f)
	if err := os.WriteFile(f, []byte(content), 0644); err != nil {
		return err
	}

	return nil
}

func dbMigrateGenerate() error {
	sch, err := newSchema()
	if err != nil {
		return err
	}

	adapter, err := newAdapter()
	if err != nil {
		return err
	}

	if err := adapter.Open(); err != nil {
		return err
	}

	database := sch.Databases[0]
	var changes []*migrate.Change
	if sch.Databases[0].Adapter == "postgresql" {
		changes, err = findSchemaChanges(adapter, database)
		if err != nil {
			return err
		}
	} else if sch.Databases[0].Adapter == "mysql" {
		changes, err = findMYSQLSchemaChanges(adapter, database)
		if err != nil {
			return err
		}
	} else if sch.Databases[0].Adapter == "sqlite" {
		changes, err = findSQLiteSchemaChanges(adapter, database)
		if err != nil {
			return err
		}
	}

	if len(changes) == 0 {
		fmt.Println("No schema changes found.")
		return nil
	}

	ups := []string{}
	downs := []string{}
	for _, change := range changes {
		if change.Cmd != "" {
			ups = append(ups, change.Cmd+";")
		}
		if change.Reverse != "" {
			downs = append([]string{change.Reverse + ";"}, downs...)
		}
	}
	// TODO: support windows line break
	up := strings.Join(ups, "\n")
	down := strings.Join(downs, "\n")

	// TODO: support migration name
	name := "auto"

	// TODO: should prompt in name conflict, support overwrite/skip flag like rails g
	var version string
	// for _, m := range migrator.UpMigrations {
	// 	if m.Name == name {
	// 		version = m.Version
	// 	}
	// }

	if version == "" {
		version = time.Now().Format("20060102150405")
	}

	if err := createFile(fmt.Sprintf("%s/migrations/%s_%s.up.sql", database.Name, version, name), up); err != nil {
		return err
	}
	if err := createFile(fmt.Sprintf("%s/migrations/%s_%s.down.sql", database.Name, version, name), down); err != nil {
		return err
	}
	return nil
}

var dbMigrateGenerateCmd = &cobra.Command{
	Use:   "db:migrate:generate",
	Short: "Generate versioned migration file according to schema",
	RunE: func(cmd *cobra.Command, args []string) error {
		return dbMigrateGenerate()
	},
}

var dbRollbackCmd = &cobra.Command{
	Use:   "db:rollback",
	Short: "Rollback last migration",
	// TODO: add --step
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

var dbMigrateStatusCmd = &cobra.Command{
	Use:   "db:migrate:status",
	Short: "List current migration status",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

var dbVersionCmd = &cobra.Command{
	Use:   "db:version",
	Short: "Prints current migration version",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}
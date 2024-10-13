package action

import (
	"bufio"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

func findSchemaChanges(a string, db adapter.Adapter, database *schema.Database) ([]*migrate.Change, error) {
	ctx := context.Background()

	var driver migrate.Driver
	var err error
	var from *sqlschema.Schema
	var desired *sqlschema.Schema

	environment := os.Getenv("QUERYX_ENV")
	if environment == "" {
		environment = "development"
	}
	config := database.LoadConfig(environment)
	dbName := adapter.NewConfig(config).Database

	if a == "postgresql" {
		driver, err = postgres.Open(db)
		if err != nil {
			return nil, err
		}
		from, err = driver.InspectSchema(ctx, "public", &sqlschema.InspectOptions{
			Tables: database.Tables(),
		})
		if err != nil {
			return nil, err
		}
		desired = database.CreatePostgreSQLSchema(dbName)
	} else if a == "mysql" {
		driver, err = mysql.Open(db)
		if err != nil {
			return nil, err
		}
		from, err = driver.InspectSchema(ctx, dbName, &sqlschema.InspectOptions{
			Tables: database.Tables(),
		})
		if err != nil {
			return nil, err
		}
		desired = database.CreateMySQLSchema(dbName)
	} else if a == "sqlite" {
		driver, err = sqlite.Open(db)
		if err != nil {
			return nil, err
		}
		from, err = driver.InspectSchema(ctx, "", &sqlschema.InspectOptions{
			Tables: database.Tables(),
		})
		if err != nil {
			return nil, err
		}
		desired = database.CreateSQLiteSchema("")
	}

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
		// defer a.Close()

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
	// defer adapter.Close()

	database := sch.Databases[0]
	changes, err := findSchemaChanges(sch.Databases[0].Adapter, adapter, database)
	if err != nil {
		return err
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

		stmts, err := change.ReverseStmts()
		if err != nil {
			return err
		}
		downs = append(stmts, downs...)
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

	upFile := filepath.Join(database.Name, "migrations", fmt.Sprintf("%s_%s.up.sql", version, name))
	downFile := filepath.Join(database.Name, "migrations", fmt.Sprintf("%s_%s.down.sql", version, name))

	if err := createFile(upFile, up); err != nil {
		return err
	}
	if err := createFile(downFile, down); err != nil {
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

func importCSV(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	if len(records) == 0 {
		return nil
	}

	columns := records[0]

	for _, record := range records[1:] {
		data := make(map[string]interface{})
		for i, col := range columns {
			data[col] = record[i]
		}

		fmt.Println("insert", columns, data)
	}

	// TODO: should generate insert statement based on adapter

	return nil
}

func importJSONL(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		var data map[string]interface{}
		if err := json.Unmarshal(scanner.Bytes(), &data); err != nil {
			return err
		}
		fmt.Println("insert", data)
	}

	return nil
}

var dbImportCmd = &cobra.Command{
	Use:   "db:import",
	Short: "Import data from a file to the database",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		replace, err := cmd.Flags().GetBool("replace")
		if err != nil {
			return err
		}
		table, err := cmd.Flags().GetString("table")
		if err != nil {
			return err
		}

		filePath := args[0]
		ext := filepath.Ext(filePath)
		switch ext {
		case ".csv":
			importCSV(filePath)
		case ".jsonl":
			importJSONL(filePath)
		}

		fmt.Println("replace", replace)
		fmt.Println("table", table)
		fmt.Println("filePath:", filePath)

		return nil
	},
}

var dbExportCmd = &cobra.Command{
	Use:   "db:export",
	Short: "Export data from the database to a file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		table, err := cmd.Flags().GetString("table")
		if err != nil {
			return err
		}

		filePath := args[0]
		fmt.Println("table", table)
		fmt.Println("filePath:", filePath)

		return nil
	},
}

func init() {
	dbImportCmd.Flags().Bool("replace", false, "Replace existing data")
	dbImportCmd.Flags().String("table", "", "Table to import data into")
	dbExportCmd.Flags().String("table", "", "Table to export data from")
}

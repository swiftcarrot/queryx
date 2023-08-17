package adapter

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Migrator struct {
	Adapter         Adapter
	MigrationsPath  string
	MigrationsTable string
	UpMigrations    UpMigrations
	DownMigrations  DownMigrations
}

func NewMigrator(adapter Adapter) (*Migrator, error) {
	// TODO: set from config
	migrationsPath := filepath.Join("db", "migrations")
	migrationsTable := "schema_migrations"
	os.MkdirAll(migrationsPath, 0766)

	m := &Migrator{
		Adapter:         adapter,
		MigrationsPath:  migrationsPath,
		MigrationsTable: migrationsTable,
	}

	err := m.FindMigrations()
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (m *Migrator) RunMigration(mg *Migration) error {
	fmt.Println("run", mg.Version, mg.Path)
	f, err := os.ReadFile(mg.Path)
	if err != nil {
		return err
	}
	sql := string(f)

	for _, line := range strings.Split(sql, "\n") {
		_, err = m.Adapter.ExecContext(context.Background(), line)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Migrator) FindMigrations() error {
	return filepath.Walk(m.MigrationsPath, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			migration, err := ParseMigrationFilename(info.Name())
			if err != nil {
				return err
			}
			migration.Path = p

			if migration.Direction == "up" {
				m.UpMigrations = append(m.UpMigrations, migration)
			} else if migration.Direction == "down" {
				m.DownMigrations = append(m.DownMigrations, migration)
			}
		}
		return nil
	})
}

func (m *Migrator) CreateTable() error {
	fmt.Println("create table", m.MigrationsTable)
	return nil
}

func (m *Migrator) Up() error {
	return m.UpWithVersion("")
}

func (m *Migrator) exists(ctx context.Context, version string) (bool, error) {
	exists := true
	rows, err := m.Adapter.QueryVersion(ctx, version)
	if err != nil {
		if err == sql.ErrNoRows {
			exists = false
		} else {
			return false, err
		}
	}
	if !rows.Next() {
		exists = false
	}
	return exists, nil
}

func (m *Migrator) UpWithVersion(version string) error {
	ctx := context.Background()

	if version != "" {
		fmt.Println("up with version", version)
	}

	err := m.Adapter.CreateMigrationsTable(ctx)
	if err != nil {
		return err
	}

	sort.Sort(m.UpMigrations)
	for _, um := range m.UpMigrations {
		exists, err := m.exists(ctx, um.Version)
		if err != nil {
			return err
		}

		if !exists {
			if err := m.RunMigration(um); err != nil {
				return err
			}
			_, err := m.Adapter.ExecContext(ctx, fmt.Sprintf("insert into %s (version) values ('%s')", "schema_migrations", um.Version))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Migrator) Down() error {
	return m.DownWithStep(1)
}

func (m *Migrator) DownWithStep(step int) error {
	fmt.Println("DownWithStep", step)
	ctx := context.Background()

	sort.Sort(m.DownMigrations)

	if step > 0 && len(m.DownMigrations) >= step {
		m.DownMigrations = m.DownMigrations[:step]
	}

	for _, dm := range m.DownMigrations {
		exists, err := m.exists(ctx, dm.Version)
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("migration version %s does not exist", dm.Version)
		}

		if err := m.RunMigration(dm); err != nil {
			return err
		}
		_, err = m.Adapter.ExecContext(ctx, fmt.Sprintf("delete from %s where version = '%s'", "schema_migrations", dm.Version))
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Migrator) Reset() error {
	return nil
}

func (m *Migrator) Status(out io.Writer) error {
	ctx := context.Background()

	rows, err := m.Adapter.QueryContext(ctx, "select * from schema_migrations")
	if err != nil {
		return err
	}

	defer rows.Close()
	if err := rows.Err(); err != nil {
		return err
	}

	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return err
		}
		log.Println(version)
	}

	if err := rows.Close(); err != nil {
		return fmt.Errorf("closing rows %w", err)
	}

	log.Println(rows)

	return nil
}

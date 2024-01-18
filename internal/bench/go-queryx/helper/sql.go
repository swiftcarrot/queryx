package helper

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

func CreateTables() error {
	queries := [][]string{
		{
			`DROP TABLE IF EXISTS models;`,
			`CREATE TABLE models (
			id SERIAL NOT NULL,
			name text NOT NULL,
			title text NOT NULL,
			fax text NOT NULL,
			web text NOT NULL,
			age integer NOT NULL,
			"righ" boolean NOT NULL,
			counter bigint NOT NULL,
			CONSTRAINT models_pkey PRIMARY KEY (id)
			) WITH (OIDS=FALSE);`,
		},
		{
			`DROP TABLE IF EXISTS model5;`,
			`CREATE TABLE model5 (
			id SERIAL NOT NULL,
			name text NOT NULL,
			title text NOT NULL,
			fax text NOT NULL,
			web text NOT NULL,
			age integer NOT NULL,
			"righ" boolean NOT NULL,
			counter bigint NOT NULL
			) WITH (OIDS=FALSE);`,
		},
	}

	db, err := sql.Open("postgres", OrmSource)
	if err != nil {
		return fmt.Errorf("init_tables: %w", err)
	}

	defer func() {
		_ = db.Close()
	}()

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("init_tables: %w", err)
	}

	for _, query := range queries {
		for _, line := range query {
			_, err = db.Exec(line)
			if err != nil {
				return fmt.Errorf("init_tables: %w", err)
			}
		}
	}

	return nil
}

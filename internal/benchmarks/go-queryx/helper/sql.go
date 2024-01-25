package helper

import (
	"github.com/swiftcarrot/queryx/internal/benchmarks/go-queryx/db"
	"log"
)

func CreateTables() error {
	client, err := db.NewClient()
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, err = client.Exec(`DROP TABLE IF EXISTS models;`)
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, err = client.Exec(
		`CREATE TABLE models (
			id SERIAL  PRIMARY KEY, 
			name text NOT NULL,
			title text NOT NULL,
			fax text NOT NULL,
			web text NOT NULL,
			age integer NOT NULL,
			righ boolean NOT NULL,
			counter bigint NOT NULL
			);`)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return err
}

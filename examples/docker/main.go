package main

import (
	"fmt"
	"log"

	"docker/db"
)

func main() {
	c, err := db.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	type Row struct {
		Version string `db:"version"`
	}
	var row Row
	err = c.QueryOne("select sqlite_version() as version").Scan(&row)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(row.Version)
}

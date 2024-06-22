package helper

import (
	"fmt"
	"github.com/swiftcarrot/queryx/internal/benchmarks/go-queryx/db"
	"log"
)

const (
	sqlSqliteCreateTable = `CREATE TABLE models (
    id integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    name varchar NULL,
    title varchar NULL,
    fax varchar NULL,
    web varchar NULL,
    age bigint NULL,
    righ boolean NULL,
    counter integer NULL
                    );`
	sqlPostgresCreateTable = `CREATE TABLE "public"."models" (
    "id" bigserial NOT NULL,
    "name" character varying NULL,
    "title" character varying NULL, 
    "fax" character varying NULL,
    "web" character varying NULL,
    "age" bigint NULL,
    "righ" boolean NULL,
    "counter" integer NULL,
    PRIMARY KEY ("id"));`
	sqlMysqlCreateTable = `CREATE TABLE test.models (
    id bigint NOT NULL AUTO_INCREMENT,
    name varchar(255) NULL, title varchar(255) NULL,
    fax varchar(255) NULL, web varchar(255) NULL,
    age bigint NULL,
    righ bool NULL,
    counter int NULL,
    PRIMARY KEY (id));`
)

func CreateTables(client *db.QXClient, adapter string) error {
	var sql string
	switch adapter {
	case "mysql":
		sql = sqlMysqlCreateTable
	case "postgresql":
		sql = sqlPostgresCreateTable
	case "sqlite":
		sql = sqlSqliteCreateTable
	default:
		return fmt.Errorf("this type of adapter is not supported:%v", adapter)
	}
	_, err := client.Exec(`DROP TABLE IF EXISTS models;`)
	if err != nil {
		return err
	}
	_, err = client.Exec(sql)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return err
}

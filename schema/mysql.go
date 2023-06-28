package schema

import (
	"fmt"
	"strings"

	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/schema"
)

// convert a mysql queryx database schema to an atlas sql schema
func (d *Database) CreateMySQLSchema(dbName string) *schema.Schema {
	public := schema.New(dbName)

	for _, model := range d.Models {
		t := schema.NewTable(model.TableName)
		columnMap := map[string]*schema.Column{}

		for _, c := range model.Columns {
			col := schema.NewColumn(c.Name)
			switch c.Type {
			case "bigint":
				col.SetType(&schema.IntegerType{T: mysql.TypeBigInt})
				if c.AutoIncrement {
					col.AddAttrs(&mysql.AutoIncrement{})
				}
			case "string":
				col.SetType(&schema.StringType{T: mysql.TypeVarchar, Size: 256})
			case "text":
				col.SetType(&schema.StringType{T: "text", Size: 0})
			case "integer":
				col.SetType(&schema.IntegerType{T: "integer"})
			case "float":
				col.SetType(&schema.FloatType{T: mysql.TypeFloat})
			case "boolean":
				col.SetType(&schema.BoolType{T: mysql.TypeBoolean})
			case "enum":
				col.SetType(&schema.StringType{T: mysql.TypeEnum, Size: 0})
			case "date":
				col.SetType(&schema.TimeType{T: mysql.TypeDate})
			case "time":
				col.SetType(&schema.TimeType{T: mysql.TypeDateTime})
			case "datetime":
				col.SetType(&schema.TimeType{T: mysql.TypeDateTime})
			case "json", "jsonb":
				col.SetType(&schema.JSONType{T: mysql.TypeJSON})
			case "uuid":
				col.SetType(&schema.StringType{T: mysql.TypeVarchar, Size: 256})
			}

			col.SetNull(c.Null)
			t.AddColumns(col)
			columnMap[c.Name] = col
		}

		cols := []*schema.Column{}
		for _, name := range model.PrimaryKey.ColumnNames {
			cols = append(cols, columnMap[name])
		}
		pk := schema.NewPrimaryKey(cols...)
		t.SetPrimaryKey(pk)

		for _, i := range model.Index {
			name := fmt.Sprintf("%s_%s_index", i.TableName, strings.Join(i.ColumnNames, "_"))
			columns := []*schema.Column{}
			for _, name := range i.ColumnNames {
				columns = append(columns, columnMap[name])
			}
			var index *schema.Index
			if i.Unique {
				index = schema.NewUniqueIndex(name).AddColumns(columns...)
			} else {
				index = schema.NewIndex(name).AddColumns(columns...)
			}
			t.AddIndexes(index)
		}

		public.AddTables(t)
	}

	return public
}

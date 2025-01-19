package schema

import (
	"fmt"
	"strconv"
	"strings"

	"ariga.io/atlas/sql/schema"
	"ariga.io/atlas/sql/sqlite"
)

func (d *Database) CreateSQLiteSchema(dbName string) *schema.Schema {
	public := schema.New(dbName)

	for _, model := range d.Models {
		columnMap := map[string]*schema.Column{}

		t := schema.NewTable(model.TableName)
		if t.Name == "sqlite_master" {
			continue
		}
		for _, c := range model.Columns {
			col := schema.NewColumn(c.Name)

			switch c.Type {
			case "string", "uuid":
				col.SetType(&schema.StringType{T: "varchar", Size: 0})
				if c.Default != nil {
					d, ok := c.Default.(string)
					if ok {
						col.SetDefault(&schema.RawExpr{X: fmt.Sprintf("'%s'", d)})
					}
				}
			case "text":
				col.SetType(&schema.StringType{T: "text", Size: 0})
			case "integer":
				col.SetType(&schema.IntegerType{T: "integer"})
				if c.Default != nil {
					d, ok := c.Default.(int)
					if ok {
						col.SetDefault(&schema.RawExpr{X: strconv.Itoa(d)})
					}
				}
			case "bigint":
				if c.AutoIncrement {
					col.SetType(&schema.IntegerType{T: "integer"})
					col.AddAttrs(&sqlite.AutoIncrement{})
				} else {
					col.SetType(&schema.IntegerType{T: "bigint"})
					d, ok := c.Default.(int)
					if ok {
						col.SetDefault(&schema.RawExpr{X: strconv.Itoa(d)})
					}
				}
			case "float":
				col.SetType(&schema.FloatType{T: "float"})
				if c.Default != nil {
					d, ok := c.Default.(float64)
					if ok {
						col.SetDefault(&schema.RawExpr{X: strconv.FormatFloat(d, 'f', 10, 64)})
					}
				}
			case "boolean":
				col.SetType(&schema.BoolType{T: "boolean"})
				if c.Default != nil {
					d, ok := c.Default.(bool)
					if ok {
						col.SetDefault(&schema.RawExpr{X: strconv.FormatBool(d)})
					}
				}
			case "enum":
				col.SetType(&schema.StringType{T: "enum", Size: 0})
			case "date":
				col.SetType(&schema.TimeType{T: "date"})
			case "time":
				col.SetType(&schema.TimeType{T: "datetime"}) // TODO: "time"
			case "datetime":
				col.SetType(&schema.TimeType{T: "datetime"})
			case "json", "jsonb":
				col.SetType(&schema.JSONType{T: "json"})
			default:
				fmt.Printf("This type is not supported:%+v", col.Type)
			}

			if c.Default == nil {
				col.SetNull(c.Null)
			}
			t.AddColumns(col)
			columnMap[c.Name] = col
		}

		if model.PrimaryKey != nil {
			cols := []*schema.Column{}
			for _, name := range model.PrimaryKey.ColumnNames {
				cols = append(cols, columnMap[name])
			}
			pk := schema.NewPrimaryKey(cols...)
			t.SetPrimaryKey(pk)
		}

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
		public.Name = "main"
	}

	return public
}

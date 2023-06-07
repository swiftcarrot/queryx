package schema

import (
	"fmt"
	"strconv"
	"strings"

	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/schema"
)

func (d *Database) CreateMYSQLSchema(dbName string) *schema.Schema {
	public := schema.New(dbName)
	for _, model := range d.Models {
		columnMap := map[string]*schema.Column{}

		t := schema.NewTable(model.TableName)

		for _, c := range model.Columns {
			col := schema.NewColumn(c.Name)
			if c.Type == "pk" {
				col.SetType(&mysql.BitType{T: "int auto_increment"}).SetNull(false)
				pk := schema.NewPrimaryKey(col)
				t.SetPrimaryKey(pk)
				continue
			}

			switch c.Type {
			case "id":
				col.SetType(&schema.IntegerType{T: mysql.TypeInt})
			case "string":
				col.SetType(&schema.StringType{T: "varchar(256)", Size: 0})
				if c.Default != nil {
					d, ok := c.Default.(string)
					if ok {
						col.SetType(&schema.StringType{T: "varchar(256)", Size: 0}).SetDefault(&schema.RawExpr{X: fmt.Sprintf("'%s'", d)})
					}
				}
			case "text":
				col.SetType(&schema.StringType{T: "text", Size: 0})
			case "integer":
				col.SetType(&schema.IntegerType{T: "integer"})
				if c.Default != nil {
					d, ok := c.Default.(int)
					if ok {
						col.SetType(&schema.IntegerType{T: "integer"}).SetDefault(&schema.RawExpr{X: strconv.Itoa(d)})
					}
				}
			case "float":
				col.SetType(&schema.FloatType{T: mysql.TypeFloat})
				if c.Default != nil {
					d, ok := c.Default.(float64)
					if ok {
						col.SetType(&schema.FloatType{T: mysql.TypeFloat}).SetDefault(&schema.RawExpr{X: strconv.FormatFloat(d, 'f', 10, 64)})
					}
				}
			case "boolean":
				col.SetType(&schema.BoolType{T: mysql.TypeBoolean})
				if c.Default != nil {
					d, ok := c.Default.(bool)
					if ok {
						col.SetType(&schema.BoolType{T: mysql.TypeBoolean}).SetDefault(&schema.RawExpr{X: strconv.FormatBool(d)})
					}
				}
			case "enum":
				col.SetType(&schema.StringType{T: mysql.TypeEnum, Size: 0})
			case "date":
				col.SetType(&schema.TimeType{T: mysql.TypeDate})
			case "time":
				col.SetType(&schema.TimeType{T: mysql.TypeDateTime})
			case "datetime":
				col.SetType(&schema.TimeType{T: mysql.TypeDateTime})
			case "jsonb":
				col.SetType(&schema.JSONType{T: mysql.TypeJSON})
			}

			col.SetNull(c.Null)
			t.AddColumns(col)
			columnMap[c.Name] = col
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
	}

	return public
}

package schema

import (
	"fmt"
	"strconv"
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
				} else {
					if c.Default != nil {
						d, ok := c.Default.(int)
						if ok {
							col.SetType(&schema.IntegerType{T: mysql.TypeBigInt}).SetDefault(&schema.RawExpr{X: strconv.Itoa(d)})
						}
					}
				}
			case "string":
				col.SetType(&schema.StringType{T: mysql.TypeVarchar, Size: 255})
				if c.Default != nil {
					d, ok := c.Default.(string)
					if ok {
						col.SetDefault(&schema.RawExpr{X: d})
					}
				}
			case "text":
				col.SetType(&schema.StringType{T: mysql.TypeText})
				if c.Default != nil {
					d, ok := c.Default.(string)
					if ok {
						col.SetDefault(&schema.RawExpr{X: d})
					}
				}
			case "integer":
				col.SetType(&schema.IntegerType{T: mysql.TypeInt})
				if c.Default != nil {
					d, ok := c.Default.(int)
					if ok {
						col.SetDefault(&schema.RawExpr{X: strconv.Itoa(d)})
					}
				}
			case "float":
				col.SetType(&schema.FloatType{T: mysql.TypeFloat})
				if c.Default != nil {
					f, ok := c.Default.(float64)
					if ok {
						str := strconv.FormatFloat(f, 'E', -10, 64)
						col.SetDefault(&schema.RawExpr{X: str})
					}
				}
			case "boolean":
				col.SetType(&schema.BoolType{T: mysql.TypeBoolean})
				if c.Default != nil {
					b, ok := c.Default.(bool)
					if ok {
						str := strconv.FormatBool(b)
						col.SetDefault(&schema.RawExpr{X: str})
					}
				}
			case "enum":
				col.SetType(&schema.StringType{T: mysql.TypeEnum, Size: 0})
			case "date":
				col.SetType(&schema.TimeType{T: mysql.TypeDate})
			case "time":
				col.SetType(&schema.TimeType{T: mysql.TypeDateTime})
			case "datetime":
				i := 6
				col.SetType(&schema.TimeType{T: mysql.TypeDateTime, Precision: &i})
			case "json", "jsonb":
				col.SetType(&schema.JSONType{T: mysql.TypeJSON})
			case "uuid":
				col.SetType(&schema.StringType{T: mysql.TypeVarchar, Size: 36})
				if c.Default != nil {
					d, ok := c.Default.(string)
					if ok {
						col.SetDefault(&schema.RawExpr{X: d})
					}
				}
			}

			col.SetNull(c.Null)
			t.AddColumns(col)
			columnMap[c.Name] = col
		}

		cols := []*schema.Column{}
		if model.PrimaryKey != nil {
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
	}

	return public
}

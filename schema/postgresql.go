package schema

import (
	"fmt"
	"strconv"
	"strings"

	"ariga.io/atlas/sql/postgres"
	"ariga.io/atlas/sql/schema"
)

// convert a postgresql queryx database schema to an atlas sql schema
func (d *Database) CreatePostgreSQLSchema(dbName string) *schema.Schema {
	public := schema.New("public")

	for _, model := range d.Models {
		t := schema.NewTable(model.TableName)
		columnMap := map[string]*schema.Column{}

		for _, c := range model.Columns {
			col := schema.NewColumn(c.Name)

			switch c.Type {
			case "bigint":
				if c.AutoIncrement {
					col.SetType(&postgres.SerialType{T: postgres.TypeBigSerial})
				} else {
					col.SetType(&schema.IntegerType{T: postgres.TypeInt8})
				}
			case "string":
				if c.Array {
					col.SetType(&postgres.ArrayType{Type: &schema.StringType{T: "character varying", Size: 0}, T: "varchar[]"})
					if c.Default != nil {
						d, ok := c.Default.(string)
						if ok {
							col.SetType(&postgres.ArrayType{Type: &schema.StringType{T: "character varying", Size: 0}, T: "varchar[]"}).SetDefault(&schema.RawExpr{X: fmt.Sprintf("'%s'", d)})
						}
					}
				} else {
					col.SetType(&schema.StringType{T: "character varying", Size: 0})
					if c.Default != nil {
						d, ok := c.Default.(string)
						if ok {
							col.SetType(&schema.StringType{T: "character varying", Size: 0}).SetDefault(&schema.RawExpr{X: fmt.Sprintf("'%s'", d)})
						}
					}
				}
			case "text":
				col.SetType(&schema.StringType{T: "character varying", Size: 0})
			case "integer":
				if c.Array {
					col.SetType(&postgres.ArrayType{Type: &schema.IntegerType{T: "integer"}, T: "int[]"})
					if c.Default != nil {
						d, ok := c.Default.(int)
						if ok {
							col.SetType(&postgres.ArrayType{Type: &schema.IntegerType{T: "integer"}, T: "int[]"}).SetDefault(&schema.RawExpr{X: strconv.Itoa(d)})
						}
					}
				} else {
					col.SetType(&schema.IntegerType{T: "integer"})
					if c.Default != nil {
						d, ok := c.Default.(int)
						if ok {
							col.SetType(&schema.IntegerType{T: "integer"}).SetDefault(&schema.RawExpr{X: strconv.Itoa(d)})
						}
					}
				}
			case "float":
				col.SetType(&schema.FloatType{T: postgres.TypeFloat8})
				if c.Default != nil {
					d, ok := c.Default.(float64)
					if ok {
						col.SetType(&schema.FloatType{T: postgres.TypeFloat8}).SetDefault(&schema.RawExpr{X: strconv.FormatFloat(d, 'f', 10, 64)})
					}
				}
			case "boolean":
				col.SetType(&schema.BoolType{T: postgres.TypeBoolean})
				if c.Default != nil {
					d, ok := c.Default.(bool)
					if ok {
						col.SetType(&schema.BoolType{T: postgres.TypeBoolean}).SetDefault(&schema.RawExpr{X: strconv.FormatBool(d)})
					}
				}
			case "enum":
				col.SetType(&schema.StringType{T: "character varying", Size: 0})
			case "date":
				col.SetType(&schema.TimeType{T: postgres.TypeDate})
			case "time":
				col.SetType(&schema.TimeType{T: postgres.TypeTime})
			case "datetime":
				col.SetType(&schema.TimeType{T: postgres.TypeTimestamp})
			case "json":
				col.SetType(&schema.JSONType{T: postgres.TypeJSON})
			case "jsonb":
				col.SetType(&schema.JSONType{T: postgres.TypeJSONB})
			case "uuid":
				col.SetType(&postgres.UUIDType{T: postgres.TypeUUID})
				if c.Default != nil {
					d, ok := c.Default.(string)
					if ok {
						col.SetType(&postgres.UUIDType{T: postgres.TypeUUID}).SetDefault(&schema.RawExpr{X: d})
					}
				}
			}

			col.SetNull(c.Null)
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
	}

	return public
}

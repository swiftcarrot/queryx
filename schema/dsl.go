package schema

import (
	"fmt"

	"github.com/swiftcarrot/queryx/inflect"
)

func NewSchema() *Schema {
	return &Schema{
		Databases: make([]*Database, 0),
	}
}

func (s *Schema) NewDatabase(name string) *Database {
	return &Database{
		Name:       name,
		Configs:    make([]*Config, 0),
		Generators: make([]*Generator, 0),
		Models:     make([]*Model, 0),
	}
}

func (d *Database) NewConfig(environment string) *Config {
	return &Config{
		Adapter:     d.Adapter,
		Environment: environment,
	}
}

// initialize new model with defaults
func (d *Database) NewModel(name string) *Model {
	m := &Model{
		Name:       name,
		TableName:  inflect.Tableize(name),
		Columns:    make([]*Column, 0),
		Attributes: make([]*Attribute, 0),
	}

	d.Models = append(d.Models, m)

	return m
}

// set default primary key with id column for a model
func (m *Model) AddDefaultPrimaryKey() {
	col := &Column{
		Name:          "id",
		Type:          "bigint",
		AutoIncrement: true,
		Null:          false,
	}
	m.DefaultPrimaryKey = true
	m.Columns = append(m.Columns, col)
	m.PrimaryKey = &PrimaryKey{
		Columns:     []*Column{col},
		ColumnNames: []string{"id"},
	}
}

// custom primary key for a model
func (m *Model) SetPrimaryKey(columnNames []string) {
	cols := []*Column{}
	cmap := make(map[string]*Column)
	for _, c := range m.Columns {
		cmap[c.Name] = c
	}
	for _, n := range columnNames {
		// TODO: check column
		col := cmap[n]
		col.Null = false // allow null on primary key fields will cause an error
		cols = append(cols, col)
	}
	m.DefaultPrimaryKey = false
	m.PrimaryKey = &PrimaryKey{
		Columns:     cols,
		ColumnNames: columnNames,
	}
}

func (m *Model) AddTimestamps() {
	m.Timestamps = true
	m.Columns = append(m.Columns, &Column{
		Name: "created_at",
		Type: "datetime",
		Null: false,
	}, &Column{
		Name: "updated_at",
		Type: "datetime",
		Null: false,
	})
}

func (m *Model) NewColumn(name string, colType string) *Column {
	return &Column{
		Name: name,
		Type: colType,
		Null: true,
	}
}

func (m *Model) AddHasMany(hasMany *HasMany) {
	m.HasMany = append(m.HasMany, hasMany)
}

func (m *Model) AddHasOne(hasOne *HasOne) {
	m.HasOne = append(m.HasOne, hasOne)
}

func (m *Model) AddBelongsTo(belongsTo *BelongsTo) {
	m.BelongsTo = append(m.BelongsTo, belongsTo)

	// TODO: support foreign key, not null
	col := &Column{
		Name: fmt.Sprintf("%s_id", belongsTo.Name),
		Type: "bigint",
		Null: true,
	}
	m.Columns = append(m.Columns, col)
}

func (m *Model) AddIndex(index *Index) {
	index.TableName = m.TableName
	m.Index = append(m.Index, index)
}

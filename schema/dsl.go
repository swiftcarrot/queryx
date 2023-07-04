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
		TimeZone:   "Local",
		Configs:    make([]*Config, 0),
		Generators: make([]*Generator, 0),
		Models:     make([]*Model, 0),
		models:     make(map[string]*Model),
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
		Database:   d,
		Name:       name,
		TableName:  inflect.Tableize(name),
		Columns:    make([]*Column, 0),
		Attributes: make([]*Attribute, 0),
	}

	d.Models = append(d.Models, m)
	d.models[m.Name] = m

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
	if hasMany.ModelName == "" {
		hasMany.ModelName = inflect.Pascal(inflect.Singular(hasMany.Name))
	}
	if hasMany.ForeignKey == "" {
		hasMany.ForeignKey = fmt.Sprintf("%s_id", inflect.Snake(m.Name))
	}
	m.HasMany = append(m.HasMany, hasMany)
	m.Database.buildAssociation()
}

func (m *Model) AddHasOne(hasOne *HasOne) {
	if hasOne.ModelName == "" {
		hasOne.ModelName = inflect.Pascal(inflect.Singular(hasOne.Name))
	}
	if hasOne.ForeignKey == "" {
		hasOne.ForeignKey = fmt.Sprintf("%s_id", inflect.Snake(m.Name))
	}
	m.HasOne = append(m.HasOne, hasOne)
	m.Database.buildAssociation()
}

func NewBelongsTo(name string) *BelongsTo {
	belongsTo := &BelongsTo{
		Name:  name,
		Type:  "bigint", // TODO: support other types such as integer?
		Null:  true,
		Index: false,
	}
	return belongsTo
}

func (m *Model) AddBelongsTo(belongsTo *BelongsTo) {
	if belongsTo.ModelName == "" {
		belongsTo.ModelName = inflect.Pascal(inflect.Singular(belongsTo.Name))
	}
	if belongsTo.ForeignKey == "" {
		belongsTo.ForeignKey = fmt.Sprintf("%s_id", belongsTo.Name)
	}

	m.BelongsTo = append(m.BelongsTo, belongsTo)

	col := &Column{
		Name: belongsTo.ForeignKey,
		Type: belongsTo.Type,
		Null: belongsTo.Null,
	}
	m.Columns = append(m.Columns, col)

	if belongsTo.Index {
		m.AddIndex(&Index{
			ColumnNames: []string{belongsTo.ForeignKey},
		})
	}

	m.Database.buildAssociation()
}

func (m *Model) AddIndex(index *Index) {
	index.TableName = m.TableName
	m.Index = append(m.Index, index)
}

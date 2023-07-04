package schema

import (
	"strings"
)

type Model struct {
	Database *Database
	// model name in camel case
	Name       string
	TableName  string
	TimeZone   string
	Timestamps bool
	Columns    []*Column
	Attributes []*Attribute
	BelongsTo  []*BelongsTo
	HasMany    []*HasMany
	HasOne     []*HasOne
	Index      []*Index
	// whether to automatically add a default primary key column named "id", default true
	DefaultPrimaryKey bool
	PrimaryKey        *PrimaryKey
}

// TODO: pretty print a model
func (m *Model) String() string {
	var b strings.Builder
	b.WriteString(m.Name)
	b.WriteString("(columns=")
	for _, c := range m.Columns {
		b.WriteString(c.Name)
		b.WriteString(",")
	}
	b.WriteString("primary_key=")
	b.WriteString(strings.Join(m.PrimaryKey.ColumnNames, ", "))
	b.WriteString(")")
	return b.String()
}

type PrimaryKey struct {
	TableName   string
	Columns     []*Column
	ColumnNames []string
}

type Column struct {
	// column name in snake case
	Name string
	Type string
	// array type
	Array bool
	// nullable, allow_null by default
	Null bool
	// sql auto_increment
	AutoIncrement bool
	Default       interface{} // TODO: support default
}

type Type struct {
	Name string
}

type Enum struct {
	Name string
}

type Attribute struct {
	Name string
	Type string
	Null bool
}

type HasMany struct {
	Name       string
	ModelName  string
	Through    string
	ForeignKey string
	Source     string
	BelongsTo  *BelongsTo
}

type HasOne struct {
	Name       string
	ModelName  string
	Through    string
	ForeignKey string
	BelongsTo  *BelongsTo
}

type BelongsTo struct {
	Name        string
	ModelName   string
	ForeignKey  string
	ForeignType string
	PrimaryKey  string
	Type        string
	Index       bool
	Null        bool
	Dependent   string
	Optional    bool
	Required    bool
	Default     bool
}

type Index struct {
	TableName   string
	ColumnNames []string
	Unique      bool
}

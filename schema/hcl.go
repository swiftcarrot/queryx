package schema

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/swiftcarrot/queryx/inflect"
	"github.com/swiftcarrot/queryx/types"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/gocty"
)

var hclSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "database", LabelNames: []string{"name"}},
	},
}

var hclClient = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "adapter"},
		{Name: "time_zone", Required: false},
	},
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "config", LabelNames: []string{"name"}},
		{Type: "generator", LabelNames: []string{"name"}},
		{Type: "model", LabelNames: []string{"name"}},
	},
}

var hclConfig = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "url", Required: false},
		{Name: "database", Required: false},
		{Name: "host", Required: false},
		{Name: "port", Required: false},
		{Name: "user", Required: false},
		{Name: "password", Required: false},
		{Name: "encoding", Required: false},
		{Name: "pool", Required: false},
		{Name: "timeout", Required: false},
		{Name: "socket", Required: false},
		{Name: "raw_options", Required: false},
	},
}

var hclModel = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "table_name"},
		{Name: "primary_key"},
		{Name: "timestamps"},
		{Name: "default_primary_key"},
	},
	Blocks: []hcl.BlockHeaderSchema{
		{Type: "column", LabelNames: []string{"name"}},
		{Type: "attribute", LabelNames: []string{"name"}},
		{Type: "belongs_to", LabelNames: []string{"name"}},
		{Type: "has_one", LabelNames: []string{"name"}},
		{Type: "has_many", LabelNames: []string{"name"}},
		{Type: "index", LabelNames: []string{}},
		{Type: "primary_key", LabelNames: []string{}},
	},
}

var hclColumn = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "type"},
		{Name: "array"},
		{Name: "null"},
		{Name: "default"},
		{Name: "unique"},
	},
	Blocks: []hcl.BlockHeaderSchema{},
}

var hclAttribute = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "type"},
		{Name: "null"},
		{Name: "default"},
	},
	Blocks: []hcl.BlockHeaderSchema{},
}

var hclPrimaryKey = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "columns", Required: true},
	},
	Blocks: []hcl.BlockHeaderSchema{},
}

var hclIndex = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "columns", Required: true},
		{Name: "unique"},
	},
	Blocks: []hcl.BlockHeaderSchema{},
}

var hclHasMany = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "class_name"},
		{Name: "through"},
		{Name: "foreign_key"},
		{Name: "source"},
	},
}

var hclHasOne = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "class_name"},
		{Name: "through"},
		{Name: "foreign_key"},
	},
}

var hclBelongsTo = &hcl.BodySchema{
	Attributes: []hcl.AttributeSchema{
		{Name: "class_name"},
		{Name: "foreign_key"},
	},
}

func (s *Schema) databaseFromBlock(block *hcl.Block, ctx *hcl.EvalContext) (*Database, error) {
	name := block.Labels[0]
	database := s.NewDatabase(name)

	content, d := block.Body.Content(hclClient)
	if d.HasErrors() {
		return nil, d.Errs()[0]
	}

	for name, attr := range content.Attributes {
		value, d := attr.Expr.Value(ctx)
		if d.HasErrors() {
			return nil, d.Errs()[0]
		}

		switch name {
		case "adapter":
			database.Adapter = value.AsString()
		case "time_zone":
			database.TimeZone = value.AsString()
		}
	}

	blocks := content.Blocks.ByType()
	for blockName := range blocks {
		switch blockName {
		case "config":
			for _, b := range blocks[blockName] {
				config, err := database.configFromBlock(b, ctx)
				if err != nil {
					return nil, err
				}
				database.Configs = append(database.Configs, config)
			}
		case "generator":
			for _, b := range blocks[blockName] {
				generator, err := generatorFromBlock(b, ctx)
				if err != nil {
					return nil, err
				}
				database.Generators = append(database.Generators, generator)
			}
		case "model":
			for _, b := range blocks[blockName] {
				_, err := database.modelFromBlock(b, ctx)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return database, nil
}

func (db *Database) configFromBlock(block *hcl.Block, ctx *hcl.EvalContext) (*Config, error) {
	environment := block.Labels[0]
	config := db.NewConfig(environment)

	content, d := block.Body.Content(hclConfig)
	if d.HasErrors() {
		return nil, d.Errs()[0]
	}

	for name, attr := range content.Attributes {
		value, d := attr.Expr.Value(ctx)
		if d.HasErrors() {
			return nil, d.Errs()[0]
		}

		switch name {
		case "url":
			// TODO:
			var v envKey
			err := gocty.FromCtyValue(value, &v)
			if err != nil {
				config.URL = types.StringOrEnv{Value: value.AsString()}
			} else {
				config.URL = types.StringOrEnv{EnvKey: v.Key}
			}
		case "database":
			config.Database = value.AsString()
		case "host":
			config.Host = value.AsString()
		case "port":
			config.Port = value.AsString()
		case "user":
			config.User = value.AsString()
		case "password":
			config.Password = value.AsString()
		case "encoding":
			config.Encoding = value.AsString()
		// case "pool":
		// 	config.Pool = value.AsString()
		// case "timeout":
		// 	config.Timeout = value.AsString()
		case "socket":
			config.Socket = value.AsString()
		case "raw_options":
			config.RawOptions = value.AsString()
		}
	}

	return config, nil
}

func generatorFromBlock(block *hcl.Block, ctx *hcl.EvalContext) (*Generator, error) {
	name := block.Labels[0]

	return &Generator{
		Name: name,
	}, nil
}

func (db *Database) modelFromBlock(block *hcl.Block, ctx *hcl.EvalContext) (*Model, error) {
	name := block.Labels[0]
	m := db.NewModel(name)

	content, d := block.Body.Content(hclModel)
	if d.HasErrors() {
		return nil, d.Errs()[0]
	}

	addDefaultPrimaryKey := true
	addTimestamps := true

	for name, attr := range content.Attributes {
		value, d := attr.Expr.Value(ctx)
		if d.HasErrors() {
			return nil, d.Errs()[0]
		}

		switch name {
		case "table_name":
			// overwrite default table name if set
			m.TableName = value.AsString()
		case "default_primary_key":
			defaultPrimaryKey := valueAsBool(value)
			if !defaultPrimaryKey {
				addDefaultPrimaryKey = false
			}
		case "timestamps":
			timestamps := valueAsBool(value)
			if !timestamps {
				addTimestamps = false
			}
		}
	}

	if addDefaultPrimaryKey {
		m.AddDefaultPrimaryKey()
	}

	for _, block := range content.Blocks {
		switch block.Type {
		case "attribute":
			attr, err := attributeFromBlock(block, ctx)
			if err != nil {
				return nil, err
			}
			m.Attributes = append(m.Attributes, attr)
		case "column":
			col, err := columnFromBlock(block, ctx)
			if err != nil {
				return nil, err
			}
			m.Columns = append(m.Columns, col)
		case "has_many":
			hasMany, err := hasManyFromBlock(block, ctx)
			if err != nil {
				return nil, err
			}
			m.AddHasMany(hasMany)
		case "has_one":
			hasOne, err := hasOneFromBlock(block, ctx)
			if err != nil {
				return nil, err
			}
			m.AddHasOne(hasOne)
		case "belongs_to":
			belongsTo, err := belongsToFromBlock(block, ctx)
			if err != nil {
				return nil, err
			}
			m.AddBelongsTo(belongsTo)
		case "index":
			index, err := indexFromBlock(block, ctx)
			if err != nil {
				return nil, err
			}
			m.AddIndex(index)
		case "primary_key":
			primaryKey, err := primaryKeyFromBlock(block, ctx)
			if err != nil {
				return nil, err
			}
			m.SetPrimaryKey(primaryKey.ColumnNames)
		}
	}

	if addTimestamps {
		m.AddTimestamps()
	}

	return m, nil
}

func indexFromBlock(block *hcl.Block, ctx *hcl.EvalContext) (*Index, error) {
	index := &Index{}

	content, d := block.Body.Content(hclIndex)
	if d.HasErrors() {
		return nil, d.Errs()[0]
	}

	for name, attr := range content.Attributes {
		value, d := attr.Expr.Value(ctx)
		if d.HasErrors() {
			return nil, d.Errs()[0]
		}

		switch name {
		case "columns":
			index.ColumnNames = valueAsStringSlice(value)
		case "unique":
			index.Unique = valueAsBool(value)
		}
	}

	return index, nil
}

func primaryKeyFromBlock(block *hcl.Block, ctx *hcl.EvalContext) (*PrimaryKey, error) {
	primaryKey := &PrimaryKey{}

	content, d := block.Body.Content(hclPrimaryKey)
	if d.HasErrors() {
		return nil, d.Errs()[0]
	}

	for name, attr := range content.Attributes {
		value, d := attr.Expr.Value(ctx)
		if d.HasErrors() {
			return nil, d.Errs()[0]
		}

		switch name {
		case "columns":
			primaryKey.ColumnNames = valueAsStringSlice(value)
		}
	}

	return primaryKey, nil
}

func belongsToFromBlock(block *hcl.Block, ctx *hcl.EvalContext) (*BelongsTo, error) {
	belongsTo := &BelongsTo{
		Name: block.Labels[0],
	}

	content, d := block.Body.Content(hclBelongsTo)
	if d.HasErrors() {
		return nil, d.Errs()[0]
	}

	for name, attr := range content.Attributes {
		value, d := attr.Expr.Value(ctx)
		if d.HasErrors() {
			return nil, d.Errs()[0]
		}

		switch name {
		case "class_name":
			belongsTo.ClassName = value.AsString()
		case "foreign_key":
			belongsTo.ForeignKey = value.AsString()
		}
	}

	return belongsTo, nil
}

func hasManyFromBlock(block *hcl.Block, ctx *hcl.EvalContext) (*HasMany, error) {
	hasMany := &HasMany{
		Name: block.Labels[0],
	}
	hasMany.ClassName = inflect.Pascal(inflect.Singular(hasMany.Name))
	content, d := block.Body.Content(hclHasMany)
	if d.HasErrors() {
		return nil, d.Errs()[0]
	}

	for name, attr := range content.Attributes {
		value, d := attr.Expr.Value(ctx)
		if d.HasErrors() {
			return nil, d.Errs()[0]
		}

		switch name {
		case "through":
			hasMany.Through = value.AsString()
		}
	}

	return hasMany, nil
}

func hasOneFromBlock(block *hcl.Block, ctx *hcl.EvalContext) (*HasOne, error) {
	hasOne := &HasOne{
		Name: block.Labels[0],
	}
	hasOne.ClassName = inflect.Pascal(inflect.Singular(hasOne.Name))
	content, d := block.Body.Content(hclHasOne)
	if d.HasErrors() {
		return nil, d.Errs()[0]
	}

	for name, attr := range content.Attributes {
		value, d := attr.Expr.Value(ctx)
		if d.HasErrors() {
			return nil, d.Errs()[0]
		}

		switch name {
		case "through":
			hasOne.Through = value.AsString()
		}
	}

	return hasOne, nil
}

func columnFromBlock(block *hcl.Block, ctx *hcl.EvalContext) (*Column, error) {
	column := &Column{
		Name: block.Labels[0],
		Null: true,
	}

	content, d := block.Body.Content(hclColumn)
	if d.HasErrors() {
		return nil, d.Errs()[0]
	}
	for name, attr := range content.Attributes {
		value, d := attr.Expr.Value(ctx)
		if d.HasErrors() {
			return nil, d.Errs()[0]
		}

		switch name {
		case "type":
			column.Type = value.AsString()
		case "array":
			column.Array = valueAsBool(value)
		case "null":
			column.Null = valueAsBool(value)
		}
	}
	for name, attr := range content.Attributes {
		value, d := attr.Expr.Value(ctx)
		if d.HasErrors() {
			return nil, d.Errs()[0]
		}
		if name == "default" {
			switch column.Type {
			case "string", "uuid":
				column.Default = value.AsString()
			case "integer":
				column.Default = valueAsInt(value)
			case "boolean":
				column.Default = valueAsBool(value)
			case "float":
				column.Default = valueAsFloat(value)
			}
		}
	}

	return column, nil
}

func attributeFromBlock(block *hcl.Block, ctx *hcl.EvalContext) (*Attribute, error) {
	attribute := &Attribute{
		Name: block.Labels[0],
	}

	content, d := block.Body.Content(hclAttribute)
	if d.HasErrors() {
		return nil, d.Errs()[0]
	}

	for name, attr := range content.Attributes {
		switch name {
		case "type":
			value, d := attr.Expr.Value(ctx)
			if d.HasErrors() {
				return nil, d.Errs()[0]
			}
			attribute.Type = value.AsString()
		}
	}

	return attribute, nil
}

type envKey struct {
	Key string `cty:"key"`
}

var envKeyType = cty.Object(map[string]cty.Type{"key": cty.String})

var env = function.New(&function.Spec{
	VarParam: &function.Parameter{
		Type: cty.String,
	},
	Type: function.StaticReturnType(envKeyType),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		return cty.ObjectVal(map[string]cty.Value{
			"key": cty.StringVal(args[0].AsString()),
		}), nil
	},
})

func Parse(body hcl.Body) (*Schema, error) {
	ctx := &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"id":        cty.StringVal("id"),
			"string":    cty.StringVal("string"),
			"text":      cty.StringVal("text"),
			"boolean":   cty.StringVal("boolean"),
			"date":      cty.StringVal("date"),
			"datetime":  cty.StringVal("datetime"),
			"decimal":   cty.StringVal("decimal"),
			"float":     cty.StringVal("float"),
			"integer":   cty.StringVal("integer"),
			"bigint":    cty.StringVal("bigint"),
			"time":      cty.StringVal("time"),
			"timestamp": cty.StringVal("timestamp"),
			"json":      cty.StringVal("json"),
			"jsonb":     cty.StringVal("jsonb"),
			"uuid":      cty.StringVal("uuid"),
		},
		Functions: map[string]function.Function{
			"env": env,
		},
	}

	sch := NewSchema()

	content, d := body.Content(hclSchema)
	// TODO: diagnostics error handling
	if d.HasErrors() {
		return nil, d.Errs()[0]
	}

	for _, block := range content.Blocks {
		switch block.Type {
		case "database":
			database, err := sch.databaseFromBlock(block, ctx)
			if err != nil {
				return nil, err
			}
			sch.Databases = append(sch.Databases, database)
		}
	}

	return sch, nil
}

func valueAsInt(value cty.Value) int {
	var i int
	gocty.FromCtyValue(value, &i)
	return i
}

func valueAsBool(value cty.Value) bool {
	var b bool
	gocty.FromCtyValue(value, &b)
	return b
}
func valueAsFloat(value cty.Value) float64 {
	var b float64
	gocty.FromCtyValue(value, &b)
	return b
}

func valueAsStringSlice(value cty.Value) []string {
	var ss []string
	value.ForEachElement(func(key, val cty.Value) (stop bool) {
		var column string
		gocty.FromCtyValue(val, &column)
		ss = append(ss, column)
		return false
	})
	return ss
}

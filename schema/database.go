package schema

type Database struct {
	Name       string
	Adapter    string
	TimeZone   string
	Generators []*Generator
	Configs    []*Config
	Models     []*Model
	models     map[string]*Model
}

func (d *Database) LoadConfig(environment string) *Config {
	for _, config := range d.Configs {
		if config.Environment == environment {
			return config
		}
	}
	return nil
}

func (d *Database) Tables() []string {
	tables := []string{}
	for _, model := range d.Models {
		tables = append(tables, model.TableName)
	}
	return tables
}

func (d *Database) buildAssociation() {
	for _, m := range d.Models { // User
		for _, h := range m.HasMany { // has_many
			if m1, ok := d.models[h.ModelName]; ok { // Post
				for _, b := range m1.BelongsTo { // belongs_to
					if b.ModelName == m.Name { // User
						h.BelongsTo = b
					}
				}
			}
		}
		for _, h := range m.HasOne {
			if m1, ok := d.models[h.ModelName]; ok {
				for _, b := range m1.BelongsTo {
					if b.ModelName == m.Name {
						h.BelongsTo = b
					}
				}
			}
		}
	}
}

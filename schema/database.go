package schema

type Database struct {
	Name       string
	Adapter    string
	TimeZone   string
	Generators []*Generator
	Configs    []*Config
	Models     []*Model
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

package schema

type Schema struct {
	Databases []*Database
}

type Generator struct {
	Name string
	Test bool
}

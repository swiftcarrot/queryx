package schema

import "testing"

func TestPostgreSQLWithoutPrimaryKey(t *testing.T) {
	schema := NewSchema()
	database := schema.NewDatabase("test")
	user := database.NewModel("User")

	user.AddColumn(&Column{Name: "name", Type: "string"})
	database.CreatePostgreSQLSchema("")
}

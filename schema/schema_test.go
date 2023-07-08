package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTableWithoutPrimaryKey(t *testing.T) {
	schema := NewSchema()
	database := schema.NewDatabase("test")
	user := database.NewModel("User")
	user.AddColumn(&Column{Name: "name", Type: "string"})

	require.Nil(t, user.PrimaryKey)
	database.CreatePostgreSQLSchema("test")
	database.CreateMySQLSchema("test")
	database.CreateSQLiteSchema("test")
}

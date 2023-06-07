package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewModel(t *testing.T) {
	schema := NewSchema()
	database := schema.NewDatabase("test")
	user := database.NewModel("User")

	require.Equal(t, "User", user.Name)
	require.Equal(t, "users", user.TableName)

	require.Equal(t, 0, len(user.Columns))

	user.AddDefaultPrimaryKey()

	require.Equal(t, 1, len(user.Columns))
	require.Equal(t, "id", user.Columns[0].Name)
	require.Equal(t, []string{"id"}, user.PrimaryKey.ColumnNames)

	user.AddTimestamps()
	require.Equal(t, 3, len(user.Columns))
	require.Equal(t, "created_at", user.Columns[1].Name)
	require.Equal(t, "updated_at", user.Columns[2].Name)

	post := database.NewModel("Post")
	post.AddBelongsTo(&BelongsTo{
		Name: "user",
	})
	require.Equal(t, 1, len(post.Columns))
	require.Equal(t, "user_id", post.Columns[0].Name)
}

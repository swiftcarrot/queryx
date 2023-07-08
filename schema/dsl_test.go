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

	require.Equal(t, "User(id: bigint, created_at: datetime, updated_at: datetime)", user.String())

	post := database.NewModel("Post")
	b := &BelongsTo{Name: "user"}
	post.AddBelongsTo(b)
	require.Equal(t, 1, len(post.Columns))
	require.Equal(t, "user_id", post.Columns[0].Name)
	require.Equal(t, []*BelongsTo{
		{Name: "user", ModelName: "User", ForeignKey: "user_id"},
	}, post.BelongsTo)

	user.AddHasMany(&HasMany{Name: "posts"})
	user.AddHasOne(&HasOne{Name: "account"})
	require.Equal(t, []*HasMany{
		{Name: "posts", ModelName: "Post", ForeignKey: "user_id", BelongsTo: b},
	}, user.HasMany)
	require.Equal(t, []*HasOne{
		{Name: "account", ModelName: "Account", ForeignKey: "user_id"},
	}, user.HasOne)
}

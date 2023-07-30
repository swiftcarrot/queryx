package inflect

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSnake(t *testing.T) {
	require.Equal(t, "user", Snake("User"))
	require.Equal(t, "user_post", Snake("UserPost"))
}

func TestCamel(t *testing.T) {
	require.Equal(t, "user", Camel("user"))
	require.Equal(t, "userPost", Camel("user_post"))
	require.Equal(t, "userPost", Camel("UserPost"))
}

func TestPascal(t *testing.T) {
	require.Equal(t, "User", Pascal("user"))
	require.Equal(t, "UserPost", Pascal("user_post"))
	require.Equal(t, "UserPost", Pascal("userPost"))
}

func TestPlural(t *testing.T) {
	require.Equal(t, "users", Plural("user"))
	require.Equal(t, "Users", Plural("User"))
	require.Equal(t, "user_posts", Plural("user_post"))
	require.Equal(t, "userPosts", Plural("userPost"))
	require.Equal(t, "moneySlice", Plural("money"))
}

func TestSingular(t *testing.T) {
	require.Equal(t, "user_post", Singular("user_posts"))
}

package queryx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewInsert(t *testing.T) {
	s1 := NewInsert().Into("users")
	sql, args := s1.ToSQL()
	require.Equal(t, "INSERT INTO users DEFAULT VALUES", sql)
	require.Equal(t, []interface{}{}, args)

	s2 := NewInsert().Into("users").Columns("name", "email").Values("test", "test@example.com")
	sql, args = s2.ToSQL()
	require.Equal(t, "INSERT INTO users(name, email) VALUES (?, ?)", sql)
	require.Equal(t, []interface{}{"test", "test@example.com"}, args)
}

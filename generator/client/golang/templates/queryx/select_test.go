package queryx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSelect(t *testing.T) {
	s := NewSelect().Select("users.*").From("users")
	sql, args := s.ToSQL()
	require.Equal(t, `SELECT users.* FROM users`, sql)
	require.Equal(t, []interface{}{}, args)

	sql, args = s.Update().Columns("name", "email").Values("test", "test@example.com").ToSQL()
	require.Equal(t, `UPDATE users SET name = ?, email = ?`, sql)
	require.Equal(t, []interface{}{"test", "test@example.com"}, args)
}

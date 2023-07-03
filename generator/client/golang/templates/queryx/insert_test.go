package queryx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewInsert(t *testing.T) {
	s := NewInsert().Into("users").Columns("name", "email").Values("test", "test@example.com")
	sql, args := s.ToSQL()
	require.Equal(t, "INSERT INTO users(name, email) VALUES (?, ?)", sql)
	require.Equal(t, []interface{}{"test", "test@example.com"}, args)
}

package queryx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewJSON(t *testing.T) {
	m := map[string]interface{}{"a": 1}
	j := NewJSON(m)
	require.Equal(t, m, j.Val)
	require.False(t, j.Null)
}

func TestNewNullableJSON(t *testing.T) {
	j1 := NewNullableJSON(nil)
	require.True(t, j1.Null)

	m := map[string]interface{}{"a": 1}
	j2 := NewNullableJSON(m)
	require.False(t, j2.Null)
}

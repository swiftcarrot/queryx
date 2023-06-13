package queryx

import (
	"encoding/json"
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

func TestJSONScan(t *testing.T) {
	m := map[string]interface{}{"a": 1}
	j := NewJSON(m)
	m2 := map[string]interface{}{"a": 2}
	bytes, err := json.Marshal(m2)
	require.NoError(t, err)
	err = j.Scan(bytes)
	require.NoError(t, err)
	require.Equal(t, 2, m2["a"])
}

func TestJSONValue(t *testing.T) {
	m := map[string]interface{}{"a": 1}
	j := NewJSON(m)
	marshal, err := json.Marshal(j.Val)
	require.NoError(t, err)
	value, err := j.Value()
	require.NoError(t, err)
	require.Equal(t, value, marshal)
}

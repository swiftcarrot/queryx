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

func TestJSONMarshalJSON(t *testing.T) {
	m := map[string]interface{}{"a": 1}
	j := NewJSON(m)
	require.Equal(t, 1, j.Val["a"])
	require.Equal(t, false, j.Null)
	_, err := j.MarshalJSON()
	require.NoError(t, err)
}

func TestJSONUnmarshalJSON(t *testing.T) {
	m := map[string]interface{}{"a": 1}
	j := NewJSON(m)
	require.Equal(t, 1, j.Val["a"])
	require.Equal(t, false, j.Null)
	bytes, err := j.MarshalJSON()
	require.NoError(t, err)
	m2 := map[string]interface{}{}
	newJson := NewJSON(m2)
	err = newJson.UnmarshalJSON(bytes)
	require.NoError(t, err)
	require.Equal(t, float64(1), newJson.Val["a"])
	require.Equal(t, false, newJson.Null)
}

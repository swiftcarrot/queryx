// Code generated by queryx, DO NOT EDIT.

package queryx

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewBoolean(t *testing.T) {
	b1 := NewBoolean(true)
	require.Equal(t, true, b1.Val)
	require.True(t, b1.Valid)

	b2 := NewNullableBoolean(nil)
	require.False(t, b2.Valid)
}

func TestBooleanJSON(t *testing.T) {
	type Foo struct {
		X Boolean `json:"x"`
		Y Boolean `json:"y"`
	}
	x := NewBoolean(true)
	y := NewNullableBoolean(nil)
	s := `{"x":true,"y":null}`

	f1 := Foo{X: x, Y: y}
	b, err := json.Marshal(f1)
	require.NoError(t, err)
	require.Equal(t, s, string(b))

	var f2 Foo
	err = json.Unmarshal([]byte(s), &f2)
	require.NoError(t, err)
	require.Equal(t, x, f2.X)
	require.Equal(t, y, f2.Y)
}

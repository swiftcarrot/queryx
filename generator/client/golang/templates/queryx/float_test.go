// Code generated by queryx, DO NOT EDIT.

package queryx

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFloat(t *testing.T) {
	f1 := NewFloat(2.1)
	require.Equal(t, 2.1, f1.Val)
	require.True(t, f1.Valid)

	f2 := NewNullableFloat(nil)
	require.False(t, f2.Valid)
}

func TestFloatJSON(t *testing.T) {
	type Foo struct {
		X Float `json:"x"`
		Y Float `json:"y"`
	}
	x := NewFloat(2.1)
	y := NewNullableFloat(nil)
	s := `{"x":2.1,"y":null}`

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

package queryx

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewVector(t *testing.T) {
	v1 := NewVector([]float32{1, 2, 3})
	require.Equal(t, []float32{1, 2, 3}, v1.Val)
	require.Equal(t, false, v1.Null)

	v2 := NewNullableVector(nil)
	require.Equal(t, true, v2.Null)
}

func TestVectorJSON(t *testing.T) {
	type Foo struct {
		X Vector `json:"x"`
		Y Vector `json:"y"`
	}
	x := NewVector([]float32{1, 2, 3})
	y := NewNullableVector(nil)
	s := `{"x":[1,2,3],"y":null}`

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

// Code generated by queryx, DO NOT EDIT.

package queryx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewFloat(t *testing.T) {
	i := NewFloat(2.1)
	require.Equal(t, 2.1, i.Val)
	require.Equal(t, false, i.Null)
}

func TestNewNullableFloat(t *testing.T) {
	i := NewNullableFloat(nil)
	require.Equal(t, true, i.Null)
}

func TestFloatScan(t *testing.T) {
	i := NewFloat(2.1)
	require.Equal(t, float64(2.1), i.Val)
	require.Equal(t, false, i.Null)
	err := i.Scan(3.1)
	require.NoError(t, err)
	require.Equal(t, float64(3.1), i.Val)
}

func TestFloatValue(t *testing.T) {
	i := NewFloat(2.1)
	require.Equal(t, float64(2.1), i.Val)
	require.Equal(t, false, i.Null)
	value, err := i.Value()
	require.NoError(t, err)
	require.Equal(t, float64(2.1), value)
}
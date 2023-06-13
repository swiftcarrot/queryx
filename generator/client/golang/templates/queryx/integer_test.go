// Code generated by queryx, DO NOT EDIT.

package queryx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewInteger(t *testing.T) {
	i := NewInteger(2)
	require.Equal(t, int32(2), i.Val)
	require.Equal(t, false, i.Null)
}

func TestNewNullableInteger(t *testing.T) {
	i := NewNullableInteger(nil)
	require.Equal(t, true, i.Null)
}

func TestIntegerScan(t *testing.T) {
	i := NewInteger(2)
	require.Equal(t, int32(2), i.Val)
	require.Equal(t, false, i.Null)
	err := i.Scan(3)
	require.NoError(t, err)
	require.Equal(t, int32(3), i.Val)
}

func TestIntegerValue(t *testing.T) {
	i := NewInteger(2)
	require.Equal(t, int32(2), i.Val)
	require.Equal(t, false, i.Null)
	value, err := i.Value()
	require.NoError(t, err)
	require.Equal(t, int64(2), value)
}
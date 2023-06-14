// Code generated by queryx, DO NOT EDIT.

package queryx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUUID(t *testing.T) {
	i := NewUUID("a81e44c5-7e18-4dfe-b9b3-d9280629d2ef")
	require.Equal(t, "a81e44c5-7e18-4dfe-b9b3-d9280629d2ef", i.Val)
	require.Equal(t, false, i.Null)
}

func TestNewNullableUUID(t *testing.T) {
	i := NewNullableUUID(nil)
	require.Equal(t, true, i.Null)
}

func TestUUIDMarshalJSON(t *testing.T) {
	i := NewUUID("a81e44c5-7e18-4dfe-b9b3-d9280629d2ef")
	require.Equal(t, "a81e44c5-7e18-4dfe-b9b3-d9280629d2ef", i.Val)
	require.Equal(t, false, i.Null)
	_, err := i.MarshalJSON()
	require.NoError(t, err)
}

func TestUUIDUnmarshalJSON(t *testing.T) {
	i := NewUUID("a81e44c5-7e18-4dfe-b9b3-d9280629d2ef")
	require.Equal(t, "a81e44c5-7e18-4dfe-b9b3-d9280629d2ef", i.Val)
	require.Equal(t, false, i.Null)
	bytes, err := i.MarshalJSON()
	require.NoError(t, err)
	u := NewUUID("")
	err = u.UnmarshalJSON(bytes)
	require.NoError(t, err)
	require.Equal(t, "a81e44c5-7e18-4dfe-b9b3-d9280629d2ef", u.Val)
	require.Equal(t, false, u.Null)
}

func TestUUIDScan(t *testing.T) {
	i := NewUUID("a81e44c5-7e18-4dfe-b9b3-d9280629d2ef")
	require.Equal(t, "a81e44c5-7e18-4dfe-b9b3-d9280629d2ef", i.Val)
	require.Equal(t, false, i.Null)
	err := i.Scan("a81e44c5-7e18-4dfe-b9b3-d9280629dfff")
	require.NoError(t, err)
	require.Equal(t, "a81e44c5-7e18-4dfe-b9b3-d9280629dfff", i.Val)
}

func TestUUIDValue(t *testing.T) {
	i := NewUUID("a81e44c5-7e18-4dfe-b9b3-d9280629d2ef")
	require.Equal(t, "a81e44c5-7e18-4dfe-b9b3-d9280629d2ef", i.Val)
	require.Equal(t, false, i.Null)
	value, err := i.Value()
	require.NoError(t, err)
	require.Equal(t, "a81e44c5-7e18-4dfe-b9b3-d9280629d2ef", value)
}

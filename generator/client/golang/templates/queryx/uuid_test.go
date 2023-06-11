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

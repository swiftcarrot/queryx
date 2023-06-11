// Code generated by queryx, DO NOT EDIT.

package queryx

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewBoolean(t *testing.T) {
	i := NewBoolean(true)
	require.Equal(t, true, i.Val)
	require.Equal(t, false, i.Null)
}

func TestNewNullableBoolean(t *testing.T) {
	i := NewNullableBoolean(nil)
	require.Equal(t, true, i.Null)
}

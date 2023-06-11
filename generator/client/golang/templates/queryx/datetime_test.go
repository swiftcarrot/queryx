// Code generated by queryx, DO NOT EDIT.

package queryx

import (
	"github.com/swiftcarrot/queryx/internal/integration/db/queryx"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDatetime(t *testing.T) {
	i := NewDatetime("2012-12-12 15:04:05")
	require.Equal(t, "2012-12-12 15:04:05", i.Val.Local().Format("2006-01-02 15:04:05"))
	require.Equal(t, false, i.Null)
}

func TestNewNullableDatetime(t *testing.T) {
	i := NewNullableDatetime(nil)
	require.Equal(t, true, i.Null)
}

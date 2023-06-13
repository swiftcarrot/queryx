// Code generated by queryx, DO NOT EDIT.

package queryx

import (
	"testing"
	"time"

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

func TestDatetimeMarshalJSON(t *testing.T) {
	i := NewDatetime("2012-12-12 15:04:05")
	require.Equal(t, "2012-12-12 15:04:05", i.Val.Local().Format("2006-01-02 15:04:05"))
	require.Equal(t, false, i.Null)
	_, err := i.MarshalJSON()
	require.NoError(t, err)
}

func TestDatetimeScan(t *testing.T) {
	i := NewDatetime("2012-12-12 15:04:05")
	require.Equal(t, "2012-12-12 15:04:05", i.Val.Local().Format("2006-01-02 15:04:05"))
	require.Equal(t, false, i.Null)
	date, err := parseDatetime("2012-12-12 15:04:05")
	require.NoError(t, err)
	err = i.Scan(*date)
	require.NoError(t, err)
	require.Equal(t, "2012-12-12 15:04:05", i.Val.Local().Format("2006-01-02 15:04:05"))
}

func TestDatetimeValue(t *testing.T) {
	i := NewDatetime("2012-12-12 15:04:05")
	require.Equal(t, "2012-12-12 15:04:05", i.Val.Local().Format("2006-01-02 15:04:05"))
	require.Equal(t, false, i.Null)
	value, err := i.Value()
	require.NoError(t, err)
	_value := value.(time.Time)
	require.Equal(t, "2012-12-12 15:04:05", _value.Local().Format("2006-01-02 15:04:05"))
}

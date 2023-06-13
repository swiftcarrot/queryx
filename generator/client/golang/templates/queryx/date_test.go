// Code generated by queryx, DO NOT EDIT.

package queryx

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewDate(t *testing.T) {
	i := NewDate("2012-12-12")
	require.Equal(t, "2012-12-12", i.Val.Local().Format("2006-01-02"))
	require.Equal(t, false, i.Null)
}

func TestNewNullableDate(t *testing.T) {
	i := NewNullableDate(nil)
	require.Equal(t, true, i.Null)
}

func TestDateMarshalJSON(t *testing.T) {
	i := NewDate("2012-12-12")
	require.Equal(t, "2012-12-12", i.Val.Local().Format("2006-01-02"))
	require.Equal(t, false, i.Null)
	_, err := i.MarshalJSON()
	require.NoError(t, err)
}

func TestDateScan(t *testing.T) {
	i := NewDate("2012-12-12")
	require.Equal(t, "2012-12-12", i.Val.Local().Format("2006-01-02"))
	require.Equal(t, false, i.Null)
	date, err := parseDate("2012-11-13")
	require.NoError(t, err)
	err = i.Scan(*date)
	require.NoError(t, err)
	require.Equal(t, "2012-11-13", i.Val.Local().Format("2006-01-02"))
}

func TestDateValue(t *testing.T) {
	i := NewDate("2012-12-12")
	require.Equal(t, "2012-12-12", i.Val.Local().Format("2006-01-02"))
	require.Equal(t, false, i.Null)
	value, err := i.Value()
	require.NoError(t, err)
	_value := value.(time.Time)
	require.Equal(t, "2012-12-12", _value.Local().Format("2006-01-02"))
}
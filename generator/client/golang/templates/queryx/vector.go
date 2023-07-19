package queryx

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Vector struct {
	Val  []float32
	Null bool
	Set  bool
}

func NewVector(v []float32) Vector {
	return Vector{Val: v, Set: true}
}

func NewNullableVector(v *[]float32) Vector {
	if v != nil {
		return Vector{Val: *v, Set: true}
	}
	return Vector{Null: true, Set: true}
}

func (v Vector) String() string {
	var buf strings.Builder
	buf.WriteString("[")

	for i := 0; i < len(v.Val); i++ {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.FormatFloat(float64(v.Val[i]), 'f', -1, 32))
	}

	buf.WriteString("]")
	return buf.String()
}

func (v *Vector) Parse(s string) error {
	v.Val = make([]float32, 0)
	sp := strings.Split(s[1:len(s)-1], ",")
	for i := 0; i < len(sp); i++ {
		n, err := strconv.ParseFloat(sp[i], 32)
		if err != nil {
			return err
		}
		v.Val = append(v.Val, float32(n))
	}
	return nil
}

// Scan implements the Scanner interface.
func (v *Vector) Scan(src interface{}) (err error) {
	switch src := src.(type) {
	case []byte:
		return v.Parse(string(src))
	case string:
		return v.Parse(src)
	default:
		return fmt.Errorf("unsupported data type: %T", src)
	}
}

// Value implements the driver Valuer interface.
func (v Vector) Value() (driver.Value, error) {
	return v.String(), nil
}

// MarshalJSON implements the json.Marshaler interface.
func (v Vector) MarshalJSON() ([]byte, error) {
	if v.Null {
		return json.Marshal(nil)
	}
	return json.Marshal(v.Val)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (v *Vector) UnmarshalJSON(data []byte) error {
	v.Set = true
	if string(data) == "null" {
		v.Null = true
		return nil
	}
	if err := json.Unmarshal(data, &v.Val); err != nil {
		return err
	}
	return nil
}

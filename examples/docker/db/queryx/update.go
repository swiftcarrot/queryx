// Code generated by queryx, DO NOT EDIT.

package queryx

import (
	"fmt"
	"strings"
)

type UpdateStatement struct {
	table     string
	columns   []string
	values    []interface{}
	where     *Clause
	returning string
}

func NewUpdate() *UpdateStatement {
	return &UpdateStatement{}
}

func (s *UpdateStatement) Table(table string) *UpdateStatement {
	s.table = table
	return s
}

func (s *UpdateStatement) Columns(columns ...string) *UpdateStatement {
	s.columns = columns
	return s
}

func (s *UpdateStatement) Values(values ...interface{}) *UpdateStatement {
	s.values = values
	return s
}

func (s *UpdateStatement) Where(expr *Clause) *UpdateStatement {
	s.where = expr
	return s
}

func (s *UpdateStatement) Returning(returning string) *UpdateStatement {
	s.returning = returning
	return s
}

func (s *UpdateStatement) ToSQL() (string, []interface{}) {
	sql, args := fmt.Sprintf("UPDATE %s SET", s.table), s.values

	sets := []string{}
	for _, col := range s.columns {
		sets = append(sets, fmt.Sprintf("%s = ?", col))
	}

	sql = fmt.Sprintf("%s %s", sql, strings.Join(sets, ", "))

	if s.where != nil {
		sql = fmt.Sprintf("%s WHERE %s", sql, s.where.fragment)
		args = append(args, s.where.args...)
	}

	if s.returning != "" {
		sql = fmt.Sprintf("%s RETURNING %s", sql, s.returning)
	}

	return sql, args
}

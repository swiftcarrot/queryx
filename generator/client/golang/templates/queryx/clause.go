// Code generated by queryx, DO NOT EDIT.

package queryx

import (
	"fmt"
	"strings"
)

type Clause struct {
	fragment string
	args     []interface{}
	err      error
}

func NewClause(fragment string, args []interface{}) *Clause {
	return &Clause{
		fragment: fragment,
		args:     args,
	}
}

func (c *Clause) Err() error {
	return c.err
}

func (c *Clause) And(clauses ...*Clause) *Clause {
	var fragments []string
	var args []interface{}
	clauses = append([]*Clause{c}, clauses...)
	for _, clause := range clauses {
		fragments = append(fragments, fmt.Sprintf("(%s)", clause.fragment))
		args = append(args, clause.args...)
	}

	return &Clause{
		fragment: strings.Join(fragments, " AND "),
		args:     args,
	}
}

func (c *Clause) Or(clauses ...*Clause) *Clause {
	var fragments []string
	var args []interface{}
	clauses = append([]*Clause{c}, clauses...)
	for _, clause := range clauses {
		fragments = append(fragments, fmt.Sprintf("(%s)", clause.fragment))
		args = append(args, clause.args...)
	}

	return &Clause{
		fragment: strings.Join(fragments, " OR "),
		args:     args,
	}
}

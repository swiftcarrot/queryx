// Code generated by queryx, DO NOT EDIT.

package queryx

import "fmt"

type StringColumn struct {
	Name  string
	Table *Table
}

func (t *Table) NewStringColumn(name string) *StringColumn {
	return &StringColumn{
		Table: t,
		Name:  name,
	}
}

func (c *StringColumn) NE(v string) *Clause {
	return &Clause{
		fragment: fmt.Sprintf("%s.%s <> ?", c.Table.Name, c.Name),
		args:     []interface{}{v},
	}
}
func (c *StringColumn) EQ(v string) *Clause {
	return &Clause{
		fragment: fmt.Sprintf("lower(%s.%s) = lower(?)", c.Table.Name, c.Name),
		args:     []interface{}{v},
	}
}

func (c *StringColumn) IEQ(v string) *Clause {
	return &Clause{
		fragment: fmt.Sprintf("%s.%s = ?", c.Table.Name, c.Name),
		args:     []interface{}{v},
	}
}

func (c *StringColumn) In(v []string) *Clause {
	if len(v) == 0 {
		return &Clause{
			fragment: "1=0",
		}
	}
	return &Clause{
		fragment: fmt.Sprintf("%s.%s IN (?)", c.Table.Name, c.Name),
		args:     []interface{}{v},
	}
}

func (c *StringColumn) NIn(v []string) *Clause {
	if len(v) == 0 {
		return &Clause{
			fragment: "1!=0",
		}
	}
	return &Clause{
		fragment: fmt.Sprintf("%s.%s NOT IN (?)", c.Table.Name, c.Name),
		args:     []interface{}{v},
	}
}

func (c *StringColumn) Like(v string) *Clause {
	return &Clause{
		fragment: fmt.Sprintf("%s.%s ilike ?", c.Table.Name, c.Name),
		args:     []interface{}{fmt.Sprintf("%s%s%s", "%", v, "%")},
	}
}

func (c *StringColumn) ILike(v string) *Clause {
	return &Clause{
		fragment: fmt.Sprintf("%s.%s like ?", c.Table.Name, c.Name),
		args:     []interface{}{fmt.Sprintf("%s%s%s", "%", v, "%")},
	}
}
func (c *StringColumn) IStartsWith(v string) *Clause {
	return &Clause{
		fragment: fmt.Sprintf("%s.%s like ?", c.Table.Name, c.Name),
		args:     []interface{}{fmt.Sprintf("%s%s", v, "%")},
	}
}

func (c *StringColumn) StartsWith(v string) *Clause {
	return &Clause{
		fragment: fmt.Sprintf("%s.%s ilike ?", c.Table.Name, c.Name),
		args:     []interface{}{fmt.Sprintf("%s%s", v, "%")},
	}
}

func (c *StringColumn) EndsWith(v string) *Clause {
	return &Clause{
		fragment: fmt.Sprintf("%s.%s ilike ?", c.Table.Name, c.Name),
		args:     []interface{}{fmt.Sprintf("%s%s", "%", v)},
	}
}

func (c *StringColumn) IEndsWith(v string) *Clause {
	return &Clause{
		fragment: fmt.Sprintf("%s.%s like ?", c.Table.Name, c.Name),
		args:     []interface{}{fmt.Sprintf("%s%s", "%", v)},
	}
}

func (c *StringColumn) Asc() string {
	return fmt.Sprintf("%s.%s ASC", c.Table.Name, c.Name)
}

func (c *StringColumn) Desc() string {
	return fmt.Sprintf("%s.%s DESC", c.Table.Name, c.Name)
}

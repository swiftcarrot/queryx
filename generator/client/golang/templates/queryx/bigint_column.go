// Code generated by queryx, DO NOT EDIT.

package queryx

import "fmt"

type BigIntColumn struct {
	Name  string
	Table *Table
}

func (t *Table) NewBigIntColumn(name string) *BigIntColumn {
	return &BigIntColumn{
		Table: t,
		Name:  name,
	}
}

func (c *BigIntColumn) EQ(v int64) *Clause {
	return &Clause{
		fragment: fmt.Sprintf("%s.%s = ?", c.Table.Name, c.Name),
		args:     []interface{}{v},
	}
}

func (c *BigIntColumn) NE(v int64) *Clause {
	return &Clause{
		fragment: fmt.Sprintf("%s.%s <> ?", c.Table.Name, c.Name),
		args:     []interface{}{v},
	}
}

func (c *BigIntColumn) LT(v int64) *Clause {
	return &Clause{
		fragment: fmt.Sprintf("%s.%s < ?", c.Table.Name, c.Name),
		args:     []interface{}{v},
	}
}

func (c *BigIntColumn) GT(v int64) *Clause {
	return &Clause{
		fragment: fmt.Sprintf("%s.%s > ?", c.Table.Name, c.Name),
		args:     []interface{}{v},
	}
}

func (c *BigIntColumn) LE(v int64) *Clause {
	return &Clause{
		fragment: fmt.Sprintf("%s.%s <= ?", c.Table.Name, c.Name),
		args:     []interface{}{v},
	}
}

func (c *BigIntColumn) GE(v int64) *Clause {
	return &Clause{
		fragment: fmt.Sprintf("%s.%s >= ?", c.Table.Name, c.Name),
		args:     []interface{}{v},
	}
}

func (c *BigIntColumn) In(v []int64) *Clause {
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

func (c *BigIntColumn) NIn(v []int64) *Clause {
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

func (c *BigIntColumn) InRange(start int64, end int64) *Clause {
	return &Clause{
		fragment: fmt.Sprintf("%s.%s >= ? and %s.%s< ?", c.Table.Name, c.Name, c.Table.Name, c.Name),
		args:     []interface{}{start, end},
	}
}

func (c *BigIntColumn) Between(start int64, end int64) *Clause {
	return &Clause{
		fragment: fmt.Sprintf("%s.%s >= ? and %s.%s<= ?", c.Table.Name, c.Name, c.Table.Name, c.Name),
		args:     []interface{}{start, end},
	}
}

func (c *BigIntColumn) Asc() string {
	return fmt.Sprintf("%s.%s ASC", c.Table.Name, c.Name)
}

func (c *BigIntColumn) Desc() string {
	return fmt.Sprintf("%s.%s DESC", c.Table.Name, c.Name)
}

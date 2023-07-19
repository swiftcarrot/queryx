package queryx

type VectorColumn struct {
	Name  string
	Table *Table
}

func (t *Table) NewVectorColumn(name string) *VectorColumn {
	return &VectorColumn{
		Table: t,
		Name:  name,
	}
}

package go_queryx

import (
	"github.com/swiftcarrot/queryx/internal/benchmarks/go-queryx/db"
	"github.com/swiftcarrot/queryx/internal/benchmarks/go-queryx/db/queryx"
	"github.com/swiftcarrot/queryx/internal/benchmarks/go-queryx/helper"
	"testing"
)

var (
	c *db.QXClient
)

const (
	queryxSelectMultiSQL = `SELECT * FROM models WHERE id > 0`
)

type Queryx struct {
	helper.ORMInterface
}

func CreateQueryx() helper.ORMInterface {
	return &Queryx{}
}

func (Queryx *Queryx) Name() string {
	return "queryx"
}

func (Queryx *Queryx) Init() (*db.QXClient, error) {
	client, err := db.NewClient()
	if err != nil {
		return nil, err
	}
	c = client
	return client, err
}

func (Queryx *Queryx) Create(b *testing.B) {
	m := NewModel()

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := c.QueryModel().Create(c.ChangeModel().SetName(m.Name).
			SetTitle(m.Title).SetFax(m.Fax).SetWeb(m.Web).SetAge(m.Age).SetCounter(m.Counter).SetRigh(m.Righ))
		if err != nil {
			helper.SetError(b, Queryx.Name(), "Create", err.Error())
		}
	}
}

func (Queryx *Queryx) InsertAll(b *testing.B) {
	m := NewModel()
	ms := make([]*queryx.ModelChange, 0, 100)
	for i := 0; i < 100; i++ {
		ms = append(ms, c.ChangeModel().SetName(m.Name).
			SetTitle(m.Title).SetFax(m.Fax).SetWeb(m.Web).SetAge(m.Age).SetCounter(m.Counter).SetRigh(m.Righ))
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := c.QueryModel().InsertAll(ms)
		if err != nil {
			helper.SetError(b, Queryx.Name(), "InsertAll", err.Error())
		}
	}
}

func (Queryx *Queryx) Update(b *testing.B) {
	m := NewModel()

	change := c.ChangeModel().SetName(m.Name).SetRigh(m.Righ).
		SetTitle(m.Title).SetFax(m.Fax).SetWeb(m.Web).SetAge(m.Age).SetCounter(m.Counter)
	m8, err := c.QueryModel().Create(change)
	if err != nil {
		helper.SetError(b, Queryx.Name(), "Update", err.Error())
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := c.QueryModel().Where(c.ModelID.EQ(m8.ID)).UpdateAll(change)
		if err != nil {
			helper.SetError(b, Queryx.Name(), "UpdateAll", err.Error())
		}
	}
}

func (Queryx *Queryx) Read(b *testing.B) {
	m := NewModel()
	change := c.ChangeModel().SetName(m.Name).SetRigh(m.Righ).
		SetTitle(m.Title).SetFax(m.Fax).SetWeb(m.Web).SetAge(m.Age).SetCounter(m.Counter)
	_, err := c.QueryModel().Create(change)
	if err != nil {
		helper.SetError(b, Queryx.Name(), "Read", err.Error())
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := c.QueryModel().FindBy(c.ModelName.EQ(m.Name))
		if err != nil {
			helper.SetError(b, Queryx.Name(), "FindBy", err.Error())
		}
	}
}

func (Queryx *Queryx) ReadSlice(b *testing.B) {
	m := NewModel()
	change := c.ChangeModel().SetName(m.Name).SetRigh(m.Righ).
		SetTitle(m.Title).SetFax(m.Fax).SetWeb(m.Web).SetAge(m.Age).SetCounter(m.Counter)
	for i := 0; i < 100; i++ {
		_, err := c.QueryModel().Create(change)
		if err != nil {
			helper.SetError(b, Queryx.Name(), "ReadSlice", err.Error())
		}
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := c.QueryModel().FindBySQL(queryxSelectMultiSQL)
		if err != nil {
			helper.SetError(b, Queryx.Name(), "FindBySQL", err.Error())
		}
	}
}

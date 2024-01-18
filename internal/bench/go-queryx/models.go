package go_queryx

// Model for go-queryx
type Model8 struct {
	ID      int    `column:"id"`
	Name    string `column:"name"`
	Title   string `column:"title"`
	Fax     string `column:"fax"`
	Web     string `column:"web"`
	Age     int    `column:"age"`
	Righ    bool   `column:"righ"`
	Counter int64  `column:"counter"`
}

func (entity *Model8) GetTableName() string {
	return "models"
}

func (entity *Model8) GetPKColumnName() string {
	return "id"
}

func NewModel8() *Model8 {
	m := new(Model8)
	m.Name = "Orm Benchmark"
	m.Title = "Just a Benchmark for fun"
	m.Fax = "99909990"
	m.Web = "http://blog.milkpod29.me"
	m.Age = 100
	m.Righ = true
	m.Counter = 1000

	return m
}

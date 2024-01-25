package go_queryx

// Model for go-queryx
type Model struct {
	ID      int    `column:"id"`
	Name    string `column:"name"`
	Title   string `column:"title"`
	Fax     string `column:"fax"`
	Web     string `column:"web"`
	Age     int    `column:"age"`
	Righ    bool   `column:"righ"`
	Counter int64  `column:"counter"`
}

func NewModel() *Model {
	m := new(Model)
	m.Name = "Orm Benchmark"
	m.Title = "Just a Benchmark for fun"
	m.Fax = "99909990"
	m.Web = "http://blog.milkpod29.me"
	m.Age = 100
	m.Righ = true
	m.Counter = 1000

	return m
}

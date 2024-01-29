package helper

import (
	"github.com/swiftcarrot/queryx/internal/benchmarks/go-queryx/db"
	"testing"
)

type ORMInterface interface {
	Name() string
	Init() (*db.QXClient, error)
	Create(b *testing.B)
	InsertAll(b *testing.B)
	Update(b *testing.B)
	Read(b *testing.B)
	ReadSlice(b *testing.B)
}

type BenchmarkResult struct {
	ORM     string
	Results []Result
}

type Result struct {
	Name     string
	Method   string
	ErrorMsg string

	N         int
	NsPerOp   int64
	MemAllocs int64
	MemBytes  int64
}

type BenchmarkReport []*Result

func (s BenchmarkReport) Len() int { return len(s) }

func (s BenchmarkReport) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s BenchmarkReport) Less(i, j int) bool {
	if s[i].ErrorMsg != "" {
		return false
	}
	if s[j].ErrorMsg != "" {
		return true
	}
	return s[i].NsPerOp < s[j].NsPerOp
}

func RunBenchmarks(adapter string, orm ORMInterface, reports map[string]BenchmarkReport) error {
	c, err := orm.Init()
	if err != nil {
		return err
	}

	operations := []func(b *testing.B){orm.InsertAll, orm.Create, orm.Update, orm.Read, orm.ReadSlice}

	for _, operation := range operations {
		err = CreateTables(c, adapter)
		if err != nil {
			return err
		}

		br := testing.Benchmark(operation)
		method := getFuncName(operation)

		gotResult := &Result{
			Name:      orm.Name(),
			Method:    method,
			ErrorMsg:  GetError(orm.Name(), method),
			N:         br.N,
			NsPerOp:   br.NsPerOp(),
			MemAllocs: br.AllocsPerOp(),
			MemBytes:  br.AllocedBytesPerOp(),
		}

		reports[method] = append(reports[method], gotResult)
	}

	return nil
}

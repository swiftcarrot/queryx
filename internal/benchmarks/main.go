package main

import (
	"flag"
	"fmt"
	"github.com/swiftcarrot/queryx/internal/benchmarks/go-queryx"
	"github.com/swiftcarrot/queryx/internal/benchmarks/go-queryx/helper"
	"os"
	"runtime"
	"sort"
	"text/tabwriter"
)

var adapter string

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.StringVar(&adapter, "adapter", "postgresql", "-orm=postgresql")
	flag.Parse()
	helper.Errors = make(map[string]map[string]string, 0)
	helper.Errors["queryx"] = make(map[string]string, 0)
	runBenchmarks(adapter)
}

func runBenchmarks(adapter string) {
	table := new(tabwriter.Writer)
	table.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)
	reports := make(map[string]helper.BenchmarkReport, 0)
	err := helper.RunBenchmarks(adapter, go_queryx.CreateQueryx(), reports)
	if err != nil {
		panic(fmt.Sprintf("An error occured while running the benchmarks: %v", err))
	}
	for _, v := range reports {
		sort.Sort(v)
	}
	_, _ = fmt.Fprintf(table, "Reports:\n\n")
	i := 1
	for method, report := range reports {
		_, _ = fmt.Fprintf(table, "%s\n", method)
		for _, result := range report {
			if result.ErrorMsg == "" {
				_, _ = fmt.Fprintf(table, "%s:\t%d\t%d ns/op\t%d B/op\t%d allocs/op\n", result.Name, result.N, result.NsPerOp, result.MemBytes, result.MemAllocs)
			} else {
				_, _ = fmt.Fprintf(table, "%s:\t%s\n", result.Name, result.ErrorMsg)
			}
		}

		if i != len(reports) {
			_, _ = fmt.Fprintf(table, "\n")
		}
		i++
	}
	_ = table.Flush()
}

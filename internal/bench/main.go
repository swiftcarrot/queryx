package main

import (
	"fmt"
	"github.com/swiftcarrot/queryx/internal/bench/go-queryx"
	"github.com/swiftcarrot/queryx/internal/bench/go-queryx/helper"
	"os"
	"runtime"
	"sort"
	"text/tabwriter"

	_ "github.com/lib/pq"
)

// VERSION constant
const VERSION = "v1.0.2"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	helper.Errors = make(map[string]map[string]string, 0)
	helper.Errors["queryx"] = make(map[string]string, 0)
	runBenchmarks()
}

func runBenchmarks() {
	// Run benchmarks
	benchmarks := map[string]helper.ORMInterface{
		"queryx": go_queryx.CreateQueryx(),
	}

	table := new(tabwriter.Writer)
	table.Init(os.Stdout, 0, 8, 2, '\t', tabwriter.AlignRight)

	reports := make(map[string]helper.BenchmarkReport, 0)

	res, err := helper.RunBenchmarks(benchmarks["queryx"], reports)
	if err != nil {
		panic(fmt.Sprintf("An error occured while running the benchmarks: %v", err))
	}

	if helper.DebugMode {

		_, _ = fmt.Fprintf(table, "%s Benchmark Results:\n", "queryx")
		for _, result := range res.Results {
			if result.ErrorMsg == "" {
				_, _ = fmt.Fprintf(table, "%s:\t%d\t%d ns/op\t%d B/op\t%d allocs/op\n", result.Method, result.N, result.NsPerOp, result.MemBytes, result.MemAllocs)
			} else {
				_, _ = fmt.Fprintf(table, "%s:\t%s\n", result.Method, result.ErrorMsg)
			}
		}
		_ = table.Flush()
	}

	// Sort results
	for _, v := range reports {
		sort.Sort(v)
	}

	// Print final reports
	if helper.DebugMode {
		_, _ = fmt.Fprint(table, "\n")
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

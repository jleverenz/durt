package main

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func displaySortedResults(nodes []*Node) {
	t := buildTable(nodes)
	t.SetOutputMirror(os.Stdout)
	t.Render()
}

func buildTable(nodes []*Node) table.Writer {
	t := table.NewWriter()

	t.AppendHeader(table.Row{"Path", "Files", "Bytes", "Pct"})

	// calc totals
	var reportTotalFiles int
	var reportTotalBytes int
	longestNonWrap := 0
	for _, value := range nodes {
		reportTotalFiles += value.ancestorCount
		reportTotalBytes += value.sumBytes
		pathLength := len(value.path)
		if pathLength < 80 && pathLength > longestNonWrap {
			longestNonWrap = pathLength
		}
	}

	isHead := globalOpts.head && len(nodes) > 20

	if isHead {
		nodes = nodes[0:20]
	}

	for _, value := range nodes {
		pct := float64(value.sumBytes) / float64(reportTotalBytes) * 100
		t.AppendRows([]table.Row{
			{value.path, value.ancestorCount, ByteSize(value.sumBytes), pct},
		})
	}

	if isHead {
		t.AppendRows([]table.Row{{"...", "...", "...", "..."}})
	}

	t.AppendFooter(table.Row{"TOTALS", reportTotalFiles, ByteSize(reportTotalBytes), "100.0 %"})

	nameTransformer := text.Transformer(func(val interface{}) string {
		return text.WrapHard(val.(string), longestNonWrap)
	})

	percentTransformer := text.Transformer(func(val interface{}) string {
		num, ok := val.(float64)

		if !ok {
			return val.(string)
		}
		if num < 1 {
			return fmt.Sprintf("< 1.0 %%")
		} else {
			return fmt.Sprintf("%.01f %%", num)
		}
	})

	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Number:      1,
			Transformer: nameTransformer,
		},
		{
			Number: 2,
			VAlign: text.VAlignBottom,
		},
		{
			Number: 3,
			VAlign: text.VAlignBottom,
		},
		{
			Number:      4,
			Align:       text.AlignRight,
			VAlign:      text.VAlignBottom,
			Transformer: percentTransformer,
		},
	})

	return t
}

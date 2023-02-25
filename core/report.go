package core

import (
	"fmt"
	"os"
	"sort"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type ByBytes []PathResult

func (a ByBytes) Len() int           { return len(a) }
func (a ByBytes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByBytes) Less(i, j int) bool { return a[j].TotalBytes < a[i].TotalBytes } // Reverse

func displaySortedResults(pathResults []PathResult) {
	sort.Sort(ByBytes(pathResults))

	t := buildTable(pathResults)
	t.SetOutputMirror(os.Stdout)
	t.Render()
}

func buildTable(pathResults []PathResult) table.Writer {
	t := table.NewWriter()

	t.AppendHeader(table.Row{"Path", "Files", "Bytes", "Pct"})

	// calc totals
	var reportTotalFiles int
	var reportTotalBytes int
	longestNonWrap := 0
	for _, value := range pathResults {
		reportTotalFiles += value.TotalFiles
		reportTotalBytes += value.TotalBytes
		pathLength := len(value.Path)
		if pathLength < 80 && pathLength > longestNonWrap {
			longestNonWrap = pathLength
		}
	}

	isHead := GlobalOpts.Head && len(pathResults) > 20

	if isHead {
		pathResults = pathResults[0:20]
	}

	for _, value := range pathResults {
		pct := float64(value.TotalBytes) / float64(reportTotalBytes) * 100
		t.AppendRows([]table.Row{
			{value.Path, value.TotalFiles, ByteSize(value.TotalBytes), pct},
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

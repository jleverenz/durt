package core

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/stretchr/testify/assert"
)

func renderToLines(tw table.Writer) []string {
	return strings.Split(tw.Render(), "\n")
}

func TestEmptyReport(t *testing.T) {
	nodes := []*Node{}
	lines := renderToLines(buildTable(nodes))
	assert.Equal(t, 6, len(lines))
}

func TestUntruncatedHeadReport(t *testing.T) {
	nodes := []*Node{}

	for i := 0; i < 10; i++ {
		node := Node{path: fmt.Sprintf("./somedir-%v", i+1)}
		node.sumBytes = 1000
		node.ancestorCount = 100
		nodes = append(nodes, &node)
	}

	GlobalOpts.Head = true
	lines := renderToLines(buildTable(nodes))
	assert.Equal(t, 16, len(lines)) // 10 rows + header/footer/border
}

func TestExactlyTruncatedHeadReport(t *testing.T) {
	nodes := []*Node{}

	for i := 0; i < 20; i++ {
		node := Node{path: fmt.Sprintf("./somedir-%v", i+1)}
		node.sumBytes = 1000
		node.ancestorCount = 100
		nodes = append(nodes, &node)
	}

	GlobalOpts.Head = true
	lines := renderToLines(buildTable(nodes))
	assert.Equal(t, 26, len(lines)) // 20 rows + header/footer/border
}

func TestTruncatedHeadReport(t *testing.T) {
	nodes := []*Node{}

	for i := 0; i < 30; i++ {
		node := Node{path: fmt.Sprintf("./somedir-%v", i+1)}
		node.sumBytes = 1000
		node.ancestorCount = 100
		nodes = append(nodes, &node)
	}

	GlobalOpts.Head = true
	lines := renderToLines(buildTable(nodes))
	assert.Equal(t, 27, len(lines)) // 20 rows + ... + header/footer/border
}

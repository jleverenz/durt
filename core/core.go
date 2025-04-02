package core

import (
	"io/fs"
	"regexp"
)

type PathStat struct {
	Path string
	Stat *fs.FileInfo
}

type ProgramOptions struct {
	Head       bool
	Exclusions []*regexp.Regexp
	MyVal      string
	Strategy   string
}

type PathResult struct {
	Path       string
	TotalBytes int
	TotalFiles int
}

var GlobalOpts ProgramOptions

func Run(pathStats []PathStat) {
	var pathResults []PathResult

	switch GlobalOpts.Strategy {
	case "walk":
		pathResults = WalkStrategy.Execute(pathStats)
	case "shell":
		pathResults = ShellStrategy.Execute(pathStats)
	default:
		panic("no strategy")
	}

	displaySortedResults(pathResults)
}

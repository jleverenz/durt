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
}

type PathResult struct {
	Path       string
	TotalBytes int
	TotalFiles int
}

var GlobalOpts ProgramOptions

func Run(pathStats []PathStat) {
	pathResults := WalkStrategy.Execute(pathStats)
	displaySortedResults(pathResults)
}

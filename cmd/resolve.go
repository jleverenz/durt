package cmd

import (
	"os"

	"github.com/jleverenz/durt/core"
)

// Resolves command line arguments args
// into an array of core.PathStat structs used as
// input for the core processing package.
func ResolveArgs(args []string) []core.PathStat {
	resolved := []core.PathStat{{Path: "."}}

	if len(args) > 0 {
		resolved = []core.PathStat{}
		for _, path := range args {
			resolved = append(resolved, core.PathStat{Path: path})
		}
	}

	if len(resolved) == 1 {
		stat, _ := os.Stat(resolved[0].Path)
		resolved[0].Stat = &stat

		if stat.IsDir() {
			os.Chdir(resolved[0].Path)
			entries, _ := os.ReadDir(".")

			// expanding := args[0]
			resolved = []core.PathStat{}
			for _, entry := range entries {
				resolved = append(resolved, core.PathStat{Path: entry.Name()})
			}
		}
	}

	return resolved
}

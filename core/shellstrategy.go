package core

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
)

var ShellStrategy = struct {
	Execute func([]PathStat) []PathResult
}{Execute: executeShellStrategy}

func executeShellStrategy(pathStats []PathStat) []PathResult {
	results := []PathResult{}

	for _, pathStat := range pathStats {
		results = append(results, PathResult{
			Path:       pathStat.Path,
			TotalBytes: calcDiskUtilization(pathStat.Path),
			TotalFiles: calcFileCount(pathStat.Path),
		})
	}

	return results
}

func calcDiskUtilization(path string) int {
	out, err := exec.Command("du", "-sk", path).Output()

	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(`\d+`)
	num, err := strconv.Atoi(string(re.Find(out)))

	if err != nil {
		panic(err)
	}

	return num * 1024
}

func calcFileCount(path string) int {
	cmd := fmt.Sprintf("find %v -type f | wc -l", path)
	out, err := exec.Command("bash", "-c", cmd).Output()

	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(`\d+`)
	num, err := strconv.Atoi(string(re.Find(out)))

	if err != nil {
		panic(err)
	}

	return num
}

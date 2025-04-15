package core

import (
	"testing"
)

type StrategySubBenchmarks struct {
	Name  string
	Strat func([]PathStat) []PathResult
}

var strategies = []StrategySubBenchmarks{
	{
		Name:  "Walk",
		Strat: WalkStrategy.Execute,
	},
	{
		Name:  "Shell",
		Strat: ShellStrategy.Execute,
	},
}

func BenchmarkStrategyFull(b *testing.B) {
	pathStats := []PathStat{{Path: "../random_files", Stat: nil}}

	for _, s := range strategies {
		b.Run(s.Name, func(b *testing.B) {
			for b.Loop() {
				s.Strat(pathStats)
			}
		})
	}
}

func BenchmarkStrategySmall(b *testing.B) {
	pathStats := []PathStat{{Path: "../random_files/small", Stat: nil}}

	for _, s := range strategies {
		b.Run(s.Name, func(b *testing.B) {
			for b.Loop() {
				s.Strat(pathStats)
			}
		})
	}
}

func BenchmarkStrategyMedium(b *testing.B) {
	pathStats := []PathStat{{Path: "../random_files/medium", Stat: nil}}

	for _, s := range strategies {
		b.Run(s.Name, func(b *testing.B) {
			for b.Loop() {
				executeWalkStrategy(pathStats)
			}
		})
	}
}

func BenchmarkStrategyLarge(b *testing.B) {
	pathStats := []PathStat{{Path: "../random_files/large", Stat: nil}}

	for _, s := range strategies {
		b.Run(s.Name, func(b *testing.B) {
			for b.Loop() {
				executeWalkStrategy(pathStats)
			}
		})
	}
}

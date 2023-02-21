package main

import "fmt"

type ByteSize float64

const (
	_           = iota // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
)

func (b ByteSize) String() string {
	switch {
	case b >= GB:
		return fmt.Sprintf("%.2f G", b/GB)
	case b >= MB:
		return fmt.Sprintf("%.2f M", b/MB)
	case b >= KB:
		return fmt.Sprintf("%.2f K", b/KB)
	}
	return fmt.Sprintf("%.0f B", b)
}

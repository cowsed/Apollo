package main

import (
	"fmt"
	"strings"
)

func percentageBar(amt float64, width int) string {
	s := "<"
	pc := fmt.Sprintf("%d", int(amt*100))
	avail := width - 2 - len(pc)
	//middle := avail / 2
	large := int(amt * float64(avail))
	large = min(large, avail)
	if large > 0 {
		s += strings.Repeat("=", large)
	}
	s += pc
	if avail-large > 0 {
		s += strings.Repeat("-", avail-large)
	}
	return s + ">"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

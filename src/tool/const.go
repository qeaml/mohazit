package tool

import (
	_ "embed"
	"strconv"
	"strings"
)

var (
	//go:embed iteration
	IterationRaw string
	Iteration    int

	//go:embed version
	Version string
)

func init() {
	Version = strings.TrimSpace(Version)
	it, err := strconv.Atoi(strings.TrimSpace(IterationRaw))
	if err != nil {
		Iteration = 0
	} else {
		Iteration = it
	}
}

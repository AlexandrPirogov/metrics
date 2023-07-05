package main

import (
	m "memtracker/internal/lint/multichecker"

	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	StaticCheckRules := m.New()
	multichecker.Main(StaticCheckRules...)
}

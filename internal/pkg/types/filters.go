package types

import (
	"fmt"
)

type Filters struct {
	Patterns []string // patterns
	Before   int      // lines before
	After    int      // lines after
}

// singleton
var filters *Filters = nil

func GetFilters() *Filters {
	if filters == nil {
		filters = new(Filters)
	}

	return filters
}

func (f *Filters) String() string {
	return fmt.Sprintf("%v", f.Patterns)
}

func (f *Filters) Type() string {
	return "strings"
}

// Set global filter
func (f *Filters) Set(p string) error {
	f.Patterns = append(f.Patterns, p)

	return nil
}

// Pop global filter
func (f *Filters) Pop() {
	f.Patterns = f.Patterns[:len(f.Patterns)-1]
}

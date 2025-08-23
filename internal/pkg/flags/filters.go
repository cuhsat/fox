package flags

import (
	"fmt"
)

type Filters struct {
	Patterns []string // patterns
	Context  int      // lines context
	Before   int      // lines before
	After    int      // lines after
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

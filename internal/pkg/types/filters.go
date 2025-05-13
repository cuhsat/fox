package types

import (
	"fmt"
)

type Filters []string

// singleton
var filters *Filters = nil

func GetFilters() *Filters {
	if filters == nil {
		filters = new(Filters)
	}

	return filters
}

func (f *Filters) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *Filters) Type() string {
	return "strings"
}

func (f *Filters) Set(p string) error {
	*f = append(*f, p)

	return nil
}

func (f *Filters) Pop() {
	*f = (*f)[:len(*f)-1]
}

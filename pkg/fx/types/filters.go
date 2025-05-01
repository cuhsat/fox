package types

import (
	"fmt"
)

type filters []string

// singleton
var _filters *filters = nil

func Filters() *filters {
	if _filters == nil {
		_filters = new(filters)
	}

	return _filters
}

func (f *filters) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *filters) Type() string {
	return "strings"
}

func (f *filters) Set(p string) error {
	*f = append(*f, p)

	return nil
}

func (f *filters) Pop() {
	*f = (*f)[:len(*f)-1]
}

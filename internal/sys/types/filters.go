package types

import (
    "strings"
)

type Filters []string

// singleton instance
var filters *Filters = nil

func GetFilters() *Filters {
    if filters == nil {
        filters = new(Filters);
    }

    return filters;
}

func (f *Filters) String() string {
    return strings.Join(*f, " > ")
}

func (f *Filters) Set(s string) error {
    *f = append(*f, s)

    return nil
}

func (f *Filters) Pop() {
    *f = (*f)[:len(*f)-1]
}

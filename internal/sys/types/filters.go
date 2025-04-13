package types

import (
    "fmt"
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

// func (f *Filters) Push(s string) {
//     *f = append(*f, s)
// }

// func (f *Filters) Pop() (s string) {
//     s = *f[len(*f)-1]
//     *f = *f[:len(*f)-1]

//     return
// }

func (f *Filters) String() string {
    return fmt.Sprintf("%v", *f)
}

func (f *Filters) Set(s string) error {
    *f = append(*f, s)
    return nil
}

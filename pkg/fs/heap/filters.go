package heap

import (
    "fmt"
)

type Filters []string

func (f *Filters) String() string {
    return fmt.Sprintf("%v", *f)
}

func (f *Filters) Set(s string) error {
    *f = append(*f, s)
    return nil
}

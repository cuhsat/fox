package data

import (
    "fmt"
)

type Filter []string

func (f *Filter) String() string {
    return fmt.Sprintf("%v", *f)
}

func (f *Filter) Set(v string) error {
    *f = append(*f, v)

    return nil
}

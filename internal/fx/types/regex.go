package types

import (
    "regexp"

    "github.com/cuhsat/fx/internal/fx"
)

func Regex(s string) (bool, *regexp.Regexp) {
    if len(s) < 2 {
        return false, nil
    }

    if s[0] != '/' || s[len(s)-1] != '/' {
        return false, nil
    }

    expr := s[1:len(s)-1]

    re, err := regexp.Compile(expr)

    if err != nil {
        fx.Error(err)
    }

    return true, re
}

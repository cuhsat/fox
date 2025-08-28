package text

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/cuhsat/fox/internal/pkg/flags"
)

var Patterns = []Pattern{
	{
		Name: "ipv4",
		Mask: regexp.MustCompile("\\b(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\\b"),
	},
	{
		Name: "ipv6",
		Mask: regexp.MustCompile("(([a-fA-F0-9]{1,4}:){7,7}[a-fA-F0-9]{1,4}|([a-fA-F0-9]{1,4}:){1,7}:|([a-fA-F0-9]{1,4}:){1,6}:[a-fA-F0-9]{1,4}|([a-fA-F0-9]{1,4}:){1,5}(:[a-fA-F0-9]{1,4}){1,2}|([a-fA-F0-9]{1,4}:){1,4}(:[a-fA-F0-9]{1,4}){1,3}|([a-fA-F0-9]{1,4}:){1,3}(:[a-fA-F0-9]{1,4}){1,4}|([a-fA-F0-9]{1,4}:){1,2}(:[a-fA-F0-9]{1,4}){1,5}|[a-fA-F0-9]{1,4}:((:[a-fA-F0-9]{1,4}){1,6})|:((:[a-fA-F0-9]{1,4}){1,7}|:)|fe80:(:[a-fA-F0-9]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([a-fA-F0-9]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))"),
	},
	{
		Name: "mac",
		Mask: regexp.MustCompile("([a-fA-F0-9]{2}[:-]){5}([a-fA-F0-9]{2})"),
	},
	{
		Name: "mail",
		Mask: regexp.MustCompile("\\b[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,6}\\b"),
	},
	{
		Name: "url",
		Mask: regexp.MustCompile("[-a-zA-Z0-9@:%._+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b[-a-zA-Z0-9()@:%_+.~#?&/=]*"),
	},
	{
		Name: "uuid",
		Mask: regexp.MustCompile("[a-fA-F0-9]{8}(?:-[a-fA-F0-9]{4}){3}-[a-fA-F0-9]{12}"),
	},
	{
		Name: "data",
		Mask: regexp.MustCompile(".*"),
	},
}

type String struct {
	Off int
	Str string
}

type Pattern struct {
	Name string
	Mask *regexp.Regexp
}

func Carve(in <-chan byte, out chan<- String, n, m int) {
	var rs []rune
	var off int

	flush := func() {
		if len(rs) >= n && len(rs) <= m {
			o := max(off-(len(rs)+1), 0)
			s := string(rs)

			if len(strings.TrimSpace(s)) > 0 {
				out <- String{o, s}
			}
		}

		rs = rs[:0]
	}

	defer close(out)
	defer flush()

	flg := flags.Get().Strings
	buf := make([]byte, 4)

	for b := range in {
		buf[0] = b
		off++

		if flg.Ascii {
			if b >= MinASCII && b <= MaxASCII {
				rs = append(rs, rune(b))
			} else {
				flush()
			}
		} else {
			l := 1
			k := 1

			if b&0x80 == 0 {
				k = 1
			} else if b&0xE0 == 0xC0 {
				k = 2
			} else if b&0xF0 == 0xE0 {
				k = 3
			} else if b&0xF8 == 0xF0 {
				k = 4
			}

			if k > 1 {
				for i := 1; i < k; i++ {
					off++

					if b, ok := <-in; ok {
						buf[i] = b
					} else {
						break
					}

					l++
				}
			}

			r, _ := utf8.DecodeRune(buf[:l])

			if r != utf8.RuneError && unicode.IsPrint(r) {
				rs = append(rs, r)
			} else {
				flush()
			}
		}
	}
}

func Match(in <-chan String, out chan<- String, i bool, re *regexp.Regexp) {
	defer close(out)

	var ps []Pattern

	// search for pattern
	if re != nil {
		ps = []Pattern{{"ioc", re}}
	}

	// search for built-in IoCs
	if i {
		ps = append(ps, Patterns...)
	}

	for s := range in {
		for _, p := range ps {
			if p.Mask.MatchString(s.Str) {
				out <- String{s.Off, fmt.Sprintf("%-4s  %s", p.Name, s.Str)}
				break
			}
		}
	}
}

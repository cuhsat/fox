package text

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/cuhsat/fox/internal/pkg/flags"
)

// https://github.com/AndrewRathbun/DFIRRegex
var patterns = []struct {
	ioc string
	re  *regexp.Regexp
}{
	{
		ioc: "IPv4",
		re:  regexp.MustCompile("\\b(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\\b"),
	},
	{
		ioc: "IPv6",
		re:  regexp.MustCompile("(([a-fA-F0-9]{1,4}:){7,7}[a-fA-F0-9]{1,4}|([a-fA-F0-9]{1,4}:){1,7}:|([a-fA-F0-9]{1,4}:){1,6}:[a-fA-F0-9]{1,4}|([a-fA-F0-9]{1,4}:){1,5}(:[a-fA-F0-9]{1,4}){1,2}|([a-fA-F0-9]{1,4}:){1,4}(:[a-fA-F0-9]{1,4}){1,3}|([a-fA-F0-9]{1,4}:){1,3}(:[a-fA-F0-9]{1,4}){1,4}|([a-fA-F0-9]{1,4}:){1,2}(:[a-fA-F0-9]{1,4}){1,5}|[a-fA-F0-9]{1,4}:((:[a-fA-F0-9]{1,4}){1,6})|:((:[a-fA-F0-9]{1,4}){1,7}|:)|fe80:(:[a-fA-F0-9]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([a-fA-F0-9]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))"),
	},
	{
		ioc: "MAC",
		re:  regexp.MustCompile("([a-fA-F0-9]{2}[:-]){5}([a-fA-F0-9]{2})"),
	},
	// {
	// 	ioc: "URL",
	// 	re:  regexp.MustCompile("[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b(?:[-a-zA-Z0-9()@:%_\\+.~#?&//=]*)"),
	// },
	// {
	// 	ioc: "Mail",
	// 	re:  regexp.MustCompile("(([a-zA-Z0-9_\\-\\.]+)@((\\[[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.)|(([a-zA-Z0-9\\-]+\\.)+))([a-zA-Z]{2,4}|[0-9]{1,3})(\\]?)(\\s*;\\s*|\\s*$))*"),
	// },
	{
		ioc: "....",
		re:  regexp.MustCompile(".*"),
	},
}

type String struct {
	Off int
	Str string
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

func Match(in <-chan String, out chan<- String) {
	defer close(out)

	for s := range in {
		for _, p := range patterns {
			if p.re.MatchString(s.Str) {
				out <- String{s.Off, fmt.Sprintf("%-4s  %s", p.ioc, s.Str)}
				break
			}
		}
	}
}

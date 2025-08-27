package text

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/cuhsat/fox/internal/pkg/flags"
)

// Patterns based on https://github.com/AndrewRathbun/DFIRRegex
var Patterns = map[string]string{
	"Hex": "#?([a-fA-F0-9]{6}|[a-fA-F0-9]{3})",
	//"Base64":   "(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?",
	//"Filename": "[^\\\\\\/:*?\"<>|\\r\\n]+$",
	//"URL":           "(https?:\\/\\/)?(www\\.)?[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b([-a-zA-Z0-9()!@:%_\\+.~#?&\\/\\/=]*)",
	"IPv4": "\\b(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\\b",
	//"IPv6":          "(([a-fA-F0-9]{1,4}:){7,7}[a-fA-F0-9]{1,4}|([a-fA-F0-9]{1,4}:){1,7}:|([a-fA-F0-9]{1,4}:){1,6}:[a-fA-F0-9]{1,4}|([a-fA-F0-9]{1,4}:){1,5}(:[a-fA-F0-9]{1,4}){1,2}|([a-fA-F0-9]{1,4}:){1,4}(:[a-fA-F0-9]{1,4}){1,3}|([a-fA-F0-9]{1,4}:){1,3}(:[a-fA-F0-9]{1,4}){1,4}|([a-fA-F0-9]{1,4}:){1,2}(:[a-fA-F0-9]{1,4}){1,5}|[a-fA-F0-9]{1,4}:((:[a-fA-F0-9]{1,4}){1,6})|:((:[a-fA-F0-9]{1,4}){1,7}|:)|fe80:(:[a-fA-F0-9]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([a-fA-F0-9]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))",
	"MAC":           "([a-fA-F0-9]{2}[:-]){5}([a-fA-F0-9]{2})",
	"Hash (MD5)":    "[a-fA-F0-9]{32}",
	"Hash (SHA1)":   "[a-fA-F0-9]{40}",
	"Hash (SHA256)": "[a-fA-F0-9]{64}",
	"Hash (SHA512)": "[a-fA-F0-9]{128}",
	//"Password":      "(?=.*?[A-Z])(?=.*?[a-z])(?=.*?[0-9])(?=.*?[#?!@$ %^&*-]).{8,}",
	//"Email":         "(([a-zA-Z0-9_\\-\\.]+)@((\\[[0-9]{1,3}\\.[0-9]{1,3}\\.[0-9]{1,3}\\.)|(([a-zA-Z0-9\\-]+\\.)+))([a-zA-Z]{2,4}|[0-9]{1,3})(\\]?)(\\s*;\\s*|\\s*$))*",
	//"Phone":         "(\\+\\d{1,2}\\s)?\\(?\\d{3}\\)?[\\s.-]?\\d{3}[\\s.-]?\\d{4}",
	//"Credit Card":   "(4[0-9]{12}(?:[0-9]{3})?$)|(^(?:5[1-5][0-9]{2}|222[1-9]|22[3-9][0-9]|2[3-6][0-9]{2}|27[01][0-9]|2720)[0-9]{12}$)|(3[47][0-9]{13})|(^3(?:0[0-5]|[68][0-9])[0-9]{11}$)|(^6(?:011|5[0-9]{2})[0-9]{12}$)|(^(?:2131|1800|35\\d{3})\\d{11}$)",
	//"Social Security": "(?!0{3})(?!6{3})[0-8]\\d{2}-(?!0{2})\\d{2}-(?!0{4})\\d{4}",
}

type String struct {
	Off int
	Str string
}

func Recognize(ch <-chan String, str chan<- String) {
	defer close(str)

	ps := make(map[string]*regexp.Regexp, len(Patterns))

	// compile patterns
	for k, v := range Patterns {
		ps[k] = regexp.MustCompile(v)
	}

	var ok bool

	for s := range ch {
		ok = false

		for k, v := range ps {
			if v.MatchString(s.Str) {
				str <- String{s.Off, fmt.Sprintf("%s (%s)", s.Str, k)}
				ok = true
				break
			}
		}

		if !ok {
			str <- s
		}
	}
}

func Carve(ch <-chan byte, str chan<- String, n, m int) {
	var rs []rune
	var off int

	flush := func() {
		if len(rs) >= n && len(rs) <= m {
			o := max(off-(len(rs)+1), 0)
			s := string(rs)

			if len(strings.TrimSpace(s)) > 0 {
				str <- String{o, s}
			}
		}

		rs = rs[:0]
	}

	defer close(str)
	defer flush()

	flg := flags.Get().Strings
	buf := make([]byte, 4)

	for b := range ch {
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

					if b, ok := <-ch; ok {
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

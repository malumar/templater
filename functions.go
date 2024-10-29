package templater

import (
	"fmt"
	"github.com/malumar/domaintools"
	"github.com/malumar/strutils"
	"golang.org/x/net/html"
	"html/template"
	"log"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// Zamiana  "   Text   More here     "
// na  "Text More here"
func ReduceWhiteSpaces(input string) string {
	final := re_leadclose_whtsp.ReplaceAllString(input, "")
	final = re_inside_whtsp.ReplaceAllString(final, " ")
	return final
}

// liczba monga
func Plural(s string) string {
	if len(s) == 0 {
		return ""
	}

	if s[len(s)-1:] == "s" {
		return s + "es"
	} else {
		return s + "s"
	}
}

func IsEmpty(v interface{}) bool {
	if v == nil || v == "" || v == false {
		return true
	}
	tv := reflect.TypeOf(v)

	vv := reflect.ValueOf(v)
	switch tv.Kind() {
	case reflect.Map, reflect.Slice:
		return vv.Len() == 0
	case reflect.String:
		return vv.String() == ""
	case reflect.Int:
		return vv.Int() == 0
	case reflect.Uint:
		return vv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return vv.Float() == 0
	case reflect.Bool:
		return vv.Bool() == false
		break
	}
	return false
}

func IsNotEmpty(v interface{}) bool {
	return !IsEmpty(v)
}

func FirstLower(s string) string {
	if len(s) == 0 {
		return ""
	}

	return strings.ToLower(s[:1]) + s[1:]
}
func FirstUpper(s string) string {
	if len(s) == 0 {
		return ""
	}

	return strings.ToUpper(s[:1]) + s[1:]
}

func IsLast(x int, a interface{}) bool {
	return x == reflect.ValueOf(a).Len()-1
}
func init() {

	functions = map[string]interface{}{
		"set_html": func(renderArgs map[string]interface{}, key string, value interface{}) template.JS {
			renderArgs[key] = value
			return template.JS("")
		},
		"append_html": func(renderArgs map[string]interface{}, key string, value interface{}) template.JS {
			if renderArgs[key] == nil {
				renderArgs[key] = []interface{}{value}
			} else {
				renderArgs[key] = append(renderArgs[key].([]interface{}), value)
			}
			return template.JS("")
		},
		"firstOf": func(args ...interface{}) interface{} {
			for _, val := range args {
				switch val.(type) {
				case nil:
					continue
				case string:
					if val == "" {
						continue
					}
					return val
				default:
					return val
				}
			}
			return nil
		},
		"printIfNotLast": func(index int, len int, what string) string {
			if index+1 < len {
				return fmt.Sprintf("%d z %d", index+1, len)
			}
			return ""
		},
		// Pads the given string with space up to the given width.
		"pad_html": func(str string, width int) template.HTML {
			if len(str) >= width {
				return template.HTML(html.EscapeString(str))
			}
			return template.HTML(html.EscapeString(str) + strings.Repeat("&nbsp;", width-len(str)))
		},
		// Pads the given string with space up to the given width.
		"pad": func(str string, width int) string {
			if len(str) >= width {
				return str
			}
			return str + strings.Repeat(" ", width-len(str))
		},
		// Pluralize, a helper for pluralizing words to correspond to data of dynamic length.
		// items - a slice of items, or an integer indicating how many items there are.
		// pluralOverrides - optional arguments specifying the output in the
		//     singular and plural cases.  by default "" and "s"
		"pluralize": func(items interface{}, pluralOverrides ...string) string {
			singular, plural := "", "s"
			if len(pluralOverrides) >= 1 {
				singular = pluralOverrides[0]
				if len(pluralOverrides) == 2 {
					plural = pluralOverrides[1]
				}
			}

			switch v := reflect.ValueOf(items); v.Kind() {
			case reflect.Int:
				if items.(int) != 1 {
					return plural
				}
			case reflect.Slice:
				if v.Len() != 1 {
					return plural
				}
			default:
				log.Println("pluralize: unexpected type: ", v)
			}
			return singular
		},
		// 2006-01-02 15:04:05
		"year": func(date time.Time) string {
			return date.Format("2006")
		},
		"yearZero": func(date time.Time) string {
			return date.Format("06")
		},
		"monthZero": func(date time.Time) string {
			return date.Format("01")
		},
		"dayZero": func(date time.Time) string {
			return date.Format("02")
		},
		"hourZero": func(date time.Time) string {
			return date.Format("15")
		},
		"minuteZero": func(date time.Time) string {
			return date.Format("04")
		},
		"secZero": func(date time.Time) string {
			return date.Format("05")
		},
		"even": func(a int) bool { return (a % 2) == 0 },
		// zwiększenie licznika
		"inc": func(v int) int {
			return v + 1
		},
		// zmniejszenie licznika
		"dec": func(v int) int {
			return v - 1
		},
		// liczba monga dla słów angielskich
		"plural_en": Plural,

		// Tekst -> tekst
		"firstLower": FirstLower,
		"firstUpper": FirstUpper,

		// Tekst -> t
		"firstCharLower": func(s string) string {
			if s == "" {
				return ""
			}
			return strings.ToLower(s[:1])
		},
		// tekst -> T
		"firstCharUpper": func(s string) string {
			if s == "" {
				return ""
			}
			return strings.ToLower(s[:1])
		},

		"bool2str": func(val bool, ifTrue, ifFalse string) string {
			if val {
				return ifTrue
			} else {
				return ifFalse
			}
		},

		"toUpper":   strings.ToUpper,
		"toLower":   strings.ToLower,
		"trimSpace": strings.TrimSpace,
		"replaceNewLineToSpace": func(s string) string {
			return strings.Replace(s, "\n", " ", -1)
		},
		"removeWhiteSpace": func(s string) string {
			return strings.Replace(strings.Replace(strings.Replace(s, "\n", "", -1), "	", "", -1), " ", "", -1)
		},
		// Znak `
		"esc": func() string {
			return "`"
		},
		"replaceString": func(old string, new string, searchIn string) string {
			return strings.Replace(searchIn, old, new, -1)
		},
		"isNotLast": func(x int, a interface{}) bool {
			return !IsLast(x, a)
		},

		// czy wartość jest pusta
		"isEmpty": IsEmpty,
		// czy wartość jest pusta
		"isNotEmpty": IsNotEmpty,
		// zmienna|islast $x czy item $x jest ju ostatni
		"isLast": IsLast,

		"concat": func(sep string, args ...string) string {
			return strings.Join(args, sep)
		},

		// zmienna|addsfx "_sfx" --> zmienna_sfx
		"addSfx": func(sfx string, s string) string {
			return s + sfx
		},

		// zmienna|addsfx "pfx_" --> pfx_zmienna
		"addPfx": func(pfx string, s string) string {
			return pfx + s
		},

		// zamiana:
		// zmienna -> zmienna
		// zmiennaInna -> zmienna_Inna
		// ZmiennnaInnaDruga -> Zmienna_Inna_Drug
		"underScoreCase": strutils.UnderscoreCase,
		"flatCase":       strutils.FlatCase,
		// zamiana:
		// zmienna -> ZMIENNA
		// zmiennaInna -> ZMIENNA_INNA
		// ZmiennnaInnaDruga -> ZMIENNA_INNA_DRUGA
		"underScoreCaseUpper": strutils.UnderscoreCaseUpper,
		"snakeCaseFirstLower": strutils.SnakeCaseFirstLower,
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"toSafeAsciiDomainName": domaintools.SafeAsciiDomainName,
		"reduceWhiteSpaces":     ReduceWhiteSpaces,
	}


}

var re_leadclose_whtsp *regexp.Regexp = regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
var re_inside_whtsp *regexp.Regexp = regexp.MustCompile(`[\s\p{Zs}]{2,}`)

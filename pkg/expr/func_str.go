package expr

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/spf13/cast"

	"github.com/PaesslerAG/gval"
	valid "github.com/asaskevich/govalidator"
)

func StringFunctions() []gval.Language {
	return []gval.Language{
		gval.Function("trim", strings.TrimSpace),
		gval.Function("trimLeft", strings.TrimLeft),
		gval.Function("trimRight", strings.TrimRight),
		gval.Function("toLower", strings.ToLower),
		gval.Function("toUpper", strings.ToUpper),
		gval.Function("shortest", shortest),
		gval.Function("longest", longest),
		gval.Function("format", fmt.Sprintf),
		gval.Function("title", title),
		gval.Function("untitle", untitle),
		gvalFunc("repeat", strings.Repeat),
		gvalFunc("replace", strings.Replace),
		gval.Function("isUrl", valid.IsURL),
		gval.Function("isEmail", valid.IsEmail),
		gval.Function("split", strings.Split),
		gval.Function("join", join),
		gval.Function("hasSubstring", hasSubstring),
		gvalFunc("substring", substring),
		gval.Function("hasPrefix", strings.HasPrefix),
		gval.Function("hasSuffix", strings.HasSuffix),
		gvalFunc("shorten", shorten),
		gval.Function("camelize", camelize),
		gval.Function("snakify", snakify),
		gval.Function("match", match),
	}
}

// gvalFunc cast any nth number of float64 param to int
func gvalFunc(name string, fn interface{}) gval.Language {
	return gval.Function(name, func(params ...interface{}) (out interface{}) {
		in := make([]reflect.Value, len(params))
		for i, param := range params {
			if reflect.TypeOf(param).Kind() == reflect.Float64 {
				param = cast.ToInt(param)
			}
			in[i] = reflect.ValueOf(param)
		}
		fun := reflect.ValueOf(fn)
		res := fun.Call(in)
		if len(res) > 0 {
			out = cast.ToString(reflect.ValueOf(res[0]).Interface())
		}
		return
	})
}

func shortest(f string, aa ...string) string {
	for _, s := range aa {
		if len(f) > len(s) {
			f = s
		}
	}

	return f
}

func longest(f string, aa ...string) string {
	for _, s := range aa {
		if len(f) < len(s) {
			f = s
		}
	}

	return f
}

// title works similarly as strings.ToTitle, with the expception
// of uppercasing only the first word in line
func title(s string) string {
	split := strings.Split(s, " ")
	return fmt.Sprintf("%s %s", strings.Title(split[0]), strings.Join(split[1:], " "))
}

// untitle works only on the first word in line
func untitle(s string) string {
	split := strings.Split(s, " ")
	first := strings.ToLower(split[0][:1]) + split[0][1:]

	return fmt.Sprintf("%s %s", first, strings.Join(split[1:], " "))
}

func join(arr interface{}, sep string) (out string, err error) {
	if arr == nil {
		// If base is empty, nothing to do
		return "", nil
	} else if i, is := arr.([]string); is {
		// If string slice, we are good to go
		return strings.Join(i, sep), nil
	} else if i, is := arr.([]interface{}); is {
		// If slice of interfaces, we can try to cast them
		var aux []string
		aux, err = CastToStringSlice(i)
		if err != nil {
			return
		}
		return strings.Join(aux, sep), nil
	} else if arr, err = toSlice(arr); err != nil {
		return
	}

	// Make an aux string slice so the join operation can use it
	stv, is := arr.([]TypedValue)
	if !is {
		return "", errors.New("could not cast array to string array")
	}

	aux := make([]string, len(stv))
	for i, rv := range stv {
		aux[i], err = CastToString(rv)
		if err != nil {
			return
		}
	}

	return strings.Join(aux, sep), nil
}

// hasSubstring checks if a substring exists in original string
// use watchCase if need case sensitivity
func hasSubstring(s, substring string, watchCase bool) bool {
	if watchCase {
		return strings.Contains(s, substring)
	}

	return strings.Contains(strings.ToLower(s), strings.ToLower(substring))
}

// substring extracts a substring from original string
// specifying end value will not match till end of string
func substring(s string, start, end int) string {
	if end == -1 {
		end = len(s)
	}

	if start >= len(s) {
		return ""
	}

	end++
	if end > len(s) {
		end = len(s)
	}

	return s[start:end]
}

// shorten trims by whole words or only chars by
// the specified amount
func shorten(s, t string, count int) string {
	var joined string

	if t == "char" {
		if count > len(s) {
			return ""
		}

		joined = s[:count]
	} else {
		fields := strings.Fields(s)

		if len(fields) == 0 {
			return ""
		}

		joined = strings.Join(fields[:count], " ")
	}

	reg, err := regexp.Compile("[^a-zA-Z0-9]$")

	if err != nil {
		return ""
	}

	joined = reg.ReplaceAllString(joined, "")

	return joined + " …"
}

func camelize(s string) string {
	return untitle(strings.Replace(strings.Title(s), " ", "", -1))
}

func snakify(s string) string {
	return strings.ToLower(strings.Replace(strings.Title(s), " ", "_", -1))
}

func match(s string, regex interface{}) (b bool, err error) {
	var (
		r *regexp.Regexp
	)

	switch v := regex.(type) {
	case *regexp.Regexp:
		r = v
	case string:
		if r, err = regexp.Compile(v); err != nil {
			return
		}
	}

	b = r.MatchString(s)
	return
}

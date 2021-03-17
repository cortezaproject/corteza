package expr

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PaesslerAG/gval"
	valid "github.com/asaskevich/govalidator"
)

func StringFunctions() []gval.Language {
	return []gval.Language{
		gval.Function("trim", strings.TrimSpace),
		gval.Function("trimLeft", strings.TrimLeft),
		gval.Function("trimRight", strings.TrimRight),
		gval.Function("length", length),
		gval.Function("toLower", strings.ToLower),
		gval.Function("toUpper", strings.ToUpper),
		gval.Function("shortest", shortest),
		gval.Function("longest", longest),
		gval.Function("format", fmt.Sprintf),
		gval.Function("title", title),
		gval.Function("untitle", untitle),
		gval.Function("repeat", strings.Repeat),
		gval.Function("replace", strings.Replace),
		gval.Function("isUrl", valid.IsURL),
		gval.Function("isEmail", valid.IsEmail),
		gval.Function("split", strings.Split),
		gval.Function("join", strings.Join),
		gval.Function("hasSubstring", hasSubstring),
		gval.Function("substring", substring),
		gval.Function("hasPrefix", strings.HasPrefix),
		gval.Function("hasSuffix", strings.HasSuffix),
		gval.Function("shorten", shorten),
		gval.Function("camelize", camelize),
		gval.Function("snakify", snakify),
	}
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

func length(s string) int {
	return len(s)
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

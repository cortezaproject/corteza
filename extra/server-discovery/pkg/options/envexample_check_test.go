package options

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cast"
)

type documented struct {
	// what the "Default:" line says
	doc string
	// what the commented assignment line holds
	val string
}

func parseEnvExample(t *testing.T) map[string]documented {
	f, err := os.Open("../../.env.example")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var (
		out    = map[string]documented{}
		assign = regexp.MustCompile(`^# ([A-Z][A-Z0-9_<>]*)=(.*)$`)
		s      = bufio.NewScanner(f)
		doc    string
	)

	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "# Default:") {
			doc = strings.TrimPrefix(strings.TrimPrefix(line, "# Default:"), " ")
			continue
		}
		if m := assign.FindStringSubmatch(line); m != nil {
			out[m[1]] = documented{doc: doc, val: m[2]}
			doc = ""
		}
	}

	return out
}

// .env.example here is written by hand, so nothing but this stops it drifting
// away from the options it describes
func TestEnvExampleMatchesDefaults(t *testing.T) {
	es, err := ES()
	if err != nil {
		t.Fatal(err)
	}

	vs, err := VectorSearch()
	if err != nil {
		t.Fatal(err)
	}

	ex := parseEnvExample(t)

	for _, o := range []interface{}{HttpServer(), Environment(), WaitFor(), es, vs} {
		v := reflect.ValueOf(o).Elem()
		for i := 0; i < v.NumField(); i++ {
			tag := v.Type().Field(i).Tag.Get("env")
			if tag == "" {
				continue
			}

			e, ok := ex[tag]
			if !ok {
				t.Errorf("%s: not documented", tag)
				continue
			}

			// a blank assignment line means the default is described in prose
			if e.val == "" && e.doc != "" {
				continue
			}

			if e.doc != e.val {
				t.Errorf("%s: Default: says %q but the assignment line says %q", tag, e.doc, e.val)
				continue
			}

			got, want, err := compare(v.Field(i), e.val)
			if err != nil {
				t.Errorf("%s: %s", tag, err)
			} else if got != want {
				t.Errorf("%s: documents %q, actual default is %q", tag, want, got)
			}
		}
	}
}

// renders the field's default and the documented value the same way, so they
// can be compared as the loader would read them
func compare(f reflect.Value, doc string) (got, want string, err error) {
	if f.Type() == reflect.TypeOf(time.Duration(1)) {
		d, e := cast.ToDurationE(doc)
		if e != nil {
			return "", "", fmt.Errorf("%q is not a duration the loader can parse", doc)
		}
		return time.Duration(f.Int()).String(), d.String(), nil
	}

	switch f.Kind() {
	case reflect.String:
		return f.String(), doc, nil
	case reflect.Bool:
		b, e := cast.ToBoolE(doc)
		if e != nil {
			return "", "", fmt.Errorf("%q is not a bool the loader can parse", doc)
		}
		return fmt.Sprint(f.Bool()), fmt.Sprint(b), nil
	case reflect.Int:
		i, e := cast.ToIntE(doc)
		if e != nil {
			return "", "", fmt.Errorf("%q is not an int the loader can parse", doc)
		}
		return fmt.Sprint(f.Int()), fmt.Sprint(i), nil
	case reflect.Slice:
		return fmt.Sprint(f.Interface()), fmt.Sprintf("[%s]", doc), nil
	}

	return "", "", fmt.Errorf("unsupported kind %s", f.Kind())
}

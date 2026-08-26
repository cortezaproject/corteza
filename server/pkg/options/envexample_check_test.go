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

// every option struct the generator emits, in its default state
func allOpts() []interface{} {
	return []interface{}{
		DB(), HTTPClient(), HttpServer(), Rbac(), SCIM(), SMTP(), ActionLog(), Apigw(),
		Auth(), Corredor(), Environment(), Eventbus(), Federation(), Limit(), Locale(),
		Log(), Messagebus(), Monitor(), ObjectStore(), Provision(), Sentry(), Template(),
		Upgrade(), WaitFor(), Websocket(), Workflow(), Discovery(), Attachment(), Webapp(),
		Observability(), Agentic(),
	}
}

type exOpt struct {
	doc, val string
	line     int
}

func parseEnvExample(t *testing.T) map[string]exOpt {
	f, err := os.Open("../../.env.example")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var (
		out    = map[string]exOpt{}
		doc    string
		assign = regexp.MustCompile(`^# ([A-Z][A-Z0-9_]*)=(.*)$`)
		s      = bufio.NewScanner(f)
		n      int
	)
	for s.Scan() {
		n++
		line := s.Text()
		if strings.HasPrefix(line, "# Default:") {
			doc = strings.TrimPrefix(strings.TrimPrefix(line, "# Default:"), " ")
			continue
		}
		if m := assign.FindStringSubmatch(line); m != nil {
			out[m[1]] = exOpt{doc: doc, val: m[2], line: n}
			doc = ""
		}
	}
	return out
}

func TestEnvExampleMatchesDefaults(t *testing.T) {
	ex := parseEnvExample(t)
	seen := map[string]bool{}
	noted := make([]string, 0, 16)

	for _, o := range allOpts() {
		v := reflect.ValueOf(o).Elem()
		for i := 0; i < v.NumField(); i++ {
			tag := v.Type().Field(i).Tag.Get("env")
			if tag == "" {
				continue
			}
			seen[tag] = true

			e, ok := ex[tag]
			if !ok {
				t.Errorf("%s: not documented in .env.example", tag)
				continue
			}

			f := v.Field(i)

			// a Default: line that is prose rather than the assignment describes a
			// default no literal can hold - computed, or dependent on another option
			if e.doc != e.val && e.doc != "" {
				noted = append(noted, tag)
				continue
			}

			if !sameValue(t, tag, f, e.val) {
				t.Errorf("%s (.env.example:%d): documents %q, actual default is %v", tag, e.line, e.val, f.Interface())
			}
		}
	}

	t.Logf("%d of %d options document their default in prose and are not compared: %s",
		len(noted), len(ex), strings.Join(noted, " "))

	for env := range ex {
		if !seen[env] {
			t.Errorf("%s: documented in .env.example but no option carries it", env)
		}
	}
}

func sameValue(t *testing.T, tag string, f reflect.Value, doc string) bool {
	if f.Type() == reflect.TypeOf(time.Duration(1)) {
		d, err := cast.ToDurationE(doc)
		if err != nil {
			t.Errorf("%s: %q is not a duration the loader can parse", tag, doc)
			return false
		}
		return time.Duration(f.Int()) == d
	}

	switch f.Kind() {
	case reflect.String:
		return f.String() == doc
	case reflect.Bool:
		if doc == "" {
			return false
		}
		b, err := cast.ToBoolE(doc)
		if err != nil {
			t.Errorf("%s: %q is not a bool the loader can parse", tag, doc)
			return false
		}
		return f.Bool() == b
	case reflect.Int, reflect.Int64:
		i, err := cast.ToIntE(doc)
		if err != nil {
			t.Errorf("%s: %q is not an int the loader can parse", tag, doc)
			return false
		}
		return f.Int() == int64(i)
	case reflect.Float32, reflect.Float64:
		fl, err := cast.ToFloat64E(doc)
		if err != nil {
			t.Errorf("%s: %q is not a float the loader can parse", tag, doc)
			return false
		}
		return f.Float() == fl
	}

	t.Errorf("%s: unsupported kind %s", tag, f.Kind())
	return false
}

var _ = fmt.Sprint

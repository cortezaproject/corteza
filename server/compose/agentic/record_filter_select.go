package agentic

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	cmpTypes "github.com/crusttech/human/server/compose/types"
)

// selectTerm matches `stage = 'Offer'` — a bare field, an exact operator and a
// quoted literal.
//
// The leading class keeps a dotted path out: in `card.name = 'Bolt'` the field
// half is preceded by a '.', which resolveRefPaths owns and this must not
// judge. Only exact comparisons are matched — a LIKE literal is a pattern and
// is not meant to be one of the declared values.
var selectTerm = regexp.MustCompile(`(^|[^A-Za-z0-9_.])([A-Za-z_][A-Za-z0-9_]*)\s*(=|!=|<>)\s*'((?:[^'\\]|\\.)*)'`)

// checkSelectValues refuses a comparison against a value a Select field does
// not declare.
//
// A Select stores the value and shows the text, and they differ by design —
// `{"value":"offer","text":"Offer"}`. Filtering on the label is therefore the
// natural mistake, and it is valid to every layer below: the expression parses,
// the query runs, and the answer is an empty set indistinguishable from "there
// are none". An agent asked which candidates are at the Offer stage answered
// that there were none while four sat in the module.
//
// The module knows its options at query time, so the mistake is knowable here.
// It is refused rather than resolved: a label is not unique — two options may
// share a text, and one option's text may be another's value — so resolving
// would guess, and guessing is what produced the wrong answer in the first
// place. The shape of the refusal follows fieldNameError: say what failed, name
// the legal values, and point at the near match.
func checkSelectValues(mod *cmpTypes.Module, subject, expr string) error {
	if mod == nil || !strings.Contains(expr, "'") {
		return nil
	}

	quoted := literalRanges(expr)

	for _, loc := range selectTerm.FindAllStringSubmatchIndex(expr, -1) {
		// A term starting inside a string literal is part of that literal.
		if inRanges(quoted, loc[4]) {
			continue
		}

		var (
			field   = expr[loc[4]:loc[5]]
			literal = expr[loc[8]:loc[9]]
		)

		f := mod.Fields.FindByName(field)
		if f == nil || f.Kind != "Select" {
			continue
		}

		// An empty literal is "no value", which every Select can hold and no
		// option declares. A prefilter interpolation is resolved by whatever
		// renders the block, so its text is never a declared value either.
		if literal == "" || strings.Contains(literal, "${") {
			continue
		}

		opts := selectOptionTexts(f)
		if len(opts) == 0 {
			continue
		}
		if _, ok := opts[literal]; ok {
			continue
		}

		return fmt.Errorf("%s failed: %s", subject, selectValueMessage(field, literal, opts))
	}

	return nil
}

// selectOptionTexts maps each declared value to the text shown in its place.
//
// The options are held as the API delivered them — a list of strings, or of
// {value, text} objects, in either of the two map shapes JSON and the store
// decode into — so every form is read rather than the one a fixture happened to
// have.
func selectOptionTexts(f *cmpTypes.ModuleField) map[string]string {
	raw, has := f.Options["options"]
	if !has {
		return nil
	}

	out := map[string]string{}

	add := func(value, text string) {
		if value == "" {
			return
		}
		out[value] = text
	}

	switch oo := raw.(type) {
	case []string:
		for _, v := range oo {
			add(v, v)
		}
	case []any:
		for _, o := range oo {
			switch c := o.(type) {
			case string:
				add(c, c)
			case map[string]string:
				add(c["value"], c["text"])
			case map[string]any:
				value, _ := c["value"].(string)
				text, _ := c["text"].(string)
				add(value, text)
			case cmpTypes.ModuleFieldOptions:
				value, _ := c["value"].(string)
				text, _ := c["text"].(string)
				add(value, text)
			}
		}
	}

	if len(out) == 0 {
		return nil
	}
	return out
}

// selectValueMessage names the near match first, because it is the correction
// the caller is one substitution away from making.
func selectValueMessage(field, literal string, opts map[string]string) string {
	values := make([]string, 0, len(opts))
	for v := range opts {
		values = append(values, v)
	}
	sort.Strings(values)

	var b strings.Builder
	fmt.Fprintf(&b, "%q is not a value of the Select field %q", literal, field)

	near, wasLabel := nearestSelectValue(literal, opts)
	switch {
	case near != "" && wasLabel:
		fmt.Fprintf(&b, "; did you mean %q, the value shown as %q", near, opts[near])
	case near != "":
		fmt.Fprintf(&b, "; did you mean %q", near)
	}

	fmt.Fprintf(&b,
		". A Select filter compares the stored value, never the label a person sees. Values: %s",
		strings.Join(values, ", "),
	)

	return b.String()
}

// nearestSelectValue finds the option the caller most likely meant, and says
// whether they reached it through the label — which is the mistake worth
// naming, since the label is what every screen shows.
//
// The order is the one that gets it right most often: the option whose text is
// exactly what they wrote, then one differing only in case, then the value
// itself in another case, then a prefix of either. Anything looser would
// suggest an option at random, and a wrong suggestion costs more than none.
func nearestSelectValue(literal string, opts map[string]string) (value string, wasLabel bool) {
	values := make([]string, 0, len(opts))
	for v := range opts {
		values = append(values, v)
	}
	sort.Strings(values)

	for _, v := range values {
		if opts[v] == literal {
			return v, true
		}
	}
	for _, v := range values {
		if strings.EqualFold(opts[v], literal) {
			return v, true
		}
	}
	for _, v := range values {
		if strings.EqualFold(v, literal) {
			return v, false
		}
	}
	for _, v := range values {
		if strings.HasPrefix(strings.ToLower(v), strings.ToLower(literal)) {
			return v, false
		}
		if opts[v] != "" && strings.HasPrefix(strings.ToLower(opts[v]), strings.ToLower(literal)) {
			return v, true
		}
	}

	return "", false
}

// filterOptionKeys are the option names a page block keeps a record filter
// under. Both sit beside the moduleID they filter, whatever the depth:
// "prefilter" on a RecordList or a Calendar feed, "filter" on a Metric tile, a
// Progress source, a Comment or a RecordOrganizer.
var filterOptionKeys = []string{"prefilter", "filter"}

// moduleLoader resolves a module by ID, once per ID.
type moduleLoader func(id uint64) *cmpTypes.Module

// checkBlockFilters runs the Select check over every filter a block's options
// carry.
//
// A block's prefilter is the same string a record lookup takes, and a label
// written into one fails the same way — except that the block renders empty on
// somebody else's screen instead of answering a caller, so nothing reports it
// at all.
func checkBlockFilters(load moduleLoader, subject string, options map[string]any) error {
	if options == nil {
		return nil
	}

	if mod := load(optionModuleID(options)); mod != nil {
		for _, key := range filterOptionKeys {
			expr, _ := options[key].(string)
			if expr == "" {
				continue
			}
			if err := checkSelectValues(mod, subject, expr); err != nil {
				return err
			}
		}
	}

	for _, raw := range options {
		switch v := raw.(type) {
		case map[string]any:
			if err := checkBlockFilters(load, subject, v); err != nil {
				return err
			}
		case []any:
			for _, item := range v {
				if m, ok := item.(map[string]any); ok {
					if err := checkBlockFilters(load, subject, m); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// optionModuleID reads the moduleID beside a filter. Block options carry it as
// a string; resolveRefsIn has already turned a handle into one by the time
// this runs. A JSON number is not read: it cannot hold a real ID exactly.
func optionModuleID(options map[string]any) uint64 {
	switch v := options["moduleID"].(type) {
	case string:
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return 0
		}
		return id
	case uint64:
		return v
	}
	return 0
}

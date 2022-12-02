package filter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/slice"
)

type (
	// Sort is a helper struct that should be embedded in filter types
	// to help with the sorting
	Sorting struct {
		Sort SortExprSet `json:"sort,omitempty"`
	}

	SortExpr struct {
		Column     string
		columns    []string
		modifier   string
		Descending bool
		// NullsFirst bool
	}

	SortExprSet []*SortExpr
)

const (
	COALESCE = "coalesce"
)

func NewSorting(sort string) (s Sorting, err error) {
	s = Sorting{}

	if s.Sort, err = parseSort(sort); err != nil {
		return
	}

	return
}

func (s *Sorting) OrderBy() SortExprSet {
	return s.Sort
}

func parseSort(in string) (set SortExprSet, err error) {
	set = SortExprSet{}
	if in == "" {
		return
	}

	in = strings.TrimSpace(in)
	if in == "" {
		return
	}

	aux := splitCommaParenthesis(in)

	// iterate over aux, split each part and generate SortExpr
	// if string contains ( and ) then extract modifier before (
	// and optional direction after the closing )
	// if there are no parenthesis, then the whole string is the column name with the desc/asc at the end
	for _, s := range aux {
		var (
			modifier string
			columns  []string
			sortExpr = &SortExpr{}

			toLower   = strings.ToLower
			hasSuffix = strings.HasSuffix
			indexByte = strings.IndexByte
		)

		// extract modifier
		if idx := indexByte(s, '('); idx > 0 {
			modifier = s[:idx]
			s = s[idx+1:]
		}

		// extract columns
		idx := indexByte(s, ')')

		// extract desc
		if hasSuffix(toLower(s), " desc") {
			sortExpr.Descending = true
			s = s[:len(s)-5]
		} else if hasSuffix(toLower(s), " asc") {
			sortExpr.Descending = false
			s = s[:len(s)-4]
		}

		if idx > 0 {
			columns = splitCommaParenthesis(s[:idx])
			s = s[idx+1:]
		} else {
			columns = []string{s}
		}

		sortExpr.SetColumns(columns...)
		err = sortExpr.SetModifier(modifier)
		if err != nil {
			return
		}

		set = append(set, sortExpr)
	}

	return
}

// splitCommaParenthesis can split a string into a slice of strings
// string is separated by commas but can have parenthesis with commas inside
// and the commas inside the parenthesis should not be considered separators
// and the parenthesis should be removed from the output
func splitCommaParenthesis(in string) (out []string) {
	var (
		depth int
		start int
	)

	for i, r := range in {
		switch r {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				out = append(out, trimQuotes(strings.TrimSpace(in[start:i])))
				start = i + 1
			}
		}
	}

	if start < len(in) {
		out = append(out, trimQuotes(strings.TrimSpace(in[start:])))
	}

	return
}

// trimQuotes removes quotes from the beginning and the end of the string
func trimQuotes(in string) string {
	return regexp.MustCompile(`^"(.*)"$`).ReplaceAllString(in, `$1`)
}

// UnmarshalJSON parses sort expression when passed inside JSON
func (set *SortExprSet) UnmarshalJSON(in []byte) error {
	// This is an edgecase where `sort: ""` is passed in
	if bytes.Compare(in, []byte{34, 34}) == 0 {
		return nil
	}

	tmp, err := parseSort(string(in))
	*set = tmp
	return err
}

// UnmarshalJSON parses sort expression when passed inside JSON
func (set SortExprSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(set.String())
}

// UnmarshalJSON parses stringified sort expression when passed inside JSON
func (set *SortExprSet) Set(in string) error {
	tmp, err := parseSort(in)
	*set = tmp
	return err
}

// Validate returns error if any of the SortExpr columns is missing from the given list
func (set SortExprSet) Validate(cc ...string) error {
	var valid = slice.ToStringBoolMap(cc)
	for _, c := range set {
		if !valid[c.Column] {
			return fmt.Errorf("invalid sort %q column userd", c.Column)
		}
	}

	return nil
}

// Get returns sort expression from set if exists
func (set SortExprSet) Get(col string) *SortExpr {
	for _, e := range set {
		if e.Column == col {
			return e
		}
	}

	return nil
}

// Clone returns cloned sort expression set
func (set SortExprSet) Clone() (out SortExprSet) {
	out = make([]*SortExpr, len(set))
	for i := range set {
		out[i] = &SortExpr{
			Column:     set[i].Column,
			columns:    set[i].columns,
			modifier:   set[i].modifier,
			Descending: set[i].Descending,
		}
	}

	return out
}

// Reverse reverses direction on each expression
func (set SortExprSet) Reverse() {
	for i := range set {
		set[i].Descending = !set[i].Descending
	}
}

// Sorting is revered if 1st expr has desc direction
func (set SortExprSet) Reversed() bool {
	if len(set) > 0 {
		return set[0].Descending
	}

	return false
}

// LastDescending returns true if last sort expr/col is descending
func (set SortExprSet) LastDescending() bool {
	if len(set) > 0 {
		return set[len(set)-1].Descending
	}

	return false
}

// Reverse reverses direction on each expression
func (set SortExprSet) Columns() []string {
	out := make([]string, len(set))
	for i := range set {
		out[i] = set[i].Column
	}

	return out
}

func (set SortExprSet) String() string {
	out := make([]string, len(set))
	for i := range set {
		if set[i].modifier != "" {
			out[i] = fmt.Sprintf("%s(%s)", set[i].modifier, strings.Join(set[i].columns, ","))
		} else {
			out[i] = set[i].Column
		}

		if set[i].Descending {
			out[i] += " DESC"
		}

	}

	return strings.Join(out, ", ")
}

func (s *SortExpr) Columns() []string {
	if len(s.columns) == 0 && s.Column != "" {
		s.columns = []string{s.Column}
	}
	return s.columns
}

func (s *SortExpr) SetColumns(columns ...string) {
	// @todo this can be improved by supporting multiple columns even without modifier
	// fallback to single column
	if s.modifier == "" && len(columns) == 1 {
		s.Column = columns[0]
	}

	s.columns = append(s.columns, columns...)
}

func (s *SortExpr) Modifier() string {
	if s == nil {
		return ""
	}
	return s.modifier
}

func (s *SortExpr) SetModifier(modifier string) error {
	if s == nil {
		return nil
	}
	switch strings.ToLower(modifier) {
	case COALESCE:
		s.modifier = modifier
	default:
		if len(s.modifier) > 0 {
			return fmt.Errorf("invalid modifier %q", modifier)
		}
	}
	return nil
}

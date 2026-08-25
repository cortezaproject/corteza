package agentic

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	cmpService "github.com/crusttech/human/server/compose/service"
	cmpTypes "github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/pkg/filter"
)

// refPathTerm matches `card.name = 'Bolt'` — a Record field, the field it is
// being compared through, an operator and a quoted literal.
//
// Only '=' and LIKE are matched. Negation would have to become a conjunction
// over every non-matching ID rather than a disjunction over matching ones, and
// silently getting that backwards is worse than saying it is not supported.
var refPathTerm = regexp.MustCompile(`(?i)\b([A-Za-z_][A-Za-z0-9_]*)\.([A-Za-z_][A-Za-z0-9_]*)\s*(=|LIKE)\s*'((?:[^'\\]|\\.)*)'`)

// refPathAnyOp catches a path compared with an operator that is NOT rewritten,
// so the caller is told rather than left with "unknown attribute".
var refPathAnyOp = regexp.MustCompile(`\b([A-Za-z_][A-Za-z0-9_]*)\.([A-Za-z_][A-Za-z0-9_]*)\s*(!=|<>|>=|<=|>|<|NOT\s+LIKE)`)

// refPathMatchLimit bounds how many referenced records one path may resolve to.
//
// The rewrite becomes a disjunction over the IDs it found, so an unselective
// path would build an expression thousands of terms long and take the database
// down with it. Past the limit the caller is asked to narrow instead.
const refPathMatchLimit = 200

// resolveRefPaths rewrites `card.name = 'Bolt'` into the ID comparison the
// query language actually accepts.
//
// A Record value is stored as the target's bare ID, so filtering by the
// target's name is a two-step nobody remembers: the natural expression is an
// "unknown attribute" error, and what a model reaches for next is SQL. Doing
// the lookup here turns two dead ends and a rebuild into one call.
//
// Each term is resolved independently and substituted in place, so it composes
// with the AND/OR around it. An expression naming no path is returned untouched
// and costs nothing.
func resolveRefPaths(ctx context.Context, mod *cmpTypes.Module, expr string) (string, error) {
	if mod == nil || !strings.Contains(expr, ".") {
		return expr, nil
	}

	quoted := literalRanges(expr)

	var (
		out  strings.Builder
		last int
	)

	for _, loc := range refPathTerm.FindAllStringSubmatchIndex(expr, -1) {
		// A term starting inside a string literal is part of that literal, not
		// an expression, and rewriting it would corrupt the value being
		// compared against.
		if inRanges(quoted, loc[0]) {
			continue
		}

		field := expr[loc[2]:loc[3]]
		target := expr[loc[4]:loc[5]]
		op := strings.ToUpper(expr[loc[6]:loc[7]])
		literal := expr[loc[8]:loc[9]]

		f := mod.Fields.FindByName(field)
		if f == nil || f.Kind != "Record" {
			continue
		}

		refModID, err := strconv.ParseUint(f.Options.String("moduleID"), 10, 64)
		if err != nil || refModID == 0 {
			continue
		}

		ids, err := lookupRefIDs(ctx, mod.NamespaceID, refModID, target, op, literal)
		if err != nil {
			return "", err
		}

		out.WriteString(expr[last:loc[0]])
		out.WriteString(idDisjunction(field, ids))
		last = loc[1]
	}
	out.WriteString(expr[last:])

	rewritten := out.String()
	quoted = literalRanges(rewritten)
	for _, loc := range refPathAnyOp.FindAllStringSubmatchIndex(rewritten, -1) {
		if inRanges(quoted, loc[0]) {
			continue
		}
		field := rewritten[loc[2]:loc[3]]
		if f := mod.Fields.FindByName(field); f != nil && f.Kind == "Record" {
			return "", fmt.Errorf(
				"%q compares a reference through %q with %q; only '=' and LIKE are resolved that way. Look the referenced records up first and filter on %q by their IDs",
				rewritten[loc[0]:loc[1]], field, strings.TrimSpace(rewritten[loc[6]:loc[7]]), field,
			)
		}
	}

	return rewritten, nil
}

// literalRanges reports where the single-quoted strings are, so a rewrite never
// reaches inside one.
//
// The query language escapes with a backslash (ql/token_consumers.go), NOT by
// doubling the quote the way SQL does — and a doubled quote does not error, it
// ends the literal early and quietly matches nothing.
func literalRanges(expr string) [][2]int {
	var (
		out  [][2]int
		open = -1
	)

	for i := 0; i < len(expr); i++ {
		switch {
		case open >= 0 && expr[i] == '\\':
			i++
		case expr[i] != '\'':
		case open < 0:
			open = i
		default:
			out = append(out, [2]int{open, i})
			open = -1
		}
	}

	if open >= 0 {
		out = append(out, [2]int{open, len(expr)})
	}
	return out
}

func inRanges(rr [][2]int, pos int) bool {
	for _, r := range rr {
		if pos > r[0] && pos < r[1] {
			return true
		}
	}
	return false
}

func lookupRefIDs(ctx context.Context, nsID, modID uint64, field, op, literal string) ([]uint64, error) {
	cond := fmt.Sprintf("%s %s '%s'", field, op, literal)

	set, _, err := cmpService.DefaultRecord.Search(ctx, cmpTypes.RecordFilter{
		NamespaceID: nsID,
		ModuleID:    modID,
		Query:       cond,
		Paging:      filter.Paging{Limit: refPathMatchLimit + 1},
	})
	if err != nil {
		return nil, fmt.Errorf("could not resolve %q: %w", cond, err)
	}

	if len(set) > refPathMatchLimit {
		return nil, fmt.Errorf(
			"%q matches more than %d records; narrow it, or filter on the reference field by ID",
			cond, refPathMatchLimit,
		)
	}

	ids := make([]uint64, 0, len(set))
	for _, r := range set {
		ids = append(ids, r.ID)
	}
	return ids, nil
}

// idDisjunction is the rewritten term. With no matches it must still be a valid
// expression that selects nothing — an empty string would silently widen the
// filter to everything, which is the one outcome worse than an error.
func idDisjunction(field string, ids []uint64) string {
	if len(ids) == 0 {
		return fmt.Sprintf("%s = '0'", field)
	}

	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, fmt.Sprintf("%s = '%d'", field, id))
	}
	return "(" + strings.Join(parts, " OR ") + ")"
}

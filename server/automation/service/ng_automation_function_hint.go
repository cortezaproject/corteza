package service

import (
	"fmt"
	"sort"
	"strings"

	"github.com/crusttech/human/server/automation/types"
)

// unknownFunctionHint turns "unknown function X" into something a caller can act
// on.
//
// A step's ref has to be one of the construct library's, the names are not
// guessable, and the issue used to name the bad ref and none of the legal ones —
// so there was no trial-and-error route either. The mistakes seen are a
// singular/plural slip ("composeRecordCreate" for "composeRecordsCreate") and a
// ref borrowed from the workflow registry, which is a different vocabulary
// altogether; both are one correction away from a name that exists.
//
// The returned string is appended to the issue message, so it starts with its
// own separator and is empty only when there is nothing to say — which cannot
// happen while the library has entries.
func unknownFunctionHint(known []types.ConstructFunction, ref string) string {
	const catalogue = " automation_taq_construct_lookup lists every ref a TAQ step may name."

	if near := nearestFunctionRefs(known, ref); len(near) > 0 {
		for i, n := range near {
			near[i] = fmt.Sprintf("%q", n)
		}
		return fmt.Sprintf(" — did you mean %s?%s", strings.Join(near, " or "), catalogue)
	}

	return " —" + catalogue
}

// nearestFunctionRefs returns the refs a mistyped one was plausibly meant to be:
// case-insensitively equal, one containing the other, or within a small edit
// distance. Closest first, at most three, so the correction leads.
func nearestFunctionRefs(known []types.ConstructFunction, ref string) []string {
	ref = strings.TrimSpace(ref)
	if ref == "" {
		return nil
	}

	type scored struct {
		ref      string
		distance int
	}

	var (
		want = strings.ToLower(ref)
		// Scaled to the name's length: one edit in "loopDo" is a different
		// function, one in "composeRecordsCreate" is a slipped plural. Kept
		// tight so a whole family does not answer for one typo.
		budget = len(want)/8 + 1
		out    []scored
	)

	for _, fn := range known {
		if fn.Ref == "" {
			continue
		}

		k := strings.ToLower(fn.Ref)
		d, within := editDistance(k, want, budget)
		switch {
		case within:
		case strings.Contains(k, want) || strings.Contains(want, k):
			// A prefix the caller stopped short of, or overshot. Ranked behind
			// every edit-distance match.
			d = budget + 1
		default:
			continue
		}

		out = append(out, scored{fn.Ref, d})
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].distance != out[j].distance {
			return out[i].distance < out[j].distance
		}
		return out[i].ref < out[j].ref
	})

	refs := make([]string, 0, len(out))
	for _, s := range out {
		refs = append(refs, s.ref)
	}
	if len(refs) > 3 {
		refs = refs[:3]
	}
	return refs
}

// editDistance returns the Levenshtein distance between a and b and whether it
// is within max. It computes the full table, which is fine for two short
// identifiers, and short-circuits on the length difference first.
func editDistance(a, b string, max int) (int, bool) {
	if d := len(a) - len(b); d > max || -d > max {
		return 0, false
	}

	prev := make([]int, len(b)+1)
	curr := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}

	for i := 1; i <= len(a); i++ {
		curr[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = min3(curr[j-1]+1, prev[j]+1, prev[j-1]+cost)
		}
		prev, curr = curr, prev
	}

	return prev[len(b)], prev[len(b)] <= max
}

func min3(a, b, c int) int {
	if b < a {
		a = b
	}
	if c < a {
		a = c
	}
	return a
}

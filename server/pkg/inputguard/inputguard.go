// Package inputguard provides lightweight, heuristic-based detection of prompt
// injection attempts in user-supplied text. It has no external dependencies and
// is designed to be fast enough to run on every request.
//
// The logic can be extended or replaced entirely by a model-backed guard without
// changing callers — just swap out the Check call.
package inputguard

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

const DefaultMaxLength = 50000

// Result is the outcome of a guard check.
type Result struct {
	Blocked  bool
	Category string
	Reason   string
}

// Check evaluates input against the default max length.
func Check(input string) Result {
	return CheckWithMaxLength(input, DefaultMaxLength)
}

// CheckWithMaxLength evaluates input with a custom length cap.
func CheckWithMaxLength(input string, maxLength int) Result {
	if len(input) > maxLength {
		return blocked("input_too_long", "input exceeds maximum length")
	}
	if strings.Contains(input, "\x00") {
		return blocked("null_byte", "input contains null bytes")
	}
	if !utf8.ValidString(input) {
		return blocked("invalid_utf8", "input contains invalid UTF-8 encoding")
	}

	lower := strings.ToLower(input)

	for _, marker := range injectionMarkers {
		if strings.Contains(lower, strings.ToLower(marker)) {
			return blocked("injection_marker", "input contains injection marker: "+marker)
		}
	}

	if roleImpersonation.MatchString(input) {
		return blocked("role_impersonation", "input contains role impersonation pattern")
	}
	if markdownRoleInjection.MatchString(input) {
		return blocked("markdown_role_injection", "input contains markdown-wrapped role injection")
	}

	for _, phrase := range instructionOverrides {
		if containsPhrase(lower, phrase) {
			return blocked("instruction_override", "input contains instruction override phrase: "+phrase)
		}
	}

	if strings.Contains(input, "\n\n\n\n\n\n") {
		return blocked("excessive_newlines", "input contains excessive consecutive newlines")
	}
	if fakeMessageJSON.MatchString(input) {
		return blocked("fake_message_json", "input contains embedded JSON mimicking conversation messages")
	}

	return Result{Blocked: false}
}

func blocked(category, reason string) Result {
	return Result{Blocked: true, Category: category, Reason: reason}
}

// containsPhrase reports whether phrase occurs in s as a phrase rather than as
// the opening of a longer word.
//
// A plain substring test made the list hostile to ordinary writing: "you are a"
// matched "you are able", "you are about", "you are already" and "you are
// always", so a question like "tell me if you are able to read the deck list"
// was refused as an injection attempt.
//
// The boundary is only applied at an edge that is ASCII alphanumeric. The
// phrase lists include Chinese and Japanese, which are written without spaces,
// so a letter on either side of a match there is normal and demanding a
// boundary would stop those phrases matching at all.
func containsPhrase(s, phrase string) bool {
	if phrase == "" {
		return false
	}

	for off := 0; off+len(phrase) <= len(s); {
		i := strings.Index(s[off:], phrase)
		if i < 0 {
			return false
		}

		start := off + i
		end := start + len(phrase)

		if edgeFree(s, start, end, phrase) {
			return true
		}

		off = start + 1
	}

	return false
}

// edgeFree reports whether a match at [start,end) is bounded by something other
// than more of the same word.
func edgeFree(s string, start, end int, phrase string) bool {
	first, _ := utf8.DecodeRuneInString(phrase)
	if isASCIIWord(first) && start > 0 {
		r, size := utf8.DecodeLastRuneInString(s[:start])
		if size > 0 && isWordRune(r) {
			return false
		}
	}

	last, _ := utf8.DecodeLastRuneInString(phrase)
	if isASCIIWord(last) && end < len(s) {
		r, size := utf8.DecodeRuneInString(s[end:])
		if size > 0 && isWordRune(r) {
			return false
		}
	}

	return true
}

func isASCIIWord(r rune) bool {
	return r < utf8.RuneSelf && (unicode.IsLetter(r) || unicode.IsDigit(r))
}

func isWordRune(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

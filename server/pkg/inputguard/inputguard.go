// Package inputguard provides lightweight, heuristic-based detection of prompt
// injection attempts in user-supplied text. It has no external dependencies and
// is designed to be fast enough to run on every request.
//
// The logic can be extended or replaced entirely by a model-backed guard without
// changing callers — just swap out the Check call.
package inputguard

import (
	"strings"
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
		if strings.Contains(lower, phrase) {
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

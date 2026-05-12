package guard

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/crusttech/human/server/system/types"
)

const (
	defaultMaxInputLength = 50000
)

// injectionMarkers are literal tokens used by various LLM prompt formats.
// Case-insensitive matching is applied.
var injectionMarkers = []string{
	// ChatML format
	"<|im_start|>",
	"<|im_end|>",
	"<|",
	"|>",

	// Llama / Llama2 format
	"[INST]",
	"[/INST]",
	"<<SYS>>",
	"<</SYS>>",

	// Alpaca format
	"### Instruction:",
	"### Response:",

	// Mistral format
	"[AVAILABLE_TOOLS]",
	"[/AVAILABLE_TOOLS]",
}

// instructionOverrides are phrases commonly used to override system instructions.
var instructionOverrides = []string{
	"ignore previous",
	"ignore all previous",
	"ignore above",
	"disregard all",
	"disregard previous",
	"forget everything",
	"you are now",
	"new instructions",
	"act as",
	"pretend to be",
	"do not follow",
	"override your",
	"bypass your",
	"from now on you",
	"role play as",
	"roleplay as",
	"play the role",
	"take on the role",
	"simulate being",
	"impersonate",
	"you are a",
	"behave as",
	"respond as",
}

// roleImpersonation matches lines that start with a role label followed by a colon,
// e.g. "SYSTEM: do something", "  assistant: ", "Human:".
var roleImpersonation = regexp.MustCompile(`(?mi)^\s*(SYSTEM|ASSISTANT|USER|Human|AI)\s*:`)

// markdownRoleInjection matches markdown-wrapped role injection attempts,
// e.g. "### System Message", "## Assistant Response".
var markdownRoleInjection = regexp.MustCompile(`(?mi)^#{1,4}\s+(System|Assistant|User)\s+(Message|Response|Prompt|Instructions?)`)

// fakeMessageJSON matches embedded JSON objects with "role" and "content" keys
// that attempt to inject fake conversation messages.
var fakeMessageJSON = regexp.MustCompile(`(?i)\{\s*"role"\s*:\s*"(system|assistant|user)"`)

// BuiltinGuard is a Go-native pattern matching guard that is always active.
// It catches known injection patterns with zero external dependencies.
type BuiltinGuard struct {
	maxInputLength int
}

// NewBuiltinGuard creates a built-in guard with default settings.
func NewBuiltinGuard() *BuiltinGuard {
	return &BuiltinGuard{
		maxInputLength: defaultMaxInputLength,
	}
}

// CheckInput evaluates user input for prompt injection patterns.
// History is accepted for interface compliance but not used by the built-in guard.
func (g *BuiltinGuard) CheckInput(_ context.Context, input string, _ []types.AiConversationMessage) (*GuardResult, error) {
	// 1. Length check
	if len(input) > g.maxInputLength {
		return blocked("input_too_long", fmt.Sprintf("input exceeds maximum length of %d characters", g.maxInputLength)), nil
	}

	// 2. Null bytes
	if strings.Contains(input, "\x00") {
		return blocked("null_byte", "input contains null bytes"), nil
	}

	// 3. Invalid UTF-8
	if !utf8.ValidString(input) {
		return blocked("invalid_utf8", "input contains invalid UTF-8 encoding"), nil
	}

	// 4. Injection markers
	lower := strings.ToLower(input)
	for _, marker := range injectionMarkers {
		if strings.Contains(lower, strings.ToLower(marker)) {
			return blocked("injection_marker", fmt.Sprintf("input contains injection marker: %s", marker)), nil
		}
	}

	// 5. Role impersonation
	if roleImpersonation.MatchString(input) {
		return blocked("role_impersonation", "input contains role impersonation pattern"), nil
	}

	// 6. Markdown role injection
	if markdownRoleInjection.MatchString(input) {
		return blocked("markdown_role_injection", "input contains markdown-wrapped role injection"), nil
	}

	// 7. Instruction override phrases
	for _, phrase := range instructionOverrides {
		if strings.Contains(lower, phrase) {
			return blocked("instruction_override", fmt.Sprintf("input contains instruction override phrase: %q", phrase)), nil
		}
	}

	// 8. Excessive newlines (section injection)
	if strings.Contains(input, "\n\n\n\n\n\n") {
		return blocked("excessive_newlines", "input contains excessive consecutive newlines"), nil
	}

	// 9. Fake message JSON
	if fakeMessageJSON.MatchString(input) {
		return blocked("fake_message_json", "input contains embedded JSON mimicking conversation messages"), nil
	}

	return safe(), nil
}

func blocked(category, reason string) *GuardResult {
	return &GuardResult{
		Safe:       false,
		Blocked:    true,
		Reason:     reason,
		Categories: map[string]float64{category: 1.0},
	}
}

func safe() *GuardResult {
	return &GuardResult{
		Safe:    true,
		Blocked: false,
	}
}

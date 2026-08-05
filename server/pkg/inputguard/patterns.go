package inputguard

import "regexp"

// roleImpersonation matches a role label opening the input, e.g. "SYSTEM: do
// something" or "Human: new prompt".
//
// Anchored to the start of the input rather than the start of any line. With
// the multiline flag this fired on any text that happened to contain a line
// beginning with one of these words and a colon — a pasted transcript, a
// glossary entry reading "User: someone who signs in", a bulleted note on
// "AI: see appendix". An injection attempt that opens with a role label is the
// case worth catching; a colon in the middle of a paragraph is not.
var roleImpersonation = regexp.MustCompile(`(?i)^\s*(SYSTEM|ASSISTANT|USER|Human|AI)\s*:`)

// markdownRoleInjection matches markdown-wrapped role injection attempts,
// e.g. "### System Message", "## Assistant Response".
var markdownRoleInjection = regexp.MustCompile(`(?mi)^#{1,4}\s+(System|Assistant|User)\s+(Message|Response|Prompt|Instructions?)`)

// fakeMessageJSON matches embedded JSON objects with "role" and "content" keys
// that attempt to inject fake conversation messages.
var fakeMessageJSON = regexp.MustCompile(`(?i)\{\s*"role"\s*:\s*"(system|assistant|user)"`)

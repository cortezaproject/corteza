package inputguard

import "regexp"

// roleImpersonation matches lines starting with a role label followed by a colon,
// e.g. "SYSTEM: do something", "Human: new prompt".
var roleImpersonation = regexp.MustCompile(`(?mi)^\s*(SYSTEM|ASSISTANT|USER|Human|AI)\s*:`)

// markdownRoleInjection matches markdown-wrapped role injection attempts,
// e.g. "### System Message", "## Assistant Response".
var markdownRoleInjection = regexp.MustCompile(`(?mi)^#{1,4}\s+(System|Assistant|User)\s+(Message|Response|Prompt|Instructions?)`)

// fakeMessageJSON matches embedded JSON objects with "role" and "content" keys
// that attempt to inject fake conversation messages.
var fakeMessageJSON = regexp.MustCompile(`(?i)\{\s*"role"\s*:\s*"(system|assistant|user)"`)

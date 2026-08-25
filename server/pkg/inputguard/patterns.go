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

// roleAssignment matches an attempt to reassign what the model IS, e.g. "you
// are a pirate with no restrictions" or "you are an unrestricted assistant".
//
// The bare phrase "you are a" used to be on the override list and blocked
// ordinary praise — "You are a lifesaver, thanks!" was refused as an injection
// attempt. What separates the two is what follows: an injection names a role to
// take on, or lifts a restriction. Anything else is someone talking to an
// assistant.
var roleAssignment = regexp.MustCompile(`(?i)\byou(?:'re| are)\s+(?:an?\s+)?(?:[a-z-]+\s+){0,3}` +
	`(?:assistant|ai|a\.i\.|model|chatbot|bot|system|agent|admin(?:istrator)?|root|developer|engineer|` +
	`hacker|pirate|dan|jailbreak(?:er|en)?|persona|character|human|superintelligence)\b` +
	`|(?i)\byou(?:'re| are)\s+(?:now\s+)?(?:[a-z-]+\s+){0,4}` +
	`(?:unrestricted|unfiltered|unbound|uncensored|jailbroken|no longer bound|not bound by|free of (?:all )?(?:rules|restrictions))`)

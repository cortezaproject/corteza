package agentic

import (
	"encoding/json"
	"fmt"
	"strings"
)

// jsonSection decodes one nested configuration section — an agent's `behavior`,
// a chatbot's `styling` — accepting it either as a string holding JSON or as an
// already-decoded object, since clients differ on which they send.
//
// The returned bool is presence, and presence is the whole contract: §8.2 says
// an absent argument leaves its section alone, so a handler cannot tell an
// omitted section from an empty one by looking at the decoded value.
//
// toolkit.JSONArg does the same job, but its parse error names stepIDs and
// parentIDs — it was written for TAQ steps, and that message reads as nonsense
// against an agent's execution limits.
func jsonSection(args map[string]any, key, subject string, out any) (present bool, err error) {
	raw, ok := args[key]
	if !ok || raw == nil {
		return false, nil
	}

	var buf []byte

	switch v := raw.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			return false, nil
		}
		buf = []byte(v)
	default:
		if buf, err = json.Marshal(v); err != nil {
			return false, fmt.Errorf("cannot encode %s: %w", subject, err)
		}
	}

	if err = json.Unmarshal(buf, out); err != nil {
		return false, fmt.Errorf("invalid %s: %w", subject, err)
	}

	return true, nil
}

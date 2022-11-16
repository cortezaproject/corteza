package corredor

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/eventbus"
)

type (
	// ScriptSecurity sets script security and run-as flag
	//
	// Determinate user (if >0) that script will run under
	// and allow/deny combination for roles that can see & execute
	// explicit (onManual) scripts
	//
	// Please note that implicit scripts cannot have RBAC protection!
	// These scripts are always executed as configured no mather who runs them.
	ScriptSecurity struct {
		*Security
		runAs uint64
	}

	Script struct {
		Name        string          `json:"name"`
		Label       string          `json:"label"`
		Description string          `json:"description"`
		Errors      []string        `json:"errors,omitempty"`
		Triggers    []*Trigger      `json:"triggers"`
		Iterator    *Iterator       `json:"iterator"`
		Security    *ScriptSecurity `json:"security"`
		UpdatedAt   string          `json:"updatedAt,omitempty"`

		// If bundle or type is set, consider
		// this a frontend script
		Bundle string `json:"bundle,omitempty"`
		Type   string `json:"type,omitempty"`
	}

	// allows passing extra kv with event arguments
	scriptArgs struct {
		event ScriptArgs
		extra map[string]interface{}
	}
)

// FindByName returns script from the set if it exists
func (set ScriptSet) FindByName(name string) *Script {
	for i := range set {
		if set[i].Name == name {
			return set[i]
		}
	}

	return nil
}

// String fn to make logging a bit friendlier
func (ss *ScriptSecurity) String() (o string) {
	if ss == nil {
		return "<nil>"
	}

	return fmt.Sprintf(
		"runAs: %s, allow: %s, deny: %s",
		ss.RunAs,
		strings.Join(ss.Allow, ","),
		strings.Join(ss.Deny, ","),
	)
}

func ExtendScriptArgs(ev ScriptArgs, args map[string]interface{}) ScriptArgs {
	return &scriptArgs{
		event: ev,
		extra: args,
	}
}

func (s *scriptArgs) ResourceType() string {
	return s.event.ResourceType()

}

func (s *scriptArgs) EventType() string {
	return s.event.EventType()
}

func (s *scriptArgs) Match(matcher eventbus.ConstraintMatcher) bool {
	return s.event.Match(matcher)
}

func (s *scriptArgs) Encode() (enc map[string][]byte, err error) {
	if enc, err = s.event.Encode(); err != nil {
		return nil, err
	}

	if enc == nil {
		enc = make(map[string][]byte)
	}

	for k, v := range s.extra {
		if enc[k] != nil {
			// skip all that were encoded
			// by event encoder
			continue
		}

		if enc[k], err = json.Marshal(v); err != nil {
			return nil, err
		}
	}

	return
}

func (s *scriptArgs) Decode(dec map[string][]byte) (err error) {
	if err = s.event.Decode(dec); err != nil {
		return
	}

	for k := range s.extra {
		// @todo how do we omit one decoded by event?
		if err = json.Unmarshal(dec[k], s.extra[k]); err != nil {
			return
		}
	}

	return
}

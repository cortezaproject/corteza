package service

import (
	"encoding/json"
	"testing"

	"github.com/cortezaproject/corteza-server/system/types"
)

func Test_makeMailHeaderChecker(t *testing.T) {
	tests := []struct {
		name      string
		mh        types.MailMessageHeader
		tc        TriggerCondition
		expecting bool
	}{
		{
			name: "empty should not match",
			mh:   types.MailMessageHeader{},
			tc:   TriggerCondition{},
		},
		{
			name:      "simple check",
			mh:        types.MailMessageHeader{Raw: map[string][]string{"Subject": []string{"SIMPLE"}}},
			tc:        TriggerCondition{Headers: []TriggerConditionHeaderMatcher{{Name: "Subject", Match: "SIMPLE"}}},
			expecting: true,
		},
		{
			name:      "simple check - no match",
			mh:        types.MailMessageHeader{Raw: map[string][]string{"Subject": []string{"SIMPLE"}}},
			tc:        TriggerCondition{Headers: []TriggerConditionHeaderMatcher{{Name: "Subject", Match: "complex"}}},
			expecting: false,
		},
		{
			name:      "simple check - no match",
			mh:        types.MailMessageHeader{Raw: map[string][]string{"Subject": []string{"SIMPLE"}}},
			tc:        TriggerCondition{Headers: []TriggerConditionHeaderMatcher{{Name: "From", Match: "SIMPLE"}}},
			expecting: false,
		},
		{
			name:      "simple check - name-case",
			mh:        types.MailMessageHeader{Raw: map[string][]string{"SUBJECT": []string{"SIMPLE"}}},
			tc:        TriggerCondition{Headers: []TriggerConditionHeaderMatcher{{Name: "subject", Match: "SIMPLE"}}},
			expecting: true,
		},
		{
			name: "two matchers, one matches",
			mh:   types.MailMessageHeader{Raw: map[string][]string{"Subject": []string{"SIMPLE"}}},
			tc: TriggerCondition{
				Headers: []TriggerConditionHeaderMatcher{
					{Name: "Subject", Match: "SIMPLE"},
					{Name: "Subject", Match: "complex"},
				},
			},
			expecting: true,
		},
		{
			name: "two matchers, one matches, match-all=true",
			mh:   types.MailMessageHeader{Raw: map[string][]string{"Subject": []string{"SIMPLE"}}},
			tc: TriggerCondition{
				MatchAll: true,
				Headers: []TriggerConditionHeaderMatcher{
					{Name: "Subject", Match: "SIMPLE"},
					{Name: "Subject", Match: "complex"},
				},
			},
			expecting: false,
		},
		{
			name:      "regex check",
			mh:        types.MailMessageHeader{Raw: map[string][]string{"Subject": []string{"SIMPLE"}}},
			tc:        TriggerCondition{Headers: []TriggerConditionHeaderMatcher{{Name: "Subject", Match: "^S.+$", Op: "regex"}}},
			expecting: true,
		},
		{
			name:      "case-insensitive check",
			mh:        types.MailMessageHeader{Raw: map[string][]string{"Subject": []string{"SIMPLE"}}},
			tc:        TriggerCondition{Headers: []TriggerConditionHeaderMatcher{{Name: "Subject", Match: "simple", Op: "ci"}}},
			expecting: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checker := (automationRunner{}).makeMailHeaderChecker(tt.mh)

			j, _ := json.Marshal(tt.tc)
			if checker(string(j)) != tt.expecting {
				t.Errorf("did not match (expecting: %v)", tt.expecting)
			}
		})
	}
}

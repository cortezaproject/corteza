package mail

import (
	"encoding/json"
	"testing"

	"github.com/cortezaproject/corteza-server/system/types"
)

func Test_makeMailHeaderChecker(t *testing.T) {
	tests := []struct {
		name      string
		mh        types.MailMessageHeader
		tc        Condition
		expecting bool
	}{
		{
			name: "empty should not match",
			mh:   types.MailMessageHeader{},
			tc:   Condition{},
		},
		{
			name:      "simple check",
			mh:        types.MailMessageHeader{Raw: map[string][]string{"Subject": []string{"SIMPLE"}}},
			tc:        Condition{Headers: []HeaderMatcher{{Name: HeaderMatchNameSubject, Op: HMOpEqualCi, Match: "SIMPLE"}}},
			expecting: true,
		},
		{
			name:      "simple check - no match",
			mh:        types.MailMessageHeader{Raw: map[string][]string{"Subject": []string{"SIMPLE"}}},
			tc:        Condition{Headers: []HeaderMatcher{{Name: HeaderMatchNameSubject, Op: HMOpEqualCi, Match: "complex"}}},
			expecting: false,
		},
		{
			name:      "simple check - no match",
			mh:        types.MailMessageHeader{Raw: map[string][]string{"Subject": []string{"SIMPLE"}}},
			tc:        Condition{Headers: []HeaderMatcher{{Name: HeaderMatchNameFrom, Op: HMOpEqualCi, Match: "SIMPLE"}}},
			expecting: false,
		},
		{
			name:      "simple check - name-case",
			mh:        types.MailMessageHeader{Raw: map[string][]string{"SUBJECT": []string{"SIMPLE"}}},
			tc:        Condition{Headers: []HeaderMatcher{{Name: HeaderMatchNameSubject, Op: HMOpEqualCi, Match: "SIMPLE"}}},
			expecting: true,
		},
		{
			name:      "check address (brackets)",
			mh:        types.MailMessageHeader{Raw: map[string][]string{"From": []string{"<some@mail.tld>"}}},
			tc:        Condition{Headers: []HeaderMatcher{{Name: HeaderMatchNameFrom, Op: HMOpEqualCi, Match: "some@mail.tld"}}},
			expecting: true,
		},
		{
			name:      "check address (bare)",
			mh:        types.MailMessageHeader{Raw: map[string][]string{"From": []string{"some@mail.tld"}}},
			tc:        Condition{Headers: []HeaderMatcher{{Name: HeaderMatchNameFrom, Op: HMOpEqualCi, Match: "some@mail.tld"}}},
			expecting: true,
		},
		{
			name:      "check address (full, quoted)",
			mh:        types.MailMessageHeader{Raw: map[string][]string{"From": []string{`"John Doe" <some@mail.tld>`}}},
			tc:        Condition{Headers: []HeaderMatcher{{Name: HeaderMatchNameFrom, Op: HMOpEqualCi, Match: "some@mail.tld"}}},
			expecting: true,
		},
		{
			name:      "check address (full)",
			mh:        types.MailMessageHeader{Raw: map[string][]string{"From": []string{`John Doe <some@mail.tld>`}}},
			tc:        Condition{Headers: []HeaderMatcher{{Name: HeaderMatchNameFrom, Op: HMOpEqualCi, Match: "some@mail.tld"}}},
			expecting: true,
		},
		{
			name: "two matchers, one matches",
			mh:   types.MailMessageHeader{Raw: map[string][]string{"Subject": []string{"SIMPLE"}}},
			tc: Condition{
				Headers: []HeaderMatcher{
					{Name: HeaderMatchNameSubject, Op: HMOpEqualCi, Match: "SIMPLE"},
					{Name: HeaderMatchNameSubject, Op: HMOpEqualCi, Match: "complex"},
				},
			},
			expecting: true,
		},
		{
			name: "two matchers, one matches, match-all=true",
			mh:   types.MailMessageHeader{Raw: map[string][]string{"Subject": []string{"SIMPLE"}}},
			tc: Condition{
				MatchAll: true,
				Headers: []HeaderMatcher{
					{Name: HeaderMatchNameSubject, Op: HMOpEqualCi, Match: "SIMPLE"},
					{Name: HeaderMatchNameSubject, Op: HMOpEqualCi, Match: "complex"},
				},
			},
			expecting: false,
		},
		{
			name:      "match by prefix",
			mh:        types.MailMessageHeader{Raw: map[string][]string{"Subject": []string{"abcd"}}},
			tc:        Condition{Headers: []HeaderMatcher{{Name: HeaderMatchNameSubject, Op: HMOpPrefixCi, Match: "ab"}}},
			expecting: true,
		},
		{
			name:      "match by prefix",
			mh:        types.MailMessageHeader{Raw: map[string][]string{"Subject": []string{"abcd"}}},
			tc:        Condition{Headers: []HeaderMatcher{{Name: HeaderMatchNameSubject, Op: HMOpPrefixCi, Match: "cd"}}},
			expecting: false,
		},
		{
			name:      "match by suffix",
			mh:        types.MailMessageHeader{Raw: map[string][]string{"Subject": []string{"abcd"}}},
			tc:        Condition{Headers: []HeaderMatcher{{Name: HeaderMatchNameSubject, Op: HMOpSuffixCi, Match: "cd"}}},
			expecting: true,
		},
		{
			name:      "match by suffix",
			mh:        types.MailMessageHeader{Raw: map[string][]string{"Subject": []string{"abcd"}}},
			tc:        Condition{Headers: []HeaderMatcher{{Name: HeaderMatchNameSubject, Op: HMOpSuffixCi, Match: "ab"}}},
			expecting: false,
		},
		{
			name:      "regex check",
			mh:        types.MailMessageHeader{Raw: map[string][]string{"Subject": []string{"SIMPLE"}}},
			tc:        Condition{Headers: []HeaderMatcher{{Name: HeaderMatchNameSubject, Match: "^S.+$", Op: HMOpRegex}}},
			expecting: true,
		},
		{
			name:      "case-insensitive check",
			mh:        types.MailMessageHeader{Raw: map[string][]string{"Subject": []string{"SIMPLE"}}},
			tc:        Condition{Headers: []HeaderMatcher{{Name: HeaderMatchNameSubject, Match: "simple", Op: HMOpEqualCi}}},
			expecting: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.tc.Prepare(); err != nil {
				t.Errorf("unable to prepare header matcher: %v", err)
			}

			checker := MakeChecker(tt.mh, nil)

			j, _ := json.Marshal(tt.tc)
			if checker(string(j)) != tt.expecting {
				t.Errorf("did not match (expecting: %v)", tt.expecting)
			}
		})
	}
}

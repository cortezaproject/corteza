package service

import (
	"net/mail"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/system/types"
)

func Test_mailProcMessage(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantOut *types.MailMessage
		wantErr bool
	}{
		// TODO: Add test cases.
		{name: "basics",
			input: `
From: <sender@testing.cortezaproject.org>
To: <rcpt@testing.cortezaproject.org>
Subject: Customer service contact info
Message-ID: <1234@local.machine.example>

Ola Corteza!
`,
			wantOut: &types.MailMessage{
				Date: time.Time{},
				Header: types.MailMessageHeader{
					From: []*mail.Address{{Address: "sender@testing.cortezaproject.org"}},
					To:   []*mail.Address{{Address: "rcpt@testing.cortezaproject.org"}},

					Raw: map[string][]string{
						"From":       []string{"<sender@testing.cortezaproject.org>"},
						"To":         []string{"<rcpt@testing.cortezaproject.org>"},
						"Subject":    []string{"Customer service contact info"},
						"Message-Id": []string{"<1234@local.machine.example>"},
					},
				},
				RawBody: []byte(`Ola Corteza!`),
			}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := strings.NewReader(strings.TrimSpace(tt.input))
			gotOut, err := mailProcMessage(input)

			if (err != nil) != tt.wantErr {
				t.Errorf("mailProcMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("mailProcMessage() \ngotOut: %v, \n  want: %v", gotOut, tt.wantOut)
			}
		})
	}
}

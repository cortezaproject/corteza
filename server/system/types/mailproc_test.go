package types

import (
	"net/mail"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_mailProcMessage(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantOut *MailMessage
		wantErr bool
	}{
		{name: "basics",
			input: `
From: <sender@testing.example.com>
To: <rcpt@testing.example.com>
Subject: Customer service contact info
Message-ID: <1234@local.machine.example>

Ola Human!
`,
			wantOut: &MailMessage{
				Date:    time.Time{},
				Subject: "Customer service contact info",
				Header: MailMessageHeader{
					From: []*mail.Address{{Address: "sender@testing.example.com"}},
					To:   []*mail.Address{{Address: "rcpt@testing.example.com"}},

					Raw: map[string][]string{
						"From":       []string{"<sender@testing.example.com>"},
						"To":         []string{"<rcpt@testing.example.com>"},
						"Subject":    []string{"Customer service contact info"},
						"Message-Id": []string{"<1234@local.machine.example>"},
					},
				},
				RawBody: []byte(`Ola Human!`),
			}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := strings.NewReader(strings.TrimSpace(tt.input))
			gotOut, err := NewMailMessage(input)

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

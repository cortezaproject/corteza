package healthcheck

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Healthy(t *testing.T) {
	tests := []struct {
		name    string
		checks  []*check
		healthy bool
		string  string
	}{
		{
			"should be healthy with handle for stringer output",
			[]*check{{func(ctx context.Context) error { return nil }, &Meta{Label: "check01"}}},
			true,
			"PASS check01\n",
		},
		{
			"should handle multiple healthy checks",
			[]*check{
				{func(ctx context.Context) error { return nil }, &Meta{Label: "check01"}},
				{func(ctx context.Context) error { return nil }, &Meta{Label: "check02"}},
			},
			true,
			"PASS check01\nPASS check02\n",
		},
		{
			"should handle healthy and unhealthy checks",
			[]*check{
				{func(ctx context.Context) error { return nil }, &Meta{Label: "check01"}},
				{func(ctx context.Context) error { return fmt.Errorf("x") }, &Meta{Label: "check02"}},
				{func(ctx context.Context) error { return nil }, &Meta{Label: "check03"}},
			},
			false,
			"PASS check01\nFAIL check02: x\nPASS check03\n",
		},
		{
			"should handle labels",
			[]*check{
				{func(ctx context.Context) error { return nil }, &Meta{Label: "check01"}},
				{func(ctx context.Context) error { return nil }, &Meta{Label: "Pretty check"}},
			},
			true,
			"PASS check01\nPASS Pretty check\n",
		},
		{
			"should handle empty check list",
			[]*check{},
			true,
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := assert.New(t)
			r := (&checks{cc: tt.checks}).Run(context.Background())
			a.Equal(tt.healthy, r.Healthy(), "healthy result failed")
			a.Equal(tt.string, r.String(), "stringer output match failed")
		})
	}
}

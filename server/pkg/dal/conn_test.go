package dal

import (
	"context"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func Test_checkIdent(t *testing.T) {
	tests := []struct {
		name  string
		ident string
		rr    []*regexp.Regexp
		want  bool
	}{
		{
			name:  "empty",
			ident: "",
			rr:    []*regexp.Regexp{},
			want:  true,
		},
		{
			name:  "one",
			ident: "foo",
			rr:    []*regexp.Regexp{regexp.MustCompile("foo")},
			want:  true,
		},
		{
			name:  "false",
			ident: "foo",
			rr:    []*regexp.Regexp{regexp.MustCompile("bar")},
			want:  false,
		},
		{
			name:  "two",
			ident: "bar",
			rr:    []*regexp.Regexp{regexp.MustCompile("foo"), regexp.MustCompile("bar")},
			want:  true,
		},
		{
			name:  "two failed",
			ident: "foo",
			rr:    []*regexp.Regexp{regexp.MustCompile("bar"), regexp.MustCompile("baz")},
			want:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkIdent(tt.ident, tt.rr...); got != tt.want {
				t.Errorf("checkIdent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getConnection_nilUnderlyingConnection(t *testing.T) {
	const (
		connID  = uint64(100)
		modelID = uint64(200)
	)

	svc, err := New(zap.NewNop(), false)
	require.NoError(t, err)

	svc.addConnection(&ConnectionWrap{ID: connID})
	svc.addModelToRegistry(&Model{ConnectionID: connID, ResourceID: modelID, Ident: "t"}, false)

	_, err = svc.Search(context.Background(), ModelRef{ConnectionID: connID, ResourceID: modelID}, nil, nil)
	require.ErrorContains(t, err, "not available")
}

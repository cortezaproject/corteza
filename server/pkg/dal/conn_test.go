package dal

import (
	"regexp"
	"testing"
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

package external

import (
	"net/url"
	"reflect"
	"testing"
)

func Test_parseExternalProviderUrl(t *testing.T) {
	mustParseURL := func(in string) *url.URL {
		r, err := url.Parse(in)
		if err != nil {
			panic(err)
		}

		return r
	}

	type args struct {
		in string
	}

	tests := []struct {
		name    string
		args    args
		wantP   *url.URL
		wantErr bool
	}{
		{
			"happy-path",
			args{"https://foo.bar"},
			mustParseURL("https://foo.bar"),
			false,
		},
		{
			"bad input",
			args{":\\"},
			nil,
			true,
		},
		{
			"add schema",
			args{"cortezaproject.org"},
			mustParseURL("https://cortezaproject.org"),
			false,
		},
		{
			"add schema and remove well-known",
			args{"cortezaproject.org/some-subdir/" + WellKnown},
			mustParseURL("https://cortezaproject.org/some-subdir/"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotP, err := parseExternalProviderUrl(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseExternalProviderUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotP, tt.wantP) {
				t.Errorf("parseExternalProviderUrl() gotP = %v, want %v", gotP, tt.wantP)
			}
		})
	}
}

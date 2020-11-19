package label

import (
	"reflect"
	"testing"
)

func TestChanged(t *testing.T) {
	tests := []struct {
		name string
		old  map[string]string
		new  map[string]string
		want bool
	}{
		{
			"2x nil",
			nil,
			nil,
			false,
		},
		{
			"2x empty",
			map[string]string{},
			map[string]string{},
			false,
		},
		{
			"nil & empty",
			nil,
			map[string]string{},
			false,
		},
		{
			"same",
			map[string]string{"a": "a"},
			map[string]string{"a": "a"},
			false,
		},
		{
			"diff1",
			map[string]string{"a": "a"},
			map[string]string{"a": "b"},
			true,
		},
		{
			"diff2",
			map[string]string{"a": "b"},
			map[string]string{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Changed(tt.old, tt.new); got != tt.want {
				t.Errorf("Changed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseStrings(t *testing.T) {
	tests := []struct {
		name    string
		labels  []string
		want    map[string]string
		wantErr bool
	}{
		{
			"set of pairs",
			[]string{"aa=b"},
			map[string]string{"aa": "b"},
			false,
		},
		{
			"empty json",
			[]string{`{}`},
			map[string]string{},
			false,
		},
		{
			"json",
			[]string{`{"aa":"b"}`},
			map[string]string{"aa": "b"},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStrings(tt.labels)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseStrings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseStrings() got = %v, want %v", got, tt.want)
			}
		})
	}
}

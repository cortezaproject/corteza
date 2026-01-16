package label

import (
	"reflect"
	"testing"

	"github.com/cortezaproject/corteza/server/pkg/label/types"
)

func TestChanged(t *testing.T) {
	tests := []struct {
		name string
		old  map[string]types.LabelValue
		new  map[string]types.LabelValue
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
			map[string]types.LabelValue{},
			map[string]types.LabelValue{},
			false,
		},
		{
			"nil & empty",
			nil,
			map[string]types.LabelValue{},
			false,
		},
		{
			"same single value",
			map[string]types.LabelValue{"a": {Val: "a"}},
			map[string]types.LabelValue{"a": {Val: "a"}},
			false,
		},
		{
			"same multi value",
			map[string]types.LabelValue{"a": {Values: []string{"1", "2"}}},
			map[string]types.LabelValue{"a": {Values: []string{"1", "2"}}},
			false,
		},
		{
			"diff single value",
			map[string]types.LabelValue{"a": {Val: "a"}},
			map[string]types.LabelValue{"a": {Val: "b"}},
			true,
		},
		{
			"diff multi value",
			map[string]types.LabelValue{"a": {Values: []string{"1", "2"}}},
			map[string]types.LabelValue{"a": {Values: []string{"1", "3"}}},
			true,
		},
		{
			"diff key removed",
			map[string]types.LabelValue{"a": {Val: "b"}},
			map[string]types.LabelValue{},
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
		want    map[string]types.LabelValue
		wantErr bool
	}{
		{
			"set of pairs",
			[]string{"aa=b"},
			map[string]types.LabelValue{"aa": {Val: "b"}},
			false,
		},
		{
			"empty json",
			[]string{`{}`},
			map[string]types.LabelValue{},
			false,
		},
		{
			"json",
			[]string{`{"aa":"b"}`},
			map[string]types.LabelValue{"aa": {Val: "b"}},
			false,
		},
		{
			"value with colons and slashes (single)",
			[]string{"ref_namespace=corteza::compose:namespace/123"},
			map[string]types.LabelValue{"ref_namespace": {Val: "corteza::compose:namespace/123"}},
			false,
		},
		{
			"value with colons and slashes (array)",
			[]string{`ref_namespace=["corteza::compose:namespace/123","corteza::compose:namespace/456"]`},
			map[string]types.LabelValue{"ref_namespace": {Values: []string{"corteza::compose:namespace/123", "corteza::compose:namespace/456"}}},
			false,
		},
		{
			"multiple reference labels",
			[]string{
				"ref_namespace=corteza::compose:namespace/123",
				"ref_module=corteza::compose:module/123/456",
			},
			map[string]types.LabelValue{
				"ref_namespace": {Val: "corteza::compose:namespace/123"},
				"ref_module":    {Val: "corteza::compose:module/123/456"},
			},
			false,
		},
		{
			"quoted string value (strips quotes)",
			[]string{`ref_module="corteza::compose:module/1239888/408"`},
			map[string]types.LabelValue{"ref_module": {Val: "corteza::compose:module/1239888/408"}},
			false,
		},
		{
			"mixed: array and quoted string",
			[]string{
				`ref_namespace=["corteza::compose:namespace/123"]`,
				`ref_module="corteza::compose:module/1239888/408"`,
			},
			map[string]types.LabelValue{
				"ref_namespace": {Values: []string{"corteza::compose:namespace/123"}},
				"ref_module":    {Val: "corteza::compose:module/1239888/408"},
			},
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

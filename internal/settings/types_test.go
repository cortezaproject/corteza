package settings

import (
	"testing"

	"github.com/jmoiron/sqlx/types"

	"github.com/cortezaproject/corteza-server/internal/test"
)

func TestKV_Bool(t *testing.T) {
	type args struct {
		k string
	}
	tests := []struct {
		name  string
		kv    KV
		args  args
		wantV bool
	}{
		{
			name:  "True value should return true",
			kv:    KV{"true-value": types.JSONText(`true`)},
			args:  args{k: "true-value"},
			wantV: true,
		},
		{
			name:  "Null value should return false",
			kv:    KV{"null-value": types.JSONText(`null`)},
			args:  args{k: "null-value"},
			wantV: false,
		},
		{
			name:  "Unexisting value should return false",
			kv:    KV{},
			args:  args{k: "unexisting"},
			wantV: false,
		},
		{
			name:  "Invalid KV should return false",
			kv:    nil,
			args:  args{k: "invalid-kv"},
			wantV: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotV := tt.kv.Bool(tt.args.k); gotV != tt.wantV {
				t.Errorf("KV.Bool() = %v, want %v", gotV, tt.wantV)
			}
		})
	}
}

func TestSettingValueAsString(t *testing.T) {
	test.NoError(t, (&Value{}).SetRawValue(`"string"`), "unable to set value as string")
	test.NoError(t, (&Value{}).SetRawValue(`false`), "unable to set value as string")
	test.NoError(t, (&Value{}).SetRawValue(`null`), "unable to set value as string")
	test.NoError(t, (&Value{}).SetRawValue(`42`), "unable to set value as string")
	test.NoError(t, (&Value{}).SetRawValue(`3.14`), "unable to set value as string")
	test.Error(t, (&Value{}).SetRawValue(`error`), "expecting error when not setting JSON")
}

func TestValueSet_Upsert(t *testing.T) {
	var vv = ValueSet{}

	test.Assert(t, len(vv) == 0, "expecting length to be 0")
	vv.Replace(&Value{Name: "name"})
	test.Assert(t, len(vv) == 1, "expecting length to be 1")
	vv.Replace(&Value{Name: "name", Value: []byte("42")})
	test.Assert(t, len(vv) == 1, "expecting length to be 1")
	test.Assert(t, string(vv[0].Value) == "42", "expecting value to be 42")
}

func TestValueSet_Changed(t *testing.T) {
	var (
		// make string value
		msv = func(n, v string) *Value {
			o := &Value{Name: n}
			_ = o.SetValue(v)
			return o
		}

		org = ValueSet{msv("a", "a1"), msv("b", "b1"), msv("d", "d1")}
		inp = ValueSet{msv("a", "a2"), msv("c", "c1"), msv("d", "d1")}

		out ValueSet
	)

	out = org.Changed(inp)

	test.Assert(t, len(out) == 2, "expecting length to be 2, got %d", len(out))
	test.Assert(t, out.First("a").String() == "a2", "expecting 'a' to have 'a2' value")
	test.Assert(t, out.First("b").String() == "", "expecting 'b' to be missing")
	test.Assert(t, out.First("c").String() == "c1", "expecting 'c' to have 'c1' value")
	test.Assert(t, out.First("d").String() == "", "expecting 'd' to be missing")
}

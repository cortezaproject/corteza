package settings

import (
	"testing"

	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/require"
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
	var req = require.New(t)

	req.NoError((&Value{}).SetRawValue(`"string"`), "unable to set value as string")
	req.NoError((&Value{}).SetRawValue(`false`), "unable to set value as string")
	req.NoError((&Value{}).SetRawValue(`null`), "unable to set value as string")
	req.NoError((&Value{}).SetRawValue(`42`), "unable to set value as string")
	req.NoError((&Value{}).SetRawValue(`3.14`), "unable to set value as string")
	req.Error((&Value{}).SetRawValue(`error`), "expecting error when not setting JSON")
}

func TestValueSet_Upsert(t *testing.T) {
	var (
		req = require.New(t)
		vv  = ValueSet{}
	)

	req.Len(vv, 0)

	vv.Replace(&Value{Name: "name"})
	req.Len(vv, 1)

	vv.Replace(&Value{Name: "name", Value: []byte("42")})
	req.Len(vv, 1)
	req.Equal("42", string(vv[0].Value))
}

func TestValueSet_Changed(t *testing.T) {
	var (
		req = require.New(t)

		// make string value
		msv = func(n, v string) *Value {
			o := &Value{Name: n}
			_ = o.SetValue(v)
			return o
		}

		// make bool value
		mbv = func(n string, v bool) *Value {
			o := &Value{Name: n}
			_ = o.SetValue(v)
			return o
		}

		org = ValueSet{msv("a", "a1"), msv("b", "b1"), msv("d", "d1"), mbv("bool", true)}
		inp = ValueSet{msv("a", "a2"), msv("c", "c1"), msv("d", "d1"), mbv("bool", false)}

		out ValueSet
	)

	out = org.Changed(inp)

	req.Len(out, 3)
	req.Equal("a2", out.First("a").String())
	req.Equal("", out.First("b").String())
	req.Equal("c1", out.First("c").String())
	req.Equal("", out.First("d").String())
	req.Equal(false, out.First("bool").Bool())
}

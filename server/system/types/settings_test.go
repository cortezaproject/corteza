package types

import (
	"testing"

	"github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/require"
)

func TestSettingsKV_Bool(t *testing.T) {
	type args struct {
		k string
	}
	tests := []struct {
		name  string
		kv    SettingsKV
		args  args
		wantV bool
	}{
		{
			name:  "True value should return true",
			kv:    SettingsKV{"true-value": types.JSONText(`true`)},
			args:  args{k: "true-value"},
			wantV: true,
		},
		{
			name:  "Null value should return false",
			kv:    SettingsKV{"null-value": types.JSONText(`null`)},
			args:  args{k: "null-value"},
			wantV: false,
		},
		{
			name:  "Unexisting value should return false",
			kv:    SettingsKV{},
			args:  args{k: "unexisting"},
			wantV: false,
		},
		{
			name:  "Invalid SettingsKV should return false",
			kv:    nil,
			args:  args{k: "invalid-kv"},
			wantV: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotV := tt.kv.Bool(tt.args.k); gotV != tt.wantV {
				t.Errorf("SettingsKV.Bool() = %v, want %v", gotV, tt.wantV)
			}
		})
	}
}

func TestSettingValueAsString(t *testing.T) {
	var req = require.New(t)

	req.NoError((&SettingValue{}).SetRawSetting(`"string"`), "unable to set value as string")
	req.NoError((&SettingValue{}).SetRawSetting(`false`), "unable to set value as string")
	req.NoError((&SettingValue{}).SetRawSetting(`null`), "unable to set value as string")
	req.NoError((&SettingValue{}).SetRawSetting(`42`), "unable to set value as string")
	req.NoError((&SettingValue{}).SetRawSetting(`3.14`), "unable to set value as string")
	req.Error((&SettingValue{}).SetRawSetting(`error`), "expecting error when not setting JSON")
}

func TestSettingValueSet_Upsert(t *testing.T) {
	var (
		req = require.New(t)
		vv  = SettingValueSet{}
	)

	req.Len(vv, 0)

	vv.Replace(&SettingValue{Name: "name"})
	req.Len(vv, 1)

	vv.Replace(&SettingValue{Name: "name", Value: []byte("42")})
	req.Len(vv, 1)
	req.Equal("42", string(vv[0].Value))
}

func TestSettingValueSet_Changed(t *testing.T) {
	var (
		req = require.New(t)

		// make string value
		msv = func(n, v string) *SettingValue {
			o := &SettingValue{Name: n}
			_ = o.SetSetting(v)
			return o
		}

		// make bool value
		mbv = func(n string, v bool) *SettingValue {
			o := &SettingValue{Name: n}
			_ = o.SetSetting(v)
			return o
		}

		org = SettingValueSet{msv("a", "a1"), msv("b", "b1"), msv("d", "d1"), mbv("bool", true)}
		inp = SettingValueSet{msv("a", "a2"), msv("c", "c1"), msv("d", "d1"), mbv("bool", false)}

		out SettingValueSet
	)

	out = org.Changed(inp)

	req.Len(out, 3)
	req.Equal("a2", out.First("a").String())
	req.Equal("", out.First("b").String())
	req.Equal("c1", out.First("c").String())
	req.Equal("", out.First("d").String())
	req.Equal(false, out.First("bool").Bool())
}

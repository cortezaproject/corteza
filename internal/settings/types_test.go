package settings

import (
	"testing"

	"github.com/jmoiron/sqlx/types"

	"github.com/crusttech/crust/internal/test"
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
	test.NoError(t, (&Value{}).SetValueAsString(`"string"`), "unable to set value as string")
	test.NoError(t, (&Value{}).SetValueAsString(`false`), "unable to set value as string")
	test.NoError(t, (&Value{}).SetValueAsString(`null`), "unable to set value as string")
	test.NoError(t, (&Value{}).SetValueAsString(`42`), "unable to set value as string")
	test.NoError(t, (&Value{}).SetValueAsString(`3.14`), "unable to set value as string")
	test.Error(t, (&Value{}).SetValueAsString(`error`), "expecting error when not setting JSON")
}

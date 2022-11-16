package values

import (
	"context"
	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/locale"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"time"
)

func makeLocaleService() localeService {
	return locale.Global()
}

func Test_validator_vDatetime(t *testing.T) {
	var (
		vldtr = validator{
			now: func() time.Time {
				t, err := time.Parse(time.RFC3339, "2020-02-20T10:10:10Z")
				if err != nil {
					panic(err)
				}
				return t
			},
			localeSvc: makeLocaleService(),
		}
		tests = []struct {
			name string
			val  string
			opt  types.ModuleFieldOptions
			want []types.RecordValueError
		}{
			{
				name: "unparsable",
				val:  "unparsable",
				want: e2s(types.RecordValueError{Kind: "invalidValue", Meta: map[string]interface{}{"field": "", "value": "unparsable"}, Message: "record-field.errors.invalidValue"}),
			},
			{
				name: "valid datetime value",
				val:  "2020-02-20T10:10:10Z",
			},
			{
				name: "valid date value",
				val:  "2020-02-20",
				opt:  types.ModuleFieldOptions{fieldOpt_Datetime_onlyDate: true},
			},
			{
				name: "valid time value",
				val:  "10:10:10",
				opt:  types.ModuleFieldOptions{fieldOpt_Datetime_onlyTime: true},
			},
			{
				name: "valid future datetime value",
				val:  "2021-02-20T10:10:10Z",
				opt:  types.ModuleFieldOptions{fieldOpt_Datetime_onlyFutureValues: true},
			},
			{
				name: "valid future date value",
				val:  "2021-02-20",
				opt:  types.ModuleFieldOptions{fieldOpt_Datetime_onlyFutureValues: true, fieldOpt_Datetime_onlyDate: true},
			},
			{
				name: "valid future time value",
				val:  "11:10:10",
				opt:  types.ModuleFieldOptions{fieldOpt_Datetime_onlyFutureValues: true, fieldOpt_Datetime_onlyTime: true},
			},
			{
				name: "valid past datetime value",
				val:  "2019-02-20T10:10:10Z",
				opt:  types.ModuleFieldOptions{fieldOpt_Datetime_onlyPastValues: true},
			},
			{
				name: "valid past date value",
				val:  "2019-02-20",
				opt:  types.ModuleFieldOptions{fieldOpt_Datetime_onlyPastValues: true, fieldOpt_Datetime_onlyDate: true},
			},
			{
				name: "valid past time value",
				val:  "09:10:10",
				opt:  types.ModuleFieldOptions{fieldOpt_Datetime_onlyPastValues: true, fieldOpt_Datetime_onlyTime: true},
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := vldtr.vDatetime(context.Background(), &types.RecordValue{Value: tt.val}, &types.ModuleField{Options: tt.opt}, nil, nil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("vDatetime() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_validator_vNumber(t *testing.T) {
	var (
		vldtr = validator{
			localeSvc: makeLocaleService(),
		}
		tests = []struct {
			name string
			val  string
			opt  types.ModuleFieldOptions
			want []types.RecordValueError
		}{
			{
				name: "unparsable",
				val:  "unparsable",
				want: e2s(types.RecordValueError{Kind: "invalidValue", Meta: map[string]interface{}{"field": "", "value": "unparsable"}, Message: "record-field.errors.invalidValue"}),
			},
			{
				name: "valid number value",
				val:  "42",
			},
			{
				name: "valid number value",
				val:  "42.123",
			},
			{
				name: "valid number value",
				val:  "412412412.322894325892365",
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := vldtr.vNumber(context.Background(), &types.RecordValue{Value: tt.val}, &types.ModuleField{Options: tt.opt}, nil, nil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("vNumber() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_validator_vUrl(t *testing.T) {
	var (
		vldtr = validator{
			localeSvc: makeLocaleService(),
		}
		tests = []struct {
			name string
			val  string
			opt  types.ModuleFieldOptions
			want []types.RecordValueError
		}{
			{
				name: "invalid-url",
				val:  "invalid-url",
				want: e2s(types.RecordValueError{Kind: "invalidValue", Meta: map[string]interface{}{"field": "", "value": "invalid-url"}, Message: "record-field.errors.invalidValue"}),
			},
			{
				name: "valid value",
				val:  "https://crust.tech/",
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := vldtr.vUrl(context.Background(), &types.RecordValue{Value: tt.val}, &types.ModuleField{Options: tt.opt}, nil, nil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("vUrl() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_validator_vEmail(t *testing.T) {
	var (
		vldtr = validator{
			localeSvc: makeLocaleService(),
		}
		tests = []struct {
			name string
			val  string
			opt  types.ModuleFieldOptions
			want []types.RecordValueError
		}{
			{
				name: "unparsable",
				val:  "un pars able",
				want: e2s(types.RecordValueError{Kind: "invalidValue", Meta: map[string]interface{}{"field": "", "value": "un pars able"}, Message: "record-field.errors.invalidValue"}),
			},
			{
				name: "valid value",
				val:  "qa@crust.tech",
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := vldtr.vEmail(context.Background(), &types.RecordValue{Value: tt.val}, &types.ModuleField{Options: tt.opt}, nil, nil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("vEmail() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_validator_vSelect(t *testing.T) {
	var (
		vldtr = validator{
			localeSvc: makeLocaleService(),
		}
		tests = []struct {
			name string
			val  string
			opt  types.ModuleFieldOptions
			want []types.RecordValueError
		}{
			{
				name: "unparsable",
				val:  "dummy",
			},
			{
				name: "valid value",
				val:  "crust",
				opt:  types.ModuleFieldOptions{"options": []string{"crust", "corteza"}},
			},
			{
				name: "valid value",
				val:  "the rest",
				opt:  types.ModuleFieldOptions{"options": []string{"crust", "corteza"}},
				want: e2s(types.RecordValueError{Kind: "invalidValue", Meta: map[string]interface{}{"field": "", "value": "the rest"}, Message: "record-field.errors.invalidValue"}),
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := vldtr.vSelect(context.Background(), &types.RecordValue{Value: tt.val}, &types.ModuleField{Options: tt.opt}, nil, nil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("vSelect() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func Test_validator_customExpr(t *testing.T) {
	var (
		vldtr = validator{}
		m     = &types.Module{}
		f     = &types.ModuleField{Name: "num", Kind: "Number"}
		r     = &types.Record{}
	)

	f.Expressions.Validators = []types.ModuleFieldValidator{
		{Test: "value < 5", Error: "value is lower than 5"},
	}
	m.Fields = append(m.Fields, f)
	r.Values = r.Values.Replace("num", "1")

	rve := vldtr.Run(context.Background(), nil, m, r)
	require.False(t, rve.IsValid())

	r.Values = r.Values.Replace("num", "10")
	rve = vldtr.Run(context.Background(), nil, m, r)
	require.True(t, rve.IsValid())
}

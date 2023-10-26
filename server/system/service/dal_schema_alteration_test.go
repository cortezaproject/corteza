package service

import (
	"reflect"
	"testing"

	"github.com/cortezaproject/corteza/server/system/types"
)

func TestAppliableAlterations(t *testing.T) {
	n := now()
	tcc := []struct {
		name string
		in   types.DalSchemaAlterationSet
		out  types.DalSchemaAlterationSet
	}{
		{
			name: "empty",
			in:   types.DalSchemaAlterationSet{},
			out:  types.DalSchemaAlterationSet{},
		},
		{
			name: "filter out completed",
			in: types.DalSchemaAlterationSet{
				{},
				{CompletedAt: n},
			},
			out: types.DalSchemaAlterationSet{
				{},
			},
		},
		{
			name: "filter out dismissed",
			in: types.DalSchemaAlterationSet{
				{},
				{DismissedAt: n},
			},
			out: types.DalSchemaAlterationSet{
				{},
			},
		},
		{
			name: "filter out when dependency on missing",
			in: types.DalSchemaAlterationSet{
				{},
				{DependsOn: 1},
			},
			out: types.DalSchemaAlterationSet{
				{},
			},
		},
		{
			name: "include when dependency present",
			in: types.DalSchemaAlterationSet{
				{ID: 1},
				{DependsOn: 1},
			},
			out: types.DalSchemaAlterationSet{
				{ID: 1},
				{DependsOn: 1},
			},
		},
		{
			name: "include when dependency present but excluded",
			in: types.DalSchemaAlterationSet{
				{ID: 1, CompletedAt: n, DismissedAt: n},
				{DependsOn: 1},
			},
			out: types.DalSchemaAlterationSet{
				{DependsOn: 1},
			},
		},
	}

	svc := dalSchemaAlteration{}
	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			got := svc.appliableAlterations(tc.in...)
			if !reflect.DeepEqual(got, tc.out) {
				t.Errorf("appliableAlterations() = %v, want %v", got, tc.out)
			}
		})
	}
}

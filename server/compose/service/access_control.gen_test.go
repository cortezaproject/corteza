package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/rbac"
	"github.com/stretchr/testify/require"
)

func TestAccessControl_ResourceLoader(t *testing.T) {
	svc := accessControl{}

	// Has wildcard resources
	testCases := []struct {
		resource string
		expected rbac.Resource
		err      error
	}{
		{
			resource: "corteza::compose:chart/1/*",
			expected: rbac.NewResource(types.ChartRbacResource(1, 0)),
			err:      nil,
		},
		{
			resource: "corteza::compose:chart/*/*",
			expected: rbac.NewResource(types.ChartRbacResource(0, 0)),
			err:      nil,
		},
		{
			resource: "corteza::compose:module/3/*",
			expected: rbac.NewResource(types.ModuleRbacResource(3, 0)),
			err:      nil,
		},
		{
			resource: "corteza::compose:module/*/*",
			expected: rbac.NewResource(types.ModuleRbacResource(0, 0)),
			err:      nil,
		},
		{
			resource: "corteza::compose:module-field/5/*/*",
			expected: rbac.NewResource(types.ModuleFieldRbacResource(5, 0, 0)),
			err:      nil,
		},
		{
			resource: "corteza::compose:module-field/*/*/*",
			expected: rbac.NewResource(types.ModuleFieldRbacResource(0, 0, 0)),
			err:      nil,
		},
		{
			resource: "corteza::compose:namespace/*",
			expected: rbac.NewResource(types.NamespaceRbacResource(0)),
			err:      nil,
		},
		{
			resource: "corteza::compose:page/9/*",
			expected: rbac.NewResource(types.PageRbacResource(9, 0)),
			err:      nil,
		},
		{
			resource: "corteza::compose:page/*/*",
			expected: rbac.NewResource(types.PageRbacResource(0, 0)),
			err:      nil,
		},
		{
			resource: "corteza::compose:page-layout/11/*/*",
			expected: rbac.NewResource(types.PageLayoutRbacResource(11, 0, 0)),
			err:      nil,
		},
		{
			resource: "corteza::compose:page-layout/*/*/*",
			expected: rbac.NewResource(types.PageLayoutRbacResource(0, 0, 0)),
			err:      nil,
		},
		{
			resource: "corteza::compose:record/14/*/*",
			expected: rbac.NewResource(types.RecordRbacResource(14, 0, 0)),
			err:      nil,
		},
		{
			resource: "corteza::compose:record/*/*/*",
			expected: rbac.NewResource(types.RecordRbacResource(0, 0, 0)),
			err:      nil,
		},
		{
			resource: "corteza::compose",
			expected: &types.Component{},
			err:      nil,
		},
		{
			resource: "unknown_resource_type:17",
			expected: nil,
			err:      fmt.Errorf("unknown resource type %q", "unknown_resource_type:17"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.resource, func(t *testing.T) {
			res, err := svc.resourceLoader(context.Background(), tc.resource)

			require.Equal(t, tc.expected, res)
			require.Equal(t, tc.err, err)
		})
	}
}

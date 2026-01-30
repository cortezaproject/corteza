package service

import (
	"context"

	"github.com/cortezaproject/corteza/server/automation/types"
	"github.com/cortezaproject/corteza/server/pkg/expr"
	"github.com/davecgh/go-spew/spew"
)

// @note this will be generated; for now it's just built manually

func populateConstructLibrary() {
	defaultConstructRegistry.triggers = []types.ConstructTrigger{{
		ResourceType: "system",
		EventType:    "onManual",
	}}

	defaultConstructRegistry.functions = []types.ConstructFunction{{
		Ref:    "usersLookup",
		Kind:   "function",
		Labels: map[string]string{"users": "step,workflow"},
		Meta: &types.ConstructFunctionMeta{
			Short:       "User lookup",
			Description: "Find specific user by ID, handle or string",
		},

		Parameters: types.ParamSet{{
			Name:  "lookup",
			Types: []string{"ID", "Handle", "String", "User"}, Required: true,
		}},

		Results: types.ParamSet{{
			Name:  "user",
			Types: []string{"User"},
		}},

		Segments: []types.ConstructSegment{{
			Meta: types.ConstructSegmentMeta{},
			Sections: []types.ConstructSection{{
				Meta: types.ConstructSectionMeta{},
				Elements: []types.SectionElement{{
					Input: types.SectionElementInput{
						Type:     "UserSelector",
						Label:    "User",
						Argument: "lookup",
					},
				}},
			}},
		}},

		Handler: func(ctx context.Context, in *expr.Vars) (out *expr.Vars, err error) {
			spew.Dump("users lookup")
			return
		},
	}}

}

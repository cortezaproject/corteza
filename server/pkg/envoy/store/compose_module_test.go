package store

import (
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/compose/types"
	st "github.com/jmoiron/sqlx/types"
	"github.com/stretchr/testify/require"
)

func TestComposeModule_Merger(t *testing.T) {
	req := require.New(t)

	now := time.Time{}
	nowP := &time.Time{}

	empty := &types.Module{}
	full := &types.Module{
		Handle:      "handle",
		Name:        "name",
		Meta:        st.JSONText{},
		NamespaceID: 1,

		CreatedAt: now,
		UpdatedAt: nowP,
		DeletedAt: nowP,
	}

	t.Run("merge on empty", func(t *testing.T) {
		c := mergeComposeModule(empty, full)
		req.Equal("name", c.Name)
		req.Equal("handle", c.Handle)
		req.NotNil(c.Meta)
		req.Equal(uint64(1), c.NamespaceID)
		req.Equal(now, c.CreatedAt)
		req.Equal(nowP, c.UpdatedAt)
		req.Equal(nowP, c.DeletedAt)
	})

	t.Run("merge with empty", func(t *testing.T) {
		c := mergeComposeModule(full, empty)
		req.Equal("name", c.Name)
		req.Equal("handle", c.Handle)
		req.NotNil(c.Meta)
		req.Equal(uint64(0), c.NamespaceID)
		req.Equal(now, c.CreatedAt)
		req.Equal(nowP, c.UpdatedAt)
		req.Equal(nowP, c.DeletedAt)
	})
}

func TestComposeModuleField_Merger(t *testing.T) {
	req := require.New(t)

	now := time.Time{}
	nowP := &time.Time{}

	empty := types.ModuleFieldSet{}
	full := types.ModuleFieldSet{
		&types.ModuleField{
			ModuleID:     1,
			Place:        2,
			Kind:         "kind",
			Name:         "name",
			Label:        "label",
			Options:      types.ModuleFieldOptions{},
			DefaultValue: types.RecordValueSet{},
			Expressions:  types.ModuleFieldExpr{},
			CreatedAt:    now,
			UpdatedAt:    nowP,
			DeletedAt:    nowP,
		},
	}

	t.Run("merge on empty", func(t *testing.T) {
		cc := mergeComposeModuleFields(empty, full)
		req.Len(cc, 1)
		c := cc[0]

		req.Equal(uint64(1), c.ModuleID)
		req.Equal(2, c.Place)
		req.Equal("kind", c.Kind)
		req.Equal("name", c.Name)
		req.Equal("label", c.Label)
		req.NotNil(c.Options)
		req.NotNil(c.DefaultValue)
		req.NotNil(c.Expressions)
		req.Equal(now, c.CreatedAt)
		req.Equal(nowP, c.UpdatedAt)
		req.Equal(nowP, c.DeletedAt)
	})

	t.Run("merge with empty", func(t *testing.T) {
		cc := mergeComposeModuleFields(full, empty)
		req.Len(cc, 1)
		c := cc[0]

		req.Equal(uint64(1), c.ModuleID)
		req.Equal(2, c.Place)
		req.Equal("kind", c.Kind)
		req.Equal("name", c.Name)
		req.Equal("label", c.Label)
		req.NotNil(c.Options)
		req.NotNil(c.DefaultValue)
		req.NotNil(c.Expressions)
		req.Equal(now, c.CreatedAt)
		req.Equal(nowP, c.UpdatedAt)
		req.Equal(nowP, c.DeletedAt)
	})
}

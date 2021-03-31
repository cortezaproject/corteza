package envoy

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza-server/compose/types"
	su "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/store"
)

func sTestComposeNamespace(ctx context.Context, t *testing.T, s store.Storer, pfx string) *types.Namespace {
	ns := &types.Namespace{
		ID:      su.NextID(),
		Name:    pfx + " namespace",
		Slug:    pfx + "_namespace",
		Enabled: true,
		Meta: types.NamespaceMeta{
			Subtitle:    "subtitle",
			Description: "description",
		},
		CreatedAt: createdAt,
		UpdatedAt: &updatedAt,
	}

	err := store.CreateComposeNamespace(ctx, s, ns)
	if err != nil {
		t.Fatal(err)
	}

	return ns
}

func sTestComposeModule(ctx context.Context, t *testing.T, s store.Storer, nsID uint64, pfx string) *types.Module {
	mID := su.NextID()

	mod := &types.Module{
		ID:          mID,
		Name:        pfx + " module",
		Handle:      pfx + "_module",
		NamespaceID: nsID,
		Fields: types.ModuleFieldSet{
			&types.ModuleField{
				ID:        su.NextID(),
				ModuleID:  mID,
				Place:     0,
				Kind:      "String",
				Name:      "module_field_string",
				Label:     "module field string",
				Private:   true,
				Required:  true,
				Visible:   true,
				Multi:     true,
				CreatedAt: createdAt,
				UpdatedAt: &updatedAt,
				Options: types.ModuleFieldOptions{
					"opt1": "opt_value_1",
				},
			},
			&types.ModuleField{
				ID:        su.NextID(),
				ModuleID:  mID,
				Place:     1,
				Kind:      "Number",
				Name:      "module_field_number",
				Label:     "module field number",
				Private:   false,
				Required:  false,
				Visible:   false,
				Multi:     false,
				CreatedAt: createdAt,
				UpdatedAt: &updatedAt,
				Options: types.ModuleFieldOptions{
					"opt1": "opt_value_1",
				},
			},
		},

		CreatedAt: createdAt,
		UpdatedAt: &updatedAt,
	}

	err := store.CreateComposeModule(ctx, s, mod)
	if err != nil {
		t.Fatal(err)
	}
	err = store.CreateComposeModuleField(ctx, s, mod.Fields...)
	if err != nil {
		t.Fatal(err)
	}

	return mod
}

func sTestComposeModuleFull(ctx context.Context, s store.Storer, t *testing.T, nsID uint64, pfx string) *types.Module {
	modID := su.NextID()
	mod := &types.Module{
		ID:          modID,
		NamespaceID: nsID,
		Handle:      pfx + "_module",
		Fields: types.ModuleFieldSet{
			{
				ID:       su.NextID(),
				ModuleID: modID,
				Kind:     "Bool",
				Name:     "BoolTrue",
			},
			{
				ID:       su.NextID(),
				ModuleID: modID,
				Kind:     "Bool",
				Name:     "BoolFalse",
			},
			{
				ID:       su.NextID(),
				ModuleID: modID,
				Kind:     "DateTime",
				Name:     "DateTime",
			},
			{
				ID:       su.NextID(),
				ModuleID: modID,
				Kind:     "Email",
				Name:     "Email",
			},
			{
				ID:       su.NextID(),
				ModuleID: modID,
				Kind:     "Select",
				Name:     "Select",
				Options: types.ModuleFieldOptions{
					"options": []string{
						"v1",
						"v2",
					},
				},
			},
			{
				ID:       su.NextID(),
				ModuleID: modID,
				Kind:     "Number",
				Name:     "Number",
				Options: types.ModuleFieldOptions{
					"precision": 2,
				},
			},
			{
				ID:       su.NextID(),
				ModuleID: modID,
				Kind:     "String",
				Name:     "String",
			},
			{
				ID:       su.NextID(),
				ModuleID: modID,
				Kind:     "Url",
				Name:     "Url",
			},
			{
				ID:       su.NextID(),
				ModuleID: modID,
				Kind:     "User",
				Name:     "User",
			},
		},
	}
	err := store.CreateComposeModule(ctx, s, mod)
	if err != nil {
		t.Fatal(err)
	}
	err = store.CreateComposeModuleField(ctx, s, mod.Fields...)
	if err != nil {
		t.Fatal(err)
	}
	return mod
}

func sTestComposePage(ctx context.Context, t *testing.T, s store.Storer, nsID uint64, pfx string) *types.Page {
	return sTestComposePageWithBlocks(ctx, t, s, nsID, pfx, types.PageBlocks{
		types.PageBlock{
			Title:       "page block content",
			Description: "description",
			Kind:        "Content",
		},
		types.PageBlock{
			Title:       "page block qwerty",
			Description: "description",
			Kind:        "Qwerty",
		},
	})
}

func sTestComposePageWithBlocks(ctx context.Context, t *testing.T, s store.Storer, nsID uint64, pfx string, bb types.PageBlocks) *types.Page {
	ns := &types.Page{
		ID:          su.NextID(),
		NamespaceID: nsID,
		Handle:      pfx + "_page",
		Title:       pfx + " page",
		Description: "description",
		Blocks:      bb,
		Visible:     true,
		Weight:      0,
		CreatedAt:   createdAt,
		UpdatedAt:   &updatedAt,
	}

	err := store.CreateComposePage(ctx, s, ns)
	if err != nil {
		t.Fatal(err)
	}

	return ns
}

func sTestComposeChart(ctx context.Context, t *testing.T, s store.Storer, nsID, modID uint64, pfx string) *types.Chart {
	chr := &types.Chart{
		ID:          su.NextID(),
		NamespaceID: nsID,
		Handle:      pfx + "_chart",
		Name:        pfx + " chart",

		Config: types.ChartConfig{
			Reports: []*types.ChartConfigReport{
				{
					Filter:   "filter",
					ModuleID: modID,
					YAxis: map[string]interface{}{
						"beginAtZero": true,
						"label":       "Euro",
					},
				},
			},
			ColorScheme: "colorscheme",
		},

		CreatedAt: createdAt,
		UpdatedAt: &updatedAt,
	}

	err := store.CreateComposeChart(ctx, s, chr)
	if err != nil {
		t.Fatal(err)
	}

	return chr
}

func sTestComposeRecord(ctx context.Context, t *testing.T, s store.Storer, nsID, modID, usrID uint64) *types.Record {
	recID := su.NextID()
	rec := &types.Record{
		ID:          recID,
		NamespaceID: nsID,
		ModuleID:    modID,

		Values: types.RecordValueSet{
			{
				RecordID: recID,
				Name:     "module_field_string",
				Value:    "string value",
			},
			{
				RecordID: recID,
				Name:     "module_field_number",
				Value:    "10",
			},
		},

		CreatedAt: createdAt,
		UpdatedAt: &updatedAt,
		OwnedBy:   usrID,
		CreatedBy: usrID,
		UpdatedBy: usrID,
	}

	mod := &types.Module{
		ID:          modID,
		NamespaceID: nsID,
	}

	err := store.CreateComposeRecord(ctx, s, mod, rec)
	if err != nil {
		t.Fatal(err)
	}

	return rec
}

package dal

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMerge(t *testing.T) {
	tcc := []struct {
		name string
		aa   AlterationSet
		bb   AlterationSet
		cc   AlterationSet
	}{
		{
			name: "empty",
			aa:   AlterationSet{},
			bb:   AlterationSet{},
			cc:   AlterationSet{},
		},
		{
			name: "aa empty",
			aa:   AlterationSet{},
			bb: AlterationSet{
				&Alteration{
					AttributeAdd: &AttributeAdd{
						Attr: &Attribute{
							Ident: "foo",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
			},
			cc: AlterationSet{
				&Alteration{
					AttributeAdd: &AttributeAdd{
						Attr: &Attribute{
							Ident: "foo",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
			},
		},
		{
			name: "bb empty",
			aa: AlterationSet{
				&Alteration{
					AttributeAdd: &AttributeAdd{
						Attr: &Attribute{
							Ident: "foo",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
			},
			bb: AlterationSet{},
			cc: AlterationSet{
				&Alteration{
					AttributeAdd: &AttributeAdd{
						Attr: &Attribute{
							Ident: "foo",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
			},
		},
		{
			name: "remove duplicates",
			aa:   AlterationSet{},
			bb: AlterationSet{
				&Alteration{
					AttributeAdd: &AttributeAdd{
						Attr: &Attribute{
							Ident: "foo",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
				&Alteration{
					AttributeAdd: &AttributeAdd{
						Attr: &Attribute{
							Ident: "foo",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
				&Alteration{
					AttributeAdd: &AttributeAdd{
						Attr: &Attribute{
							Ident: "foo",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
			},
			cc: AlterationSet{
				&Alteration{
					AttributeAdd: &AttributeAdd{
						Attr: &Attribute{
							Ident: "foo",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
			},
		},
		{
			name: "un matching types",
			aa: AlterationSet{
				&Alteration{
					AttributeAdd: &AttributeAdd{
						Attr: &Attribute{
							Ident: "foo",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
			},
			bb: AlterationSet{
				&Alteration{
					AttributeDelete: &AttributeDelete{
						Attr: &Attribute{
							Ident: "foo",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
			},
			cc: AlterationSet{
				&Alteration{
					AttributeAdd: &AttributeAdd{
						Attr: &Attribute{
							Ident: "foo",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
				&Alteration{
					AttributeDelete: &AttributeDelete{
						Attr: &Attribute{
							Ident: "foo",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
			},
		},
		{
			name: "match AttributeAdd",
			aa: AlterationSet{
				&Alteration{
					AttributeAdd: &AttributeAdd{
						Attr: &Attribute{
							Ident: "foo",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
			},
			bb: AlterationSet{
				&Alteration{
					AttributeAdd: &AttributeAdd{
						Attr: &Attribute{
							Ident: "foo",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
			},
			cc: AlterationSet{
				&Alteration{
					AttributeAdd: &AttributeAdd{
						Attr: &Attribute{
							Ident: "foo",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
			},
		},
		{
			name: "not match AttributeAdd",
			aa: AlterationSet{
				&Alteration{
					AttributeAdd: &AttributeAdd{
						Attr: &Attribute{
							Ident: "foo",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
			},
			bb: AlterationSet{
				&Alteration{
					AttributeAdd: &AttributeAdd{
						Attr: &Attribute{
							Ident: "bar",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
			},
			cc: AlterationSet{
				&Alteration{
					AttributeAdd: &AttributeAdd{
						Attr: &Attribute{
							Ident: "foo",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
				&Alteration{
					AttributeAdd: &AttributeAdd{
						Attr: &Attribute{
							Ident: "bar",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
			},
		},

		{
			name: "match AttributeDelete",
			aa: AlterationSet{
				&Alteration{
					AttributeDelete: &AttributeDelete{
						Attr: &Attribute{
							Ident: "foo",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
			},
			bb: AlterationSet{
				&Alteration{
					AttributeDelete: &AttributeDelete{
						Attr: &Attribute{
							Ident: "foo",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
			},
			cc: AlterationSet{
				&Alteration{
					AttributeDelete: &AttributeDelete{
						Attr: &Attribute{
							Ident: "foo",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
			},
		},
		{
			name: "not match AttributeDelete",
			aa: AlterationSet{
				&Alteration{
					AttributeDelete: &AttributeDelete{
						Attr: &Attribute{
							Ident: "foo",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
			},
			bb: AlterationSet{
				&Alteration{
					AttributeDelete: &AttributeDelete{
						Attr: &Attribute{
							Ident: "bar",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
			},
			cc: AlterationSet{
				&Alteration{
					AttributeDelete: &AttributeDelete{
						Attr: &Attribute{
							Ident: "foo",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
				&Alteration{
					AttributeDelete: &AttributeDelete{
						Attr: &Attribute{
							Ident: "bar",
							Type:  &TypeJSON{Nullable: false},
						},
					},
				},
			},
		},

		{
			name: "match AttributeReType",
			aa: AlterationSet{
				&Alteration{
					AttributeReType: &AttributeReType{
						Attr: &Attribute{
							Ident: "foo",
						},
						To: &TypeBoolean{Nullable: false},
					},
				},
			},
			bb: AlterationSet{
				&Alteration{
					AttributeReType: &AttributeReType{
						Attr: &Attribute{
							Ident: "foo",
						},
						To: &TypeBoolean{Nullable: false},
					},
				},
			},
			cc: AlterationSet{
				&Alteration{
					AttributeReType: &AttributeReType{
						Attr: &Attribute{
							Ident: "foo",
						},
						To: &TypeBoolean{Nullable: false},
					},
				},
			},
		},
		{
			name: "not match AttributeReType",
			aa: AlterationSet{
				&Alteration{
					AttributeReType: &AttributeReType{
						Attr: &Attribute{
							Ident: "foo",
						},
						To: &TypeBoolean{Nullable: false},
					},
				},
			},
			bb: AlterationSet{
				&Alteration{
					AttributeReType: &AttributeReType{
						Attr: &Attribute{
							Ident: "foo",
						},
						To: &TypeText{Nullable: false},
					},
				},
			},
			cc: AlterationSet{
				&Alteration{
					AttributeReType: &AttributeReType{
						Attr: &Attribute{
							Ident: "foo",
						},
						To: &TypeBoolean{Nullable: false},
					},
				},
				&Alteration{
					AttributeReType: &AttributeReType{
						Attr: &Attribute{
							Ident: "foo",
						},
						To: &TypeText{Nullable: false},
					},
				},
			},
		},

		{
			name: "match AttributeReEncode",
			aa: AlterationSet{
				&Alteration{
					AttributeReEncode: &AttributeReEncode{
						Attr: &Attribute{Ident: "foo"},
						To:   &CodecPlain{},
					},
				},
			},
			bb: AlterationSet{
				&Alteration{
					AttributeReEncode: &AttributeReEncode{
						Attr: &Attribute{Ident: "foo"},
						To:   &CodecPlain{},
					},
				},
			},
			cc: AlterationSet{
				&Alteration{
					AttributeReEncode: &AttributeReEncode{
						Attr: &Attribute{Ident: "foo"},
						To:   &CodecPlain{},
					},
				},
			},
		},
		{
			name: "not match AttributeReEncode",
			aa: AlterationSet{
				&Alteration{
					AttributeReEncode: &AttributeReEncode{
						Attr: &Attribute{Ident: "foo"},
						To:   &CodecPlain{},
					},
				},
			},
			bb: AlterationSet{
				&Alteration{
					AttributeReEncode: &AttributeReEncode{
						Attr: &Attribute{Ident: "foo"},
						To:   &CodecAlias{Ident: "foo2"},
					},
				},
			},
			cc: AlterationSet{
				&Alteration{
					AttributeReEncode: &AttributeReEncode{
						Attr: &Attribute{Ident: "foo"},
						To:   &CodecPlain{},
					},
				},
				&Alteration{
					AttributeReEncode: &AttributeReEncode{
						Attr: &Attribute{Ident: "foo"},
						To:   &CodecAlias{Ident: "foo2"},
					},
				},
			},
		},

		{
			name: "match ModelAdd",
			aa: AlterationSet{
				&Alteration{
					ModelAdd: &ModelAdd{
						Model: &Model{Ident: "foo"},
					},
				},
			},
			bb: AlterationSet{
				&Alteration{
					ModelAdd: &ModelAdd{
						Model: &Model{Ident: "foo"},
					},
				},
			},
			cc: AlterationSet{
				&Alteration{
					ModelAdd: &ModelAdd{
						Model: &Model{Ident: "foo"},
					},
				},
			},
		},
		{
			name: "not match ModelAdd",
			aa: AlterationSet{
				&Alteration{
					ModelAdd: &ModelAdd{
						Model: &Model{Ident: "foo"},
					},
				},
			},
			bb: AlterationSet{
				&Alteration{
					ModelAdd: &ModelAdd{
						Model: &Model{Ident: "bar"},
					},
				},
			},
			cc: AlterationSet{
				&Alteration{
					ModelAdd: &ModelAdd{
						Model: &Model{Ident: "foo"},
					},
				},
				&Alteration{
					ModelAdd: &ModelAdd{
						Model: &Model{Ident: "bar"},
					},
				},
			},
		},

		{
			name: "match ModelDelete",
			aa: AlterationSet{
				&Alteration{
					ModelDelete: &ModelDelete{
						Model: &Model{Ident: "foo"},
					},
				},
			},
			bb: AlterationSet{
				&Alteration{
					ModelDelete: &ModelDelete{
						Model: &Model{Ident: "foo"},
					},
				},
			},
			cc: AlterationSet{
				&Alteration{
					ModelDelete: &ModelDelete{
						Model: &Model{Ident: "foo"},
					},
				},
			},
		},
		{
			name: "not match ModelDelete",
			aa: AlterationSet{
				&Alteration{
					ModelDelete: &ModelDelete{
						Model: &Model{Ident: "foo"},
					},
				},
			},
			bb: AlterationSet{
				&Alteration{
					ModelDelete: &ModelDelete{
						Model: &Model{Ident: "bar"},
					},
				},
			},
			cc: AlterationSet{
				&Alteration{
					ModelDelete: &ModelDelete{
						Model: &Model{Ident: "foo"},
					},
				},
				&Alteration{
					ModelDelete: &ModelDelete{
						Model: &Model{Ident: "bar"},
					},
				},
			},
		},
	}

	req := require.New(t)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			cc := tc.aa.Merge(tc.bb)
			req.Len(cc, len(tc.cc))
			for i, c := range cc {
				cmp := c.compare(*tc.cc[i])
				req.True(cmp)
			}
		})
	}
}

func TestReTypeMarshling(t *testing.T) {
	a := &AttributeReType{
		Attr: &Attribute{
			Ident: "foo",
			Type:  &TypeBlob{},
			Store: &CodecPlain{},
		},
		To: &TypeBoolean{},
	}

	bb, err := json.Marshal(a)
	require.NoError(t, err)

	b := &AttributeReType{}
	err = json.Unmarshal(bb, &b)
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(a, b))
}

func TestReEncodeMarshling(t *testing.T) {
	a := &AttributeReEncode{
		Attr: &Attribute{
			Ident: "foo",
			Type:  &TypeBlob{},
			Store: &CodecPlain{},
		},
		To: &CodecAlias{Ident: "bar"},
	}

	bb, err := json.Marshal(a)
	require.NoError(t, err)

	b := &AttributeReEncode{}
	err = json.Unmarshal(bb, &b)
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(a, b))
}

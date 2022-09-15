package postgres

import (
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/ddl"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_parseIndexDefinition(t *testing.T) {
	tests := []struct {
		name    string
		def     string
		index   *ddl.Index
		wantErr bool
	}{
		{
			name: "simple index",
			def:  "CREATE INDEX users_email ON public.users USING btree (email)",
			index: &ddl.Index{
				Ident:      "users_email",
				TableIdent: "users",
				Type:       "btree",
				Fields: []*ddl.IndexField{
					{
						Column: "email",
					},
				},
			},
		},
		{
			name: "with predicate",
			def:  "CREATE UNIQUE INDEX users_email ON public.users USING btree (email) WHERE (deleted_at IS NULL)",
			index: &ddl.Index{
				Ident:      "users_email",
				TableIdent: "users",
				Type:       "btree",
				Unique:     true,
				Fields: []*ddl.IndexField{
					{
						Column: "email",
					},
				},
				Predicate: "deleted_at IS NULL",
			},
		},
		{
			name: "expressions with predicate",
			def:  "CREATE UNIQUE INDEX unique_language_handle ON public.templates USING btree (language, lower((handle)::text)) WHERE ((length((handle)::text) > 0) AND (deleted_at IS NULL))",
			index: &ddl.Index{
				Ident:      "unique_language_handle",
				TableIdent: "templates",
				Type:       "btree",
				Unique:     true,
				Fields: []*ddl.IndexField{
					{
						Column: "language",
					},
					{
						Expression: "lower((handle)::text)",
					},
				},
				Predicate: "(length((handle)::text) > 0) AND (deleted_at IS NULL)",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index := &ddl.Index{}
			if err := parseIndexDefinition(tt.def, index); (err != nil) != tt.wantErr {
				t.Errorf("parseIndexDefinition() error = %v, wantErr %v", err, tt.wantErr)
			}

			require.Equal(t, tt.index, index)
		})
	}
}

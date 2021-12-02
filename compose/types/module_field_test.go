package types

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/locale"
	"github.com/stretchr/testify/require"
)

func TestModuleField_decodeTranslationsOptionsOptionTexts(t *testing.T) {
	rti := locale.ResourceTranslationIndex{
		"meta.options.val1.text": &locale.ResourceTranslation{Msg: "TEXT-1"},
	}
	tests := []struct {
		name string
		opts ModuleFieldOptions
		out  interface{}
	}{
		// TODO: Add test cases.
		{"all-nil", nil, nil},
		{"options-nil",
			ModuleFieldOptions{
				"options": nil,
			},
			nil,
		},
		{"options-foo",
			ModuleFieldOptions{
				"options": "foo",
			},
			nil,
		},
		{"options-foo-number",
			ModuleFieldOptions{
				"options": []interface{}{1},
			},
			nil,
		},
		{"options-valid-slice",
			ModuleFieldOptions{
				"options": []interface{}{
					"val1",
					"val2",
				},
			},
			[]interface{}{
				map[string]string{"value": "val1", "text": "TEXT-1"},
				map[string]string{"value": "val2", "text": "val2"},
			},
		},
		{"options-valid-map",
			ModuleFieldOptions{
				"options": []interface{}{
					map[string]interface{}{"value": "val1", "text": "Text1"},
					map[string]interface{}{"value": "val2", "text": "Text2"},
				},
			},
			[]interface{}{
				map[string]string{"value": "val1", "text": "TEXT-1"},
				map[string]string{"value": "val2", "text": "Text2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				req = require.New(t)
				f   = &ModuleField{Options: tt.opts}
			)

			f.decodeTranslationsOptionsOptionTexts(rti)

			if tt.out != nil {
				req.Equal(tt.out, f.Options["options"])
			}
		})
	}
}

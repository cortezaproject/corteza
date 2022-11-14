package renderer

import (
	"bytes"
	"context"
	"io"
	"regexp"

	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	genericText struct {
		def DriverDefinition
	}
	genericTextDriver struct{}
)

var (
	plainTextRegex = regexp.MustCompile("text(/)?.*")
)

func newGenericText() driverFactory {
	return &genericText{
		def: DriverDefinition{
			Name: "genericText",
			InputTypes: []types.DocumentType{
				types.DocumentTypePlain,
				types.DocumentTypeHTML,
			},
			OutputTypes: []types.DocumentType{
				types.DocumentTypePlain,
			},
		},
	}
}

func (d *genericText) Define() DriverDefinition {
	return d.def
}

func (d *genericText) CanRender(t types.DocumentType) bool {
	for _, i := range d.def.InputTypes {
		if i == t {
			return true
		}
	}
	return false
}

func (d *genericText) CanProduce(t types.DocumentType) bool {
	for _, o := range d.def.OutputTypes {
		if o == t {
			return true
		}
	}
	return false
}

func (d *genericText) Driver() driver {
	return &genericTextDriver{}
}

func (d *genericTextDriver) Render(ctx context.Context, pl *driverPayload) (io.ReadSeeker, error) {
	t, err := preprocPlainTemplate(pl.Template, pl.Partials)
	if err != nil {
		return nil, err
	}

	dd := &bytes.Buffer{}
	err = t.Execute(dd, pl.Variables)

	return bytes.NewReader(dd.Bytes()), err
}

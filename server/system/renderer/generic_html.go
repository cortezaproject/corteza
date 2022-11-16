package renderer

import (
	"bytes"
	"context"
	"io"

	"github.com/cortezaproject/corteza/server/system/types"
)

type (
	genericHTML struct {
		def DriverDefinition
	}
	genericHTMLDriver struct{}
)

func newGenericHTML() driverFactory {
	return &genericHTML{
		def: DriverDefinition{
			Name: "genericHTML",
			InputTypes: []types.DocumentType{
				types.DocumentTypePlain,
				types.DocumentTypeHTML,
			},
			OutputTypes: []types.DocumentType{
				types.DocumentTypeHTML,
			},
		},
	}
}

func (d *genericHTML) Define() DriverDefinition {
	return d.def
}

func (d *genericHTML) CanRender(t types.DocumentType) bool {
	for _, i := range d.def.InputTypes {
		if i == t {
			return true
		}
	}
	return false
}

func (d *genericHTML) CanProduce(t types.DocumentType) bool {
	for _, o := range d.def.OutputTypes {
		if o == t {
			return true
		}
	}
	return false
}

func (d *genericHTML) Driver() driver {
	return &genericHTMLDriver{}
}

func (d *genericHTMLDriver) Render(ctx context.Context, pl *driverPayload) (io.ReadSeeker, error) {
	t, err := preprocHTMLTemplate(pl)
	if err != nil {
		return nil, err
	}

	dd := &bytes.Buffer{}
	err = t.Execute(dd, pl.Variables)

	return bytes.NewReader(dd.Bytes()), err
}

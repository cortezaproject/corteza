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
	if pl.Header == nil && pl.Footer == nil {
		t, err := preprocHTMLTemplate(pl)
		if err != nil {
			return nil, err
		}

		dd := &bytes.Buffer{}
		err = t.Execute(dd, pl.Variables)

		return bytes.NewReader(dd.Bytes()), err
	}

	// Header and footer are rendered with the same variables and joined with
	// the content into a single document; partials are one-shot readers, so
	// they are buffered to let each of the three templates parse its own copy
	partials := make(map[string][]byte, len(pl.Partials))
	for h, r := range pl.Partials {
		bb, err := io.ReadAll(r)
		if err != nil {
			return nil, err
		}
		partials[h] = bb
	}

	dd := &bytes.Buffer{}
	for _, src := range []io.Reader{pl.Header, pl.Template, pl.Footer} {
		if src == nil {
			continue
		}

		aux := *pl
		aux.Template = src
		aux.Partials = make(map[string]io.Reader, len(partials))
		for h, bb := range partials {
			aux.Partials[h] = bytes.NewReader(bb)
		}

		t, err := preprocHTMLTemplate(&aux)
		if err != nil {
			return nil, err
		}

		if err = t.Execute(dd, aux.Variables); err != nil {
			return nil, err
		}
	}

	return bytes.NewReader(dd.Bytes()), nil
}

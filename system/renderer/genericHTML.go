package renderer

import (
	"bytes"
	"context"
	"io"

	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	genericHTML       struct{}
	genericHTMLDriver struct{}
)

func newGenericHTML() driverFactory {
	return &genericHTML{}
}

func (d *genericHTML) CanRender(t types.DocumentType) bool {
	return t == types.DocumentTypeHTML || t == types.DocumentTypePlain
}

func (d *genericHTML) CanProduce(t types.DocumentType) bool {
	return t == types.DocumentTypeHTML
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

package renderer

import (
	"bytes"
	"context"
	"io"
	"regexp"

	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	genericTXT       struct{}
	genericTXTDriver struct{}
)

var (
	plainTextRegex = regexp.MustCompile("text(/)?.*")
)

func newGenericTXT() driverFactory {
	return &genericTXT{}
}

func (d *genericTXT) CanRender(t types.DocumentType) bool {
	return t == types.DocumentTypePlain || t == types.DocumentTypeHTML
}

func (d *genericTXT) CanProduce(t types.DocumentType) bool {
	return t == types.DocumentTypePlain
}

func (d *genericTXT) Driver() driver {
	return &genericTXTDriver{}
}

func (d *genericTXTDriver) Render(ctx context.Context, pl *driverPayload) (io.ReadSeeker, error) {
	t, err := preprocPlainTemplate(pl.Template, pl.Partials)
	if err != nil {
		return nil, err
	}

	dd := &bytes.Buffer{}
	err = t.Execute(dd, pl.Variables)

	return bytes.NewReader(dd.Bytes()), err
}

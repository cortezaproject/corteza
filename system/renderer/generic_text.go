package renderer

import (
	"bytes"
	"context"
	"io"
	"regexp"

	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	genericText       struct{}
	genericTextDriver struct{}
)

var (
	plainTextRegex = regexp.MustCompile("text(/)?.*")
)

func newGenericText() driverFactory {
	return &genericText{}
}

func (d *genericText) CanRender(t types.DocumentType) bool {
	return t == types.DocumentTypePlain || t == types.DocumentTypeHTML
}

func (d *genericText) CanProduce(t types.DocumentType) bool {
	return t == types.DocumentTypePlain
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

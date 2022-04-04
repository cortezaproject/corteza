package locale

import (
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

// tested with
// go test -count 10 -race -run TestServiceReloadAndTranslate ./pkg/locale/...
func TestServiceReloadAndTranslate(t *testing.T) {
	var (
		req      = require.New(t)
		svc, err = Service(zap.NewNop(), options.LocaleOpt{
			Languages: "en",
		})

		tag = language.English
	)

	req.NoError(err)
	req.NotNil(svc)
	go svc.ResourceTranslations(tag, "resource")
	go svc.ReloadStatic()
}

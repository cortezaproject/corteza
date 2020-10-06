package federation

import (
	"context"

	"github.com/cortezaproject/corteza-server/federation/importer"
	impAux "github.com/cortezaproject/corteza-server/pkg/importer"
	"github.com/cortezaproject/corteza-server/store"
	"go.uber.org/zap"
)

func Provision(ctx context.Context, log *zap.Logger, s store.Storer) error {
	log.Info("provisioning federation")
	readers, err := impAux.ReadStatic(Asset)
	if err != nil {
		return err
	}

	return importer.Import(ctx, readers...)
}

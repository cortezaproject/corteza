package compose

import (
	"context"
	"github.com/cortezaproject/corteza-server/compose/importer"
	"github.com/cortezaproject/corteza-server/compose/service"
	"github.com/cortezaproject/corteza-server/compose/types"
	impAux "github.com/cortezaproject/corteza-server/pkg/importer"
	"go.uber.org/zap"
)

func Provision(ctx context.Context, log *zap.Logger) error {
	if provisioned, err := notProvisioned(ctx); err != nil {
		return err
	} else if !provisioned {
		log.Info("provisioning compose")
		readers, err := impAux.ReadStatic(Asset)
		if err != nil {
			return err
		}

		return importer.Import(ctx, nil, readers...)
	}

	return nil
}

// provision only where there are no namespaces
func notProvisioned(ctx context.Context) (bool, error) {
	f := types.NamespaceFilter{}
	f.Limit = 1
	set, _, err := service.DefaultNamespace.With(ctx).Find(f)
	return len(set) == 0, err
}

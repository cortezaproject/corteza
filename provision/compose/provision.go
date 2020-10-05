package compose

import (
	"context"

	"github.com/cortezaproject/corteza-server/compose/importer"
	"github.com/cortezaproject/corteza-server/compose/types"
	impAux "github.com/cortezaproject/corteza-server/pkg/importer"
	"github.com/cortezaproject/corteza-server/store"
	"go.uber.org/zap"
)

// provision only where there are no namespaces
func hasNamespaces(ctx context.Context, s store.Storer) (bool, error) {
	if set, _, err := store.SearchComposeNamespaces(ctx, s, types.NamespaceFilter{}); err != nil {
		return false, err
	} else {
		return len(set) > 0, nil
	}
}

func Provision(ctx context.Context, log *zap.Logger, s store.Storer) error {
	if namespacesExist, err := hasNamespaces(ctx, s); err != nil {
		return err
	} else if !namespacesExist {
		log.Info("provisioning compose")
		readers, err := impAux.ReadStatic(Asset)
		if err != nil {
			return err
		}

		return importer.Import(ctx, nil, readers...)
	} else {
		log.Info("compose already provisioned")
	}

	return nil
}

package compose

import (
	"context"

	"github.com/cortezaproject/corteza-server/provision/util"
	"github.com/cortezaproject/corteza-server/store"
	"go.uber.org/zap"
)

func hasNamespaces(ctx context.Context, s store.Storer, slug string) (bool, error) {
	ns, err := store.LookupComposeNamespaceBySlug(ctx, s, slug)
	if err == store.ErrNotFound {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return ns != nil, nil
}

func Provision(ctx context.Context, log *zap.Logger, s store.Storer) error {
	log.Info("provisioning compose: system")
	if err := util.EncodeStatik(ctx, s, Asset, "/system"); err != nil {
		return err
	}
	log.Info("provisioned compose: system")

	if has, err := hasNamespaces(ctx, s, "crm"); !has && err == nil {
		log.Info("provisioning compose: crm")
		if err := util.EncodeStatik(ctx, s, Asset, "/crm"); err != nil {
			return err
		}
		log.Info("provisioned compose: crm")
	} else if err != nil {
		return err
	} else {
		log.Info("[skip] provisioning compose: crm")
	}

	if has, err := hasNamespaces(ctx, s, "service-solution"); !has && err == nil {
		log.Info("provisioning compose: service-solution")
		if err := util.EncodeStatik(ctx, s, Asset, "/service-solution"); err != nil {
			return err
		}
		log.Info("provisioned compose: service-solution")
	} else if err != nil {
		return err
	} else {
		log.Info("[skip] provisioning compose: service-solution")
	}

	return nil
}

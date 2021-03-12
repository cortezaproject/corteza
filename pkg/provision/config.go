package provision

import (
	"context"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"path/filepath"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/envoy/directory"
	"github.com/cortezaproject/corteza-server/pkg/envoy/resource"
	es "github.com/cortezaproject/corteza-server/pkg/envoy/store"
	"github.com/cortezaproject/corteza-server/pkg/envoy/yaml"
	"github.com/cortezaproject/corteza-server/store"
	"go.uber.org/zap"
)

// imports configuration files from path(s)
//
// paths can be colon delimited list of absolute or relative paths and/or with glob pattern
func importConfig(ctx context.Context, log *zap.Logger, s store.Storer, paths string) error {
	if can, err := canImportConfig(ctx, s); !can {
		log.Info("already provisioned, skipping full config import")
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to check if config import can be done: %w", err)
	}

	var (
		yd  = yaml.Decoder()
		nn  = make([]resource.Interface, 0, 200)
		se  = es.NewStoreEncoder(s, &es.EncoderConfig{OnExisting: resource.MergeLeft})
		bld = envoy.NewBuilder(se)

		sources = make([]string, 0, 16)
	)

	log.Info("importing config", zap.String("paths", paths))

	// verify all paths before doing the actual import
	for _, path := range strings.Split(paths, ":") {
		if aux, err := filepath.Glob(path); err != nil {
			return err
		} else {
			sources = append(sources, aux...)
		}
	}

	for _, path := range sources {
		log.Info("provisioning from path", zap.String("path", path))
		if mm, err := directory.Decode(ctx, path, yd); err != nil {
			return err
		} else {
			nn = append(nn, mm...)
		}
	}

	if g, err := bld.Build(ctx, nn...); err != nil {
		return err
	} else if err = envoy.Encode(ctx, g, se); err != nil {
		return err
	}

	return nil
}

// canImportConfig checks state of the store and
// verifies if Corteza should be provisioned (ie config should be imported)
func canImportConfig(ctx context.Context, s store.Storer) (bool, error) {
	rr, _, err := store.SearchRbacRules(ctx, s, rbac.RuleFilter{})
	return len(rr) == 0, err
}

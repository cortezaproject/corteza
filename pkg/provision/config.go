package provision

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/filter"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/system/types"

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
	canImportFull, err := canImportConfig(ctx, s)
	if !canImportFull {
		log.Info("already provisioned, skipping full config import")
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

	// verify all paths before doing the actual import
	for _, path := range strings.Split(paths, ":") {
		if aux, err := filepath.Glob(path); err != nil {
			return err
		} else {
			sources = append(sources, aux...)
		}
	}

	if canImportFull {
		log.Info("importing all configs", zap.String("paths", paths))
		for _, path := range sources {
			log.Info("provisioning from path", zap.String("path", path))
			if mm, err := directory.Decode(ctx, path, yd); err != nil {
				return err
			} else {
				nn = append(nn, mm...)
			}
		}
	} else {
		nn, err = collectUnimportedConfigs(ctx, log, s, sources, yd)
		if err != nil {
			return err
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

func collectUnimportedConfigs(ctx context.Context, log *zap.Logger, s store.Storer, sources []string, dec directory.Decoder) (nn []resource.Interface, err error) {
	var (
		searchPartialDirectories = []string{
			"002_templates",
			"003_auth",
		}
	)

	return nn, store.Tx(ctx, s, func(ctx context.Context, s store.Storer) (err error) {
		log.Debug("verifying partial config import for templates")
		set, _, err := store.SearchTemplates(ctx, s, types.TemplateFilter{Deleted: filter.StateInclusive})
		// Import only of no templates exist
		if err != nil || len(set) == 0 {
			return err
		}

		for _, d := range searchPartialDirectories {
			if list, e := decodeDirectory(ctx, sources, d, dec); e != nil {
				return fmt.Errorf("failed to decode template configs: %w", err)
			} else if len(list) == 0 {
				log.Error("failed to execute partial config import for templates, directory not found or no configs", zap.Any("dir", d))
				return
			} else {
				nn = append(nn, list...)
			}
		}

		return
	})
}

func hasSourceDir(sources []string, dir string) (string, bool) {
	for _, source := range sources {
		if strings.HasSuffix(source, dir) {
			return source, true
		}
	}
	return "", false
}

func decodeDirectory(ctx context.Context, sources []string, dir string, dec directory.Decoder) (res []resource.Interface, err error) {
	if source, has := hasSourceDir(sources, dir); has {
		res, err = directory.Decode(ctx, source, dec)
	}

	return
}

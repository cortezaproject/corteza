package provision

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cortezaproject/corteza/server/pkg/dal"
	"github.com/cortezaproject/corteza/server/pkg/envoyx"
	se "github.com/cortezaproject/corteza/server/system/envoy"
	"github.com/cortezaproject/corteza/server/system/types"

	"github.com/cortezaproject/corteza/server/pkg/rbac"

	"github.com/cortezaproject/corteza/server/store"
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
		nn   envoyx.NodeSet
		auxN envoyx.NodeSet
		evy  = envoyx.Global()

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
			auxN, _, err = evy.Decode(ctx, envoyx.DecodeParams{
				Type: envoyx.DecodeTypeURI,
				Params: map[string]any{
					"uri": "file://" + path,
				},
			})
			if err != nil {
				return err
			}
			nn = append(nn, auxN...)
		}
	} else {
		nn, err = collectUnimportedConfigs(ctx, log, s, sources, evy)
		if err != nil {
			return err
		}
	}

	// Get potentially missing refs
	//
	// @todo replace this with getting missing refs and just fetching those.
	//       For now this is the only scenario so we can get away with this just fine.
	rr, _, err := store.SearchRoles(ctx, s, types.RoleFilter{})
	if err != nil {
		return err
	}
	for _, r := range rr {
		aux, err := se.RoleToEnvoyNode(r)
		if err != nil {
			return err
		}
		aux.Placeholder = true
		nn = append(nn, aux)
	}
	// ----------------------------------------------------------------------

	return store.Tx(ctx, s, func(ctx context.Context, s store.Storer) (err error) {
		ep := envoyx.EncodeParams{
			Type: envoyx.EncodeTypeStore,
			Params: map[string]any{
				"storer": s,
				"dal":    dal.Service(),
			},
			Envoy: envoyx.EnvoyConfig{
				MergeAlg: envoyx.OnConflictSkip,
			},
		}

		gg, err := evy.Bake(ctx, ep,
			nil,
			nn...,
		)
		if err != nil {
			return err
		}

		err = evy.Encode(ctx, ep, gg)
		if err != nil {
			return err
		}

		return nil
	})
}

// canImportConfig checks state of the store and
// verifies if Corteza should be provisioned (ie config should be imported)
func canImportConfig(ctx context.Context, s store.Storer) (bool, error) {
	rr, _, err := store.SearchRbacRules(ctx, s, rbac.RuleFilter{})
	return len(rr) == 0, err
}

func collectUnimportedConfigs(ctx context.Context, log *zap.Logger, s store.Storer, sources []string, evy *envoyx.Service) (nn envoyx.NodeSet, err error) {
	var (
		searchPartialDirectories = []uConfig{
			{dir: "000_base", fn: provisionPartialBase},
			{dir: "002_templates", fn: provisionPartialTemplates},
			{dir: "003_auth", fn: provisionPartialAuthClients},
			{dir: "200_federation", fn: nil},
			{dir: "300_automation", fn: nil},
		}
	)

	err = store.Tx(ctx, s, func(ctx context.Context, s store.Storer) (err error) {
		nn, err = collectPartialDirectories(log, searchPartialDirectories,
			func(d uConfig) bool { return d.fn == nil || d.fn(ctx, s, log) },
			func(dir string) (envoyx.NodeSet, error) { return decodeDirectory(ctx, sources, dir, evy) },
		)
		return
	})

	return
}

// collectPartialDirectories decodes each directory that needs importing; one
// that is missing or empty is reported and skipped, not the end of the list
func collectPartialDirectories(log *zap.Logger, dirs []uConfig, needed func(uConfig) bool, decode func(dir string) (envoyx.NodeSet, error)) (nn envoyx.NodeSet, err error) {
	for _, d := range dirs {
		if !needed(d) {
			log.Debug("skipping partial config import, no changes", zap.String("dir", d.dir))
			continue
		}

		list, err := decode(d.dir)
		if err != nil {
			return nil, fmt.Errorf("failed to decode %s configs: %w", d.dir, err)
		}

		if len(list) == 0 {
			log.Error("skipping partial config import, directory not found or no configs", zap.String("dir", d.dir))
			continue
		}

		log.Debug("partial import ready", zap.String("dir", d.dir))
		nn = append(nn, list...)
	}

	return
}

func hasSourceDir(sources []string, dir string) (string, bool) {
	for _, source := range sources {
		if strings.HasSuffix(source, dir) {
			return source, true
		}
	}
	return "", false
}

func decodeDirectory(ctx context.Context, sources []string, dir string, evy *envoyx.Service) (res envoyx.NodeSet, err error) {
	if source, has := hasSourceDir(sources, dir); has {
		res, _, err = evy.Decode(ctx, envoyx.DecodeParams{
			Type: envoyx.DecodeTypeURI,
			Params: map[string]any{
				"uri": "file://" + source,
			},
		})
	}

	return
}

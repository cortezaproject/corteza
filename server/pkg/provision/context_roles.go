package provision

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/crusttech/human/server/store"
	"github.com/crusttech/human/server/system/types"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// contextRoleTypes gives a stored context role the resource-type list its
// expression is meant to be tried on, where it has none.
//
// A context role whose list is empty is skipped for every resource
// (pkg/rbac/roles.go), so it holds its rules and never activates. Roles import
// with envoyx.OnConflictSkip, which keeps the stored row verbatim, so a
// re-import restores the rules and never the role — nothing else reconciles the
// two.
//
// Only an absent list is filled: a deployment that has narrowed a context role
// keeps its own.
func contextRoleTypes(ctx context.Context, log *zap.Logger, s store.Storer, paths string) error {
	declared, err := declaredContextRoles(paths)
	if err != nil {
		return err
	}

	if len(declared) == 0 {
		log.Debug("no base roles config found, skipping")
		return nil
	}

	stored, err := loadRoles(ctx, s)
	if err != nil {
		return err
	}

	repaired := fillContextRoleTypes(declared, stored)
	if len(repaired) == 0 {
		return nil
	}

	for _, r := range repaired {
		r.UpdatedAt = now()

		log.Info("context role had no resource types and could never activate; filling them in",
			zap.String("handle", r.Handle),
			zap.Strings("resourceTypes", r.Meta.Context.Resource))
	}

	if err = store.UpdateRole(ctx, s, repaired...); err != nil {
		return fmt.Errorf("failed to fill in context role resource types: %w", err)
	}

	return nil
}

// declaredContextRoles reads the base roles config and returns the resource-type
// list of every context role it declares, by handle.
func declaredContextRoles(paths string) (map[string][]string, error) {
	var sources []string

	for _, path := range strings.Split(paths, ":") {
		aux, err := filepath.Glob(path)
		if err != nil {
			return nil, err
		}

		sources = append(sources, aux...)
	}

	dir, has := hasSourceDir(sources, "000_base")
	if !has {
		return nil, nil
	}

	f, err := os.ReadFile(filepath.Join(dir, "roles.yaml"))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, err
	}

	var doc struct {
		Roles map[string]*types.Role `yaml:"roles"`
	}

	if err = yaml.Unmarshal(f, &doc); err != nil {
		return nil, fmt.Errorf("failed to read base roles config: %w", err)
	}

	out := make(map[string][]string, len(doc.Roles))

	for handle, r := range doc.Roles {
		if r == nil || r.Meta == nil || r.Meta.Context == nil {
			continue
		}

		if len(r.Meta.Context.Resource) == 0 {
			continue
		}

		out[handle] = r.Meta.Context.Resource
	}

	return out, nil
}

// fillContextRoleTypes sets the declared resource types on every stored context
// role that has an expression and no types, and returns the ones it changed.
//
// A role the config does not declare, one that is not contextual in the store,
// and one that already carries a list are all left exactly as they are.
func fillContextRoleTypes(declared map[string][]string, stored map[string]*types.Role) types.RoleSet {
	var out types.RoleSet

	for handle, rr := range declared {
		r := stored[handle]
		if r == nil || r.Meta == nil || r.Meta.Context == nil {
			continue
		}

		if r.Meta.Context.Expr == "" || len(r.Meta.Context.Resource) > 0 {
			continue
		}

		r.Meta.Context.Resource = rr
		out = append(out, r)
	}

	return out
}

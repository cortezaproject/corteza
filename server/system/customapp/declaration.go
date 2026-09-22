// Package customapp resolves what a custom application's page declares.
//
// It sits beside the callers rather than in system/service: resolving reaches
// compose, and system/service cannot — the two services would import each
// other. The service checks the result is resolved, so a caller that skips
// this is refused rather than storing a declaration no viewer can follow.
package customapp

import (
	"context"
	"strconv"

	cmpService "github.com/crusttech/human/server/compose/service"
	sysService "github.com/crusttech/human/server/system/service"
	"github.com/crusttech/human/server/system/types"
)

// Resolve fills in the IDs for the namespace slug and module handles a
// declaration names, as the caller: what it names is checked while whoever is
// storing the page can still correct it, and every viewer afterwards reads by
// ID, which asks nothing of them beyond the data itself.
func Resolve(ctx context.Context, meta *types.ApplicationSourceMeta) error {
	if meta == nil {
		return nil
	}

	meta.NamespaceID = 0
	meta.ModuleIDs = nil

	if meta.Namespace == "" {
		return nil
	}

	ns, err := cmpService.DefaultNamespace.FindByHandle(ctx, meta.Namespace)
	if err != nil || ns == nil {
		return sysService.ApplicationErrUnknownNamespace()
	}
	meta.NamespaceID = ns.ID

	meta.ModuleIDs = make(map[string]string, len(meta.Modules))
	for _, handle := range meta.Modules {
		mod, err := cmpService.DefaultModule.FindByHandle(ctx, ns.ID, handle)
		if err != nil || mod == nil {
			return sysService.ApplicationErrUnknownModule()
		}
		meta.ModuleIDs[handle] = strconv.FormatUint(mod.ID, 10)
	}

	return nil
}

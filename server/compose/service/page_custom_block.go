package service

import (
	"context"
	"fmt"
	"slices"
	"strconv"

	"github.com/crusttech/human/server/compose/types"
	"github.com/crusttech/human/server/store"
	systemService "github.com/crusttech/human/server/system/service"
	"github.com/spf13/cast"
)

// checkCustomBlocks holds a Custom block's own HTML to what a custom
// application's page is held to when it is stored: the sandbox must be able to
// run it, what it may change must be among what it reads, and every module it
// names must exist in the page's namespace. What it declares is stored the way
// the sandbox reads it — origins reduced to bare origins, and each module
// handle beside its ID, so a viewer reads by ID and needs no search. A block
// naming an application carries no HTML and is checked when that
// application's page is stored.
func checkCustomBlocks(ctx context.Context, s store.ComposeModules, namespaceID uint64, blocks types.PageBlocks) error {
	for i, b := range blocks {
		if b.Kind != "Custom" {
			continue
		}

		fail := func(err error) error { return fmt.Errorf("block %d (Custom): %w", i+1, err) }

		source := cast.ToString(b.Options["source"])
		if cast.ToUint64(b.Options["applicationID"]) > 0 && source == "" {
			continue
		}

		origins, err := systemService.NormalizeSourceOrigins(cast.ToStringSlice(b.Options["origins"]))
		if err != nil {
			return fail(err)
		}

		if err = systemService.CheckApplicationSource(source, origins); err != nil {
			return fail(err)
		}

		modules := cast.ToStringSlice(b.Options["modules"])
		for _, w := range cast.ToStringSlice(b.Options["writes"]) {
			if !slices.Contains(modules, w) {
				return fail(fmt.Errorf("it may change module %q but does not read it; add it to modules too", w))
			}
		}

		moduleIDs := make(map[string]string, len(modules))
		for _, handle := range modules {
			mod, err := store.LookupComposeModuleByNamespaceIDHandle(ctx, s, namespaceID, handle)
			if err != nil || mod == nil || mod.DeletedAt != nil {
				return fail(fmt.Errorf("there is no module %q in this page's namespace", handle))
			}
			moduleIDs[handle] = strconv.FormatUint(mod.ID, 10)
		}

		if b.Options != nil {
			b.Options["origins"] = origins
			b.Options["moduleIDs"] = moduleIDs
		}
	}

	return nil
}

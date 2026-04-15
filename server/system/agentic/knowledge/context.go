package knowledge

import (
	"context"
	"fmt"
	"strings"

	cmpTypes "github.com/cortezaproject/corteza/server/compose/types"
	sysTypes "github.com/cortezaproject/corteza/server/system/types"
)

type (
	KnowledgeBaseStore interface {
		FindByID(ctx context.Context, ID uint64) (*sysTypes.KnowledgeBase, error)
	}

	NamespaceLookup interface {
		FindByID(ctx context.Context, namespaceID uint64) (*cmpTypes.Namespace, error)
	}

	ModuleLookup interface {
		FindByID(ctx context.Context, namespaceID, moduleID uint64) (*cmpTypes.Module, error)
	}
)

func BuildContext(ctx context.Context, kbStore KnowledgeBaseStore, ns NamespaceLookup, mod ModuleLookup, ids []uint64) string {
	if len(ids) == 0 || ns == nil || mod == nil {
		return ""
	}

	var parts []string
	for _, id := range ids {
		kb, err := kbStore.FindByID(ctx, id)
		if err != nil || kb == nil {
			continue
		}

		text, err := render(ctx, ns, mod, kb)
		if err != nil || text == "" {
			continue
		}

		parts = append(parts, text)
	}

	return strings.Join(parts, "\n\n")
}

func sanitize(s string) string {
	s = strings.ReplaceAll(s, "\x00", "")
	for strings.Contains(s, "\n\n\n") {
		s = strings.ReplaceAll(s, "\n\n\n", "\n\n")
	}
	for _, marker := range []string{"<|", "|>", "###", "##", "---", "SYSTEM:", "ASSISTANT:", "USER:"} {
		s = strings.ReplaceAll(s, marker, "")
	}
	return strings.TrimSpace(s)
}

func render(ctx context.Context, ns NamespaceLookup, mod ModuleLookup, kb *sysTypes.KnowledgeBase) (string, error) {
	var parts []string

	if kb.Description != "" {
		parts = append(parts, sanitize(kb.Description))
	}

	if kb.Context != nil {
		composeCtx, err := renderComposeContext(ctx, ns, mod, kb.Title, kb.Context)
		if err != nil {
			return "", err
		}
		if composeCtx != "" {
			parts = append(parts, composeCtx)
		}
	}

	return strings.Join(parts, "\n\n"), nil
}

func renderComposeContext(ctx context.Context, ns NamespaceLookup, mod ModuleLookup, title string, c *sysTypes.KnowledgeBaseContext) (string, error) {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# %s\n", sanitize(title)))

	for _, nsCtx := range c.Namespaces {
		namespace, err := ns.FindByID(ctx, nsCtx.NamespaceID)
		if err != nil || namespace == nil {
			continue
		}

		sb.WriteString(fmt.Sprintf("\nNamespace: %s\n", sanitize(namespace.Name)))

		for _, moduleID := range nsCtx.ModuleIDs {
			module, err := mod.FindByID(ctx, namespace.ID, moduleID)
			if err != nil || module == nil {
				continue
			}

			sb.WriteString(fmt.Sprintf("\nModule: %s\n", sanitize(module.Name)))
			for _, f := range module.Fields {
				sb.WriteString(fmt.Sprintf("- %s (%s)\n", sanitize(f.Name), f.Kind))
			}
		}
	}

	return sb.String(), nil
}

package agentic

import (
	"strings"
	"testing"

	cmpTypes "github.com/crusttech/human/server/compose/types"
)

// A slug is optional on a namespace and plenty of real ones have none. The
// webapp resolves the :slug segment by slug OR by namespaceID
// (lib/vue/src/stores/useNamespaceStore.js getByUrlPart), so the ID is a valid
// url part — and keying on the slug alone returned no link at all for those
// namespaces, which is how this went unnoticed.
func TestNamespaceLinksFallBackToTheIDWhenThereIsNoSlug(t *testing.T) {
	links := namespaceLinks(&cmpTypes.Namespace{ID: 509895166871863297, Slug: ""})

	for key, want := range map[string]string{
		"url":     "/compose/namespace/509895166871863297",
		"editUrl": "/compose/namespaces/edit/509895166871863297",
	} {
		if got := links[key]; !strings.HasSuffix(got, want) {
			t.Errorf("%s is %q, want it to end in %q", key, got, want)
		}
	}
}

func TestNamespaceLinksPreferTheSlug(t *testing.T) {
	links := namespaceLinks(&cmpTypes.Namespace{ID: 509895166871863297, Slug: "crm"})

	if got := links["url"]; !strings.HasSuffix(got, "/compose/namespace/crm") {
		t.Errorf("url is %q, want it to end in /compose/namespace/crm", got)
	}
}

// A link is a convenience: a nil resource produces no links rather than a panic
// or a link to a list, which would read as "here is the thing" and show
// something else.
func TestLinksOfNothingAreNothing(t *testing.T) {
	if links := namespaceLinks(nil); links != nil {
		t.Errorf("got %v, want no links", links)
	}
}

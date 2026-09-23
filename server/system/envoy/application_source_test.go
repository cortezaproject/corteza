package envoy

import (
	"context"
	"testing"

	"github.com/crusttech/human/server/pkg/envoyx"
	"github.com/crusttech/human/server/system/types"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// A custom application's page travels with it, or "build it here, ship it
// there" does not work. The encoder writes the document and the decoder reads
// it back by field name, so the two have to agree on every key — a camel-cased
// one exports and never imports again, which loses the page's declaration
// while the page itself arrives, leaving an app that renders and reads nothing.
// The encoder asks a traverser about a node's neighbours; one application on
// its own has none.
type noTraverser struct{}

func (noTraverser) ParentForRef(*envoyx.Node, envoyx.Ref) *envoyx.Node          { return nil }
func (noTraverser) ParentForRT(*envoyx.Node, string) envoyx.NodeSet             { return nil }
func (noTraverser) ChildrenForResourceType(*envoyx.Node, string) envoyx.NodeSet { return nil }
func (noTraverser) Children(*envoyx.Node) envoyx.NodeSet                        { return nil }
func (noTraverser) NodeForRef(envoyx.Ref) *envoyx.Node                          { return nil }

func TestApplicationSourceSurvivesTheDocument(t *testing.T) {
	src := &types.Application{
		ID:      42,
		Name:    "Contacts desk",
		Enabled: true,
		Unify:   &types.ApplicationUnify{Name: "Contacts desk", Listed: true, Kind: "custom", Url: "app/42"},
		Source:  "<title>Desk</title><p>hello</p>",
		SourceMeta: &types.ApplicationSourceMeta{
			Namespace: "crm",
			Modules:   []string{"Lead", "Account"},
			Writes:    []string{"Lead"},
		},
	}

	node, err := YamlEncoder{}.encodeApplication(
		context.Background(),
		envoyx.EncodeParams{},
		&envoyx.Node{Resource: src, ResourceType: types.ApplicationResourceType},
		noTraverser{},
	)
	require.NoError(t, err)

	doc, err := yaml.Marshal(node)
	require.NoError(t, err)

	var back types.Application
	require.NoError(t, yaml.Unmarshal(doc, &back))

	require.Equal(t, src.Source, back.Source, "the page itself")
	require.NotNil(t, back.SourceMeta, "the declaration, without which the app reads nothing")
	require.Equal(t, src.SourceMeta.Namespace, back.SourceMeta.Namespace)
	require.Equal(t, src.SourceMeta.Modules, back.SourceMeta.Modules)
	require.Equal(t, src.SourceMeta.Writes, back.SourceMeta.Writes)
	require.Equal(t, "custom", back.Unify.Kind)
	require.True(t, back.Enabled)
}

package plugin

// import (
// 	"context"
// 	"testing"

// 	pb "github.com/cortezaproject/corteza/server/automation/plugin/bundler"
// 	"github.com/cortezaproject/corteza/server/automation/types"
// 	"github.com/cortezaproject/corteza/server/pkg/plugin"
// 	"github.com/davecgh/go-spew/spew"
// 	hcp "github.com/hashicorp/go-plugin"
// 	"go.uber.org/zap"
// )

// type (
// 	wfs struct{}

// 	automationRegistry struct{}
// 	pluginService      struct{}
// )

// func TestBundler_discoverBundles(t *testing.T) {
// 	// loop through designated dir and find bundles
// 	// valid bundle should have a meta.y?ml file
// 	// it can be empty otherwise

// 	// make sure to undispense the dispensed clients

// 	var ctx = context.Background()

// 	b := New(zap.NewNop())

// 	plugin.Setup(zap.NewNop().Named("test"), &automationRegistry{})

// 	b.registerBundler(
// 		// plugin.Service(),
// 		pb.New(&wfs{}),
// 		pb.NewAf(pluginService{}),
// 	)

// 	// we need another implementation to be able to read all files and extract them
// 	// loop .zip files (configurable) and .crust files

// 	// first we need to find all the folders
// 	// loop folders and check with each Bundler
// 	// after those are validated and registered, handle the zip files and do the same

// 	// b.handlers = append(b.handlers, packageHandlerZip{})
// 	b.handlers = append(b.handlers, packageHandlerFs{})

// 	b.Discover(ctx, "../../bundle")
// 	b.Validate(ctx)
// 	b.Register(ctx)

// 	t.Fail()
// }

// func (ar *automationRegistry) AddFunctions(ff ...*types.Function) {
// 	// spew.Dump("add function")
// }

// func (w *wfs) Create(context.Context, *types.Workflow) (*types.Workflow, error) {
// 	spew.Dump("persisting a new workflow")
// 	return nil, nil
// }

// func (ps pluginService) RegisterPlugin(ctx context.Context, path string, p hcp.Plugin) error {
// 	spew.Dump("registering plugin", path)
// 	return nil
// }

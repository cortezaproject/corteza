package plugin

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cortezaproject/corteza/server/pkg/plugin/automation"
	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/go-hclog"
	hcp "github.com/hashicorp/go-plugin"
	"go.uber.org/zap"
)

type (
	plugin struct {
		log *zap.Logger
		ar  automation.AutomationRegistry

		rawPlugins map[string]any
		// rawPlugins map[string]*cp
		clients   map[string]*hcp.Client
		pluginMap map[string]hcp.Plugin
	}
)

var (
	pl *plugin
)

const PLUGIN_PATH = "./plugin"

func Service() *plugin {
	return pl
}

func Setup(l *zap.Logger, ar automation.AutomationRegistry) {
	if pl != nil {
		return
	}

	pl = New(l, ar)
}

func New(l *zap.Logger, ar automation.AutomationRegistry) *plugin {
	return &plugin{
		log: l.Named("plugin"),
		ar:  ar,

		rawPlugins: make(map[string]any),
		// rawPlugins: make(map[string]*cp),
		clients:   make(map[string]*hcp.Client),
		pluginMap: make(map[string]hcp.Plugin),
	}
}

func (p *plugin) register(path string) {
	var (
		logger = hclog.New(&hclog.LoggerOptions{
			Name:   "plugin",
			Output: os.Stdout,
			Level:  hclog.Debug,
		})
		config = &hcp.ClientConfig{
			HandshakeConfig:  handshakeConfig,
			Plugins:          p.pluginMap,
			Cmd:              exec.Command("sh", "-c", path),
			AllowedProtocols: allowedProtocols,
			Logger:           logger,
		}
	)

	slug := makeSlug(path)

	p.clients[slug] = hcp.NewClient(config)

	p.log.Info("registered new GRPC client", zap.String("slug", slug))
}

func (p *plugin) unregister() {
	for _, c := range p.clients {
		c.Kill()
	}
}

// finds the registered client and prepares it in the internal registry
func (p *plugin) dispense2(slug string, ff func() any) (plugin any, err error) {
	if _, is := p.clients[slug]; !is {
		err = fmt.Errorf("could not find grpc client %s", slug)
		return
	}

	c, err := p.clients[slug].Client()

	if err != nil {
		p.log.Error("client register error", zap.Error(err))
		return
	}

	plugin, err = c.Dispense(slug)
	// raw, err := c.Dispense(slug)

	if err != nil {
		p.log.Error("client dispense error", zap.Error(err))
	}

	// cc := &cp{}

	// switch s := raw.(type) {
	// case automation.AutomationFunction:
	// 	cc.af = s
	// }

	// p.rawPlugins[slug] = ff()
	// err = nil

	return
}

// func (p *plugin) dispense(slug string) (err error) {
// 	if _, is := p.clients[slug]; !is {
// 		err = fmt.Errorf("could not find grpc client %s", slug)
// 		return
// 	}

// 	c, err := p.clients[slug].Client()

// 	if err != nil {
// 		p.log.Error("client register error", zap.Error(err))
// 		return
// 	}

// 	raw, err := c.Dispense(slug)

// 	if err != nil {
// 		p.log.Error("client dispense error", zap.Error(err))
// 		return
// 	}

// 	cc := &cp{}

// 	switch s := raw.(type) {
// 	case automation.AutomationFunction:
// 		cc.af = s
// 	}

// 	p.rawPlugins[slug] = cc

// 	return
// }

func (p *plugin) findPlugins(path string) (ss []string, err error) {
	ss, err = filepath.Glob(filepath.Join(path, "*"))
	return
}

func (p *plugin) registerPlugins(ctx context.Context, path string) {

	filePaths, err := p.findPlugins(path)

	spew.Dump("file paths", filePaths)

	if err != nil {
		p.log.Error("invalid file glob pattern", zap.Error(err))
		return
	}

	// for _, filePath := range filePaths {
	// 	// todo handle error
	// 	_ = p.RegisterPlugin(ctx, filePath, &automation.AutomationFunctionPlugin{})
	// }
}

func (p *plugin) SetRawPlugin(ctx context.Context, path string, rawPlugin any) (err error) {
	var (
		slug = makeSlug(path)
	)

	p.rawPlugins[slug] = rawPlugin
	err = nil

	return
}

func (p *plugin) RegisterPlugin(ctx context.Context, path string, pl hcp.Plugin, ff func() any) (rawPlugin any, err error) {
	var (
		slug = makeSlug(path)
	)

	// todo - this should be handled via Manager
	p.pluginMap[slug] = pl

	p.register(path)

	if rawPlugin, err = p.dispense2(slug, ff); err != nil {
		p.log.Error("could not dispense", zap.String("slug", slug))
	}

	return
}

func (p *plugin) registerAutomations() {
	// for _, raw := range p.rawPlugins {
	// 	p.ar.AddFunctions(raw.(*cp).automationFunction())
	// }
}

func (p *plugin) RegisterAll(ctx context.Context) {
	p.registerPlugins(ctx, PLUGIN_PATH)
	p.registerAutomations()
}

func makeSlug(path string) string {
	return filepath.Base(path)
}

// satisfy bundler
func (p *plugin) Register(path string) error {
	// temp
	path = filepath.Join(path, "automation_function")
	p.registerPlugins(context.Background(), path)

	return nil
}

func (p *plugin) Dispense() error {
	p.registerAutomations()
	return nil
}

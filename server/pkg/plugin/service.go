package plugin

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	hcp "github.com/hashicorp/go-plugin"
	"go.uber.org/zap"
)

type (
	plugin struct {
		log *zap.Logger

		rawPlugins map[string]any
		clients    map[string]*hcp.Client
		pluginMap  map[string]hcp.Plugin
	}
)

var (
	pl *plugin
)

const PLUGIN_PATH = "./plugin"

func Service() *plugin {
	return pl
}

func Setup(l *zap.Logger) {
	if pl != nil {
		return
	}

	pl = New(l)
}

func New(l *zap.Logger) *plugin {
	return &plugin{
		log: l.Named("plugin"),

		rawPlugins: make(map[string]any),
		clients:    make(map[string]*hcp.Client),
		pluginMap:  make(map[string]hcp.Plugin),
	}
}

func (p *plugin) register(path string) {
	var (
		slug = makeSlug(path)

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

	p.clients[slug] = hcp.NewClient(config)

	p.log.Info("registered new GRPC client", zap.String("slug", slug))
}

func (p *plugin) Unregister() {
	for _, c := range p.clients {
		c.Kill()
	}
}

// dispense finds the registered client and prepares it in the internal registry
func (p *plugin) dispense(slug string, ff func() any) (plugin any, err error) {
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

	if err != nil {
		p.log.Error("client dispense error", zap.Error(err))
	}

	return
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

	p.pluginMap[slug] = pl
	p.register(path)

	if rawPlugin, err = p.dispense(slug, ff); err != nil {
		p.log.Error("could not dispense", zap.String("slug", slug))
	}

	return
}

func (p *plugin) DeregisterPlugin(ctx context.Context, path string) (err error) {
	if c, is := p.clients[makeSlug(path)]; is {
		c.Kill()
	}

	return
}

func makeSlug(path string) string {
	return filepath.Base(path)
}

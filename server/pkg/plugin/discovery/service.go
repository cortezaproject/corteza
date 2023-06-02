package discovery

import (
	"context"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

const DIR = "./bundle"

const (
	STATUS_VALID status = iota
	STATUS_INVALID
	STATUS_REGISTERED
)

type (
	status int8

	Discovery interface {
		AddBundler(...Bundler)
		Discover(context.Context, string) error
	}

	bundle struct {
		id     byte
		handle string
		path   string
		status status
		sub    []string
		meta   map[string]string
		s      map[byte]conf
	}

	conf struct {
		files []string
	}

	discovery struct {
		log *zap.Logger

		bundlers []Bundler
		handlers []PackageHandler
		bundles  map[string]*bundle
	}
)

var (
	pDsc *discovery
)

func New(l *zap.Logger) *discovery {
	return &discovery{
		log: l.Named("plugin"),

		bundles:  make(map[string]*bundle),
		bundlers: make([]Bundler, 0),
		handlers: make([]PackageHandler, 0),
	}
}

func Setup(l *zap.Logger) {
	if pDsc != nil {
		return
	}

	pDsc = New(l)
}

func Service() *discovery {
	return pDsc
}

// todo
func (b *discovery) Dispense(ctx context.Context) {
	b.log.Info("disconnecting plugins")

	for _, bnd := range b.bundles {
		for _, bb := range b.bundlers {
			sett := bnd.s[bb.Type()]
			err := bb.Deregister(ctx, sett.files)

			if err != nil {
				b.log.Error("could not disconnect plugin", zap.Error(err), zap.Any("bundler", bb.Type()))
			}
		}
	}
}

func (b *discovery) Discover(ctx context.Context, path string) (err error) {
	// temp
	b.handlers = append(b.handlers, packageHandlerFs{})

	ss, err := os.ReadDir(path)

	if err != nil {
		return
	}

	for _, s := range ss {
		for _, h := range b.handlers {
			if !h.Match(s) {
				continue
			}

			err = h.Validate(s, path)

			if err != nil {
				continue
			}

			// if u, is := h.(Unboxer); is {
			// 	if err = u.Unbox(); err != nil {
			// 		continue
			// 	}
			// }

			pp := filepath.Join(path, s.Name())
			ss := STATUS_VALID

			m, err := filepath.Glob(pp + "/*")

			if err != nil || len(m) == 0 {
				ss = STATUS_INVALID
			}

			b.bundles[s.Name()] = &bundle{
				handle: s.Name(),
				path:   pp,
				status: ss,
				sub:    m,
				meta:   make(map[string]string),
				s:      make(map[byte]conf),
			}
		}

	}

	return
}

func (b *discovery) validate(ctx context.Context) {
	for _, bb := range b.bundles {
		if bb.status != STATUS_VALID {
			continue
		}

		for _, h := range b.bundlers {
			c, list := h.Validate(ctx, bb.sub)

			bb.id |= c
			bb.s[c] = conf{
				files: list,
			}
		}
	}
}

func (b *discovery) register(ctx context.Context) {
	var (
		s   conf
		err error
	)

	for _, bnd := range b.bundles {
		if bnd.status != STATUS_VALID {
			continue
		}

		for _, bb := range b.bundlers {
			s = bnd.s[bb.Type()]
			err = bb.Register(ctx, s.files)

			if err != nil {
				b.log.Error("could not register bundler", zap.Error(err), zap.Any("bundler", bb.Type()))
			}

			bnd.status = STATUS_REGISTERED
		}
	}
}

func (b *discovery) AddBundler(bb ...Bundler) {
	b.bundlers = append(b.bundlers, bb...)
}

func (b *discovery) Activate(ctx context.Context) {
	b.log.Debug("starting plugin discovery ...")

	if err := b.Discover(ctx, DIR); err != nil {
		b.log.Error(err.Error())
		return
	}

	b.log.Debug("validating discovered plugins")
	b.validate(ctx)

	b.log.Debug("registering plugins")
	b.register(ctx)
}

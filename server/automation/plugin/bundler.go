package plugin

import (
	"context"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
)

// what does this one do?
//  - finds proper files under a specific dir
//  - pushes the found plugin to the parent registry service for a specific bundle

// get dir from bundle discovery
// find the specific dir inside it (automation_function)
// loop through all the files in it
// populate pluginmap
// create a new client (hashicorp)
// dispense plugin (raw plugin to an internal list)
// how do we register the automation functions? do we inject automation registry into this service?

// where do the bundles that are coming in via api gonna be persisted and how?

const DIR = "./bundle"

const (
	AUTOMATION_WORKFLOW = 1
	AUTOMATION_FUNCTION = 2
	EXPR_TYPE           = 4
	GATEWAY_FILTER      = 8
)

type (
	Bundler interface {
		Register([]string) error
		Type() byte
	}

	FsValidator interface {
		Validate([]string) (byte, []string)
	}

	bundle struct {
		handle string
		path   string
		status string
		sub    []string
		meta   map[string]string
		conf   byte
		s      map[byte]BundlerSettings
	}

	BundlerSettings struct {
		files []string
	}

	bundler struct {
		bundlers []Bundler
		handlers []PackageHandler
		bundles  map[string]*bundle
	}
)

var (
	bnd *bundler
)

func Service() *bundler {
	return bnd
}

func Setup(l *zap.Logger) {
	if bnd != nil {
		return
	}

	bnd = New(l)
}

func New(l *zap.Logger) *bundler {
	return &bundler{
		bundles:  make(map[string]*bundle),
		bundlers: make([]Bundler, 0),
		handlers: make([]PackageHandler, 0),
	}
}

func (b *bundler) Discover(ctx context.Context, path string) (err error) {
	spew.Dump(">>>> discover")
	// temp
	b.handlers = append(b.handlers, packageHandlerFs{})

	// find via path
	// ss, err := filepath.Glob(filepath.Join(path, "*"))
	ss, err := os.ReadDir(path)

	if err != nil {
		spew.Dump(err)
		return
	}

	for _, s := range ss {
		for _, h := range b.handlers {
			if !h.Match(s) {
				continue
			}

			spew.Dump("matched ", s.Name())

			err = h.Validate(s, path)

			if err != nil {
				spew.Dump("could not validate", err)
				continue
			}

			if u, is := h.(Unboxer); is {
				if err = u.Unbox(); err != nil {
					spew.Dump("could not unbox")
					continue
				}
			}

			pp := filepath.Join(path, s.Name())
			ss := "valid"

			m, err := filepath.Glob(pp + "/*")

			if err != nil || len(m) == 0 {
				ss = "invalid"
			}

			b.bundles[s.Name()] = &bundle{
				handle: s.Name(),
				path:   pp,
				status: ss,
				sub:    m,
				meta:   make(map[string]string),
				s:      make(map[byte]BundlerSettings),
			}
		}

	}

	return
}

func (b *bundler) Validate(ctx context.Context) {
	for _, bb := range b.bundles {
		if bb.status != "valid" {
			continue
		}

		for _, h := range b.bundlers {
			if vv, is := h.(FsValidator); is {
				c, list := vv.Validate(bb.sub)

				bb.conf |= c
				bb.s[c] = BundlerSettings{
					files: list,
				}
			}
		}
	}

	spew.Dump(b.bundles)
}

func (b *bundler) Register(ctx context.Context) {
	for _, bnd := range b.bundles {
		for _, bb := range b.bundlers {
			sett := bnd.s[bb.Type()]
			spew.Dump("setings", sett)
			err := bb.Register(sett.files)

			if err != nil {
				spew.Dump(err)
			}
		}
	}
}

// func (b *bundler) Dispense() {
// 	for _, _ = range b.bundles {
// 		for _, bb := range b.bundlers {
// 			_ = bb.Dispense()
// 		}
// 	}
// }

// register a bundler
func (b *bundler) registerBundler(bb ...Bundler) {
	b.bundlers = append(b.bundlers, bb...)
}

func (b *bundler) AddBundler(bb ...Bundler) {
	b.registerBundler(bb...)
}

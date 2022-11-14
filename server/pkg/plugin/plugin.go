package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"
)

type (
	// Set of plugins
	Set []*item

	// item represents a plugin
	item struct {
		src string
		def interface{}
	}
)

// Resolve string with colon separated paths
func Resolve(paths string) (out []string, err error) {
	var (
		matches []string
	)

	for _, part := range strings.Split(paths, ":") {
		matches, err = filepath.Glob(part)

		if err != nil {
			return
		}

		out = append(out, matches...)
	}

	return
}

// Load loads plugins from all given paths
func Load(paths ...string) (Set, error) {
	var set = Set{}

	for _, path := range paths {

		if info, err := os.Lstat(path); err != nil {
			return nil, err
		} else if info.IsDir() {
			return nil, fmt.Errorf("can not use directory %s as a plugin", path)
		}

		if i, err := load(path); err != nil {
			return nil, err
		} else {
			set = append(set, i)
		}
	}

	return set, nil
}

// load single plugin from the given path
func load(path string) (i *item, err error) {
	i = &item{}
	p, err := plugin.Open(path)
	if err != nil {
		return
	}

	aux, err := p.Lookup("CortezaPlugin")
	if err != nil {
		return
	}

	fn, is := aux.(func() interface{})
	if !is {
		return nil, fmt.Errorf("incompatible plugin definition")
	}

	i.def = fn()
	return
}

package jsenv

import (
	"fmt"

	"github.com/evanw/esbuild/pkg/api"
)

const (
	// limiting the support for loaders and targets
	LoaderJS  TransformLoader = TransformLoader(api.LoaderJS)
	LoaderJSX TransformLoader = TransformLoader(api.LoaderJSX)
	LoaderTS  TransformLoader = TransformLoader(api.LoaderTS)

	TargetNoop   TransformTarget = 0
	TargetES5    TransformTarget = TransformTarget(api.ES5)
	TargetES2016 TransformTarget = TransformTarget(api.ES2016)
)

type (
	TransformLoader uint8
	TransformTarget uint8

	t struct {
		ldr TransformLoader
		tr  TransformTarget
	}

	noop struct{}

	Transformer interface {
		Transform(string) ([]byte, error)
	}
)

func NewTransformer(loader TransformLoader, target TransformTarget) Transformer {
	if target == TargetNoop {
		return &noop{}
	}

	return &t{
		ldr: loader,
		tr:  target,
	}
}

// Transform uses the loaders and targets and transpiles
func (tt t) Transform(p string) (b []byte, err error) {
	result := api.Transform(p, api.TransformOptions{
		Loader: api.Loader(tt.ldr),
		Target: api.Target(tt.tr),
	})

	if len(result.Errors) > 0 {
		return []byte{}, fmt.Errorf("could not transform payload: %s", result.Errors[0].Text)
	}

	return result.Code, nil
}

// Fallback transform that keeps the original intact
func (tt noop) Transform(p string) ([]byte, error) {
	return []byte(p), nil
}

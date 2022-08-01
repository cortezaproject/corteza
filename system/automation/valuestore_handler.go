package automation

import (
	"context"
)

type (
	envGetter interface {
		Env(k string) (v any)
	}

	valuestoreHandler struct {
		reg    valuestoreHandlerRegistry
		getter envGetter
	}
)

func ValuestoreHandler(reg valuestoreHandlerRegistry, getter envGetter) *valuestoreHandler {
	h := &valuestoreHandler{
		reg:    reg,
		getter: getter,
	}

	h.register()
	return h
}

func (h valuestoreHandler) env(ctx context.Context, args *valuestoreEnvArgs) (results *valuestoreEnvResults, err error) {
	results = &valuestoreEnvResults{
		Value: h.getter.Env(args.Key),
	}
	return
}

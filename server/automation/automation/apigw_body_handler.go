package automation

import (
	"context"
	"io"
)

type (
	apigwBodyHandler struct {
		reg apigwBodyHandlerRegistry
	}
)

func ApigwBodyHandler(reg queueHandlerRegistry) *apigwBodyHandler {
	h := &apigwBodyHandler{
		reg: reg,
	}

	h.register()
	return h
}

func (h apigwBodyHandler) read(ctx context.Context, args *apigwBodyReadArgs) (res *apigwBodyReadResults, err error) {
	res = &apigwBodyReadResults{}

	bb, err := io.ReadAll(args.Request.Body)

	if err != nil {
		return
	}

	res.Body = string(bb)

	return
}

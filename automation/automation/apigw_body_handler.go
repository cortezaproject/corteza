package automation

import (
	"context"
	"fmt"
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

	if !args.hasRequest {
		err = fmt.Errorf("could not read body, contents missing")
		return
	}

	bb, err := io.ReadAll(args.Request.Body)

	if err != nil {
		return
	}

	res.Body = string(bb)

	return
}

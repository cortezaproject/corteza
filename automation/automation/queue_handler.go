package automation

import (
	"context"
	"fmt"

	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/cortezaproject/corteza-server/pkg/messagebus"
)

type (
	queueHandler struct {
		reg queueHandlerRegistry
	}
)

func QueueHandler(reg queueHandlerRegistry) *queueHandler {
	h := &queueHandler{
		reg: reg,
	}

	h.register()
	return h
}

func (h queueHandler) write(ctx context.Context, args *queueWriteArgs) (err error) {
	if !args.hasQueue {
		return fmt.Errorf("could not send message to queue, queue empty")
	}

	if !args.hasPayload {
		return fmt.Errorf("could not send message to queue, payload empty")
	}

	p, err := expr.CastToString(args.Payload)

	if err != nil {
		return err
	}

	go func() {
		queue := messagebus.Service().Write(args.Queue)
		queue <- []byte(p)
	}()

	return nil
}

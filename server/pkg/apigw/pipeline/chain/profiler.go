package chain

import (
	"net/http"

	"github.com/cortezaproject/corteza/server/pkg/apigw/profiler"
)

type (
	chiHandlerChainProfiler struct {
		chain middleware
	}
)

func NewProfiler() *chiHandlerChainProfiler {
	return &chiHandlerChainProfiler{}
}

func (cc *chiHandlerChainProfiler) Chain(m middleware) {
	cc.chain = m
}

func (cc *chiHandlerChainProfiler) Handler() http.Handler {
	// wrap handlers around the chain so we get some profiling info
	cc.chain = append([]func(http.Handler) http.Handler{profiler.StartHandler}, cc.chain...)
	cc.chain = append(cc.chain, profiler.FinishHandler)

	return chiHandler(cc.chain)
}

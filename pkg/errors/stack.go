package errors

import (
	"runtime"
	"strings"
)

type (
	frame struct {
		Func string `json:"func"`
		File string `json:"file"`
		Line int    `json:"line"`
	}
)

const (
	stackBuf = 128
)

func collectStack(skip int) []*frame {
	var (
		stack     = make([]*frame, 0, stackBuf)
		cc        = make([]uintptr, stackBuf)
		collected = runtime.Callers(skip, cc)
		frames    = runtime.CallersFrames(cc[:collected])
	)

	for {
		f, n := frames.Next()
		if !n {
			break
		}

		aux := &frame{
			Func: f.Function,
			File: f.File,
			Line: f.Line,
		}

		if li := strings.LastIndex(f.Function, "/"); li > 0 {
			// remove prefix
			aux.Func = aux.Func[li+1:]
		}

		stack = append(stack, aux)
	}

	return stack
}

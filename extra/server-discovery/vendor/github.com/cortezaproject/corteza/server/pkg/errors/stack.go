package errors

import (
	"runtime"
	"strconv"
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

// Converts node stack trace (from Error().stack) to internal structure
//
// Node stack traces are received with errors from Corredor (node.js) automation server
func convertNodeStack(nf []string) []*frame {
	conv := make([]*frame, len(nf))
	for i := range nf {
		var (
			f         = &frame{}
			file, lno string
		)

		if p := strings.Index(nf[i], " ("); p > 0 {
			f.Func = nf[i][:p]

			file = nf[i][p+2 : len(nf[i])-1]
		} else {
			file = nf[i]
		}

		parts := strings.Split(file, ":")
		f.File = parts[0]

		if len(parts) >= 2 {
			lno = parts[1]
		}

		if lno != "" {
			f.Line, _ = strconv.Atoi(lno)
		}

		conv[i] = f
	}
	return conv
}

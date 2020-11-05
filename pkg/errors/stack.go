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

// Converts node stack trace (from Error().stack) to internal structure
//
// Node stack traces are received with errors from Corredor (node.js) automation server
func convertNodeStack(stack []string) []*frame {
	// @todo

	//(metadata.MD) (len=2) {
	// (string) (len=12) "content-type": ([]string) (len=1 cap=1) {
	//  (string) (len=16) "application/grpc"
	// },
	// (string) (len=5) "stack": ([]string) (len=10 cap=10) {
	//  (string) (len=93) "Object.exec (/Users/darh/dev/corteza/server-corredor/usr/testing/server-scripts/foo.js:12:11)",
	//  (string) (len=127) "Object.<anonymous> (/Users/darh/dev/corteza/server-corredor/node_modules/@cortezaproject/corteza-js/src/corredor/exec.ts:26:35)",
	//  (string) (len=28) "Generator.next (<anonymous>)",
	//  (string) (len=100) "/Users/darh/dev/corteza/server-corredor/node_modules/@cortezaproject/corteza-js/dist/index.js:277:71",
	//  (string) (len=25) "new Promise (<anonymous>)",
	//  (string) (len=112) "__awaiter (/Users/darh/dev/corteza/server-corredor/node_modules/@cortezaproject/corteza-js/dist/index.js:273:12)",
	//  (string) (len=114) "Object.Exec (/Users/darh/dev/corteza/server-corredor/node_modules/@cortezaproject/corteza-js/dist/index.js:483:12)",
	//  (string) (len=97) "ServerScripts.exec (/Users/darh/dev/corteza/server-corredor/src/services/server-scripts.ts:86:17)",
	//  (string) (len=96) "Object.Exec (/Users/darh/dev/corteza/server-corredor/src/grpc-handlers/server-scripts.ts:142:11)",
	//  (string) (len=78) "/Users/darh/dev/corteza/server-corredor/node_modules/grpc/src/server.js:593:13"
	// }
	//}

	return nil
}

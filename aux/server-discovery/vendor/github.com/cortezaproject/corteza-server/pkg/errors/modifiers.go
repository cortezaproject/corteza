package errors

type (
	mfn func(*Error)
)

// Adds meta
func Meta(k, v interface{}) mfn {
	return func(e *Error) {
		if e.meta == nil {
			e.meta = meta{}
		}

		e.meta[k] = v
	}
}

// Trim all keys from meta
func MetaTrim(kk ...interface{}) mfn {
	return func(e *Error) {
		for _, k := range kk {
			delete(e.meta, k)
		}
	}
}

// StackSkip skips n frames in the stack
//
func StackSkip(n int) mfn {
	return func(e *Error) {
		if n > len(e.stack) {
			e.stack = nil
		}

		e.stack = e.stack[n:]
	}
}

// StackTrim removes n frames from the end of the stack
func StackTrim(n int) mfn {
	return func(e *Error) {
		if len(e.stack) < n {
			e.stack = nil
		}

		e.stack = e.stack[:len(e.stack)-n]
	}
}

// StackTrimAtFn removes all frames after (inclusive) the (1st) function match
func StackTrimAtFn(fn string) mfn {
	return func(e *Error) {
		for i := range e.stack {
			if i > 2 && e.stack[i].Func == fn {
				e.stack = e.stack[:i]
				break
			}
		}
	}
}

// Wrap embeds error
func Wrap(w error) mfn {
	return func(e *Error) {
		e.wrap = w
	}
}

// Converts and attaches node.js stack to error
func AddNodeStack(stack []string) mfn {
	return func(e *Error) {
		e.stack = convertNodeStack(stack)
	}
}

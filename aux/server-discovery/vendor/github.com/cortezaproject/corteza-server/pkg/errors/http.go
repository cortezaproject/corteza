package errors

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/cortezaproject/corteza-server/pkg/locale"
)

type (
	translatable interface {
		Translate(tr func(string, string, ...string) string) error
	}
)

// ServeHTTP Prepares and encodes given error for HTTP transport
//
// mask arg hides extra/debug info
//
// Proper HTTP status codes are generally not used in the API due to compatibility issues
// This should be addressed in the future versions when/if we restructure the API
func ServeHTTP(w http.ResponseWriter, r *http.Request, err error, mask bool) {
	// due to backward compatibility,
	// custom HTTP statuses are disabled for now.
	serveHTTP(w, r, http.StatusOK, err, mask)
}

// ProperlyServeHTTP Prepares and encodes given error for HTTP transport, same as ServeHTTP but with proper status codes
func ProperlyServeHTTP(w http.ResponseWriter, r *http.Request, err error, mask bool) {
	var (
		code = http.StatusInternalServerError
	)

	if e, is := err.(*Error); is {
		code = e.kind.httpStatus()
	}

	serveHTTP(w, r, code, err, mask)
}

// Serves error via
func serveHTTP(w http.ResponseWriter, r *http.Request, code int, err error, mask bool) {
	var (
		// Very naive approach on parsing accept headers
		acceptsJson = strings.Contains(r.Header.Get("accept"), "application/json")
	)

	if !mask && !acceptsJson {
		// Prettify error for plain text debug output
		w.Header().Set("Content-Type", "plain/text")
		w.WriteHeader(code)
		writeHttpPlain(w, err, mask)
		_, _ = fmt.Fprintln(w, "Note: you are seeing this because system is running in development mode")
		_, _ = fmt.Fprintln(w, "and HTTP request is made without \"Accept: .../json\" headers")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	writeHttpJSON(r.Context(), w, err, mask)
}

func writeHttpPlain(w io.Writer, err error, mask bool) {
	var (
		dlmt = strings.Repeat("-", 80)

		write = func(s string, a ...interface{}) {
			if _, err := fmt.Fprintf(w, s, a...); err != nil {
				panic(fmt.Errorf("failed to write HTTP response: %w", err))
			}
		}
	)

	write("Error: ")
	write(err.Error())
	write("\n")

	if _, is := err.(interface{ Safe() bool }); !is || mask {
		// do not output any details on un-safe errors or
		// when a masked output is preferred
		return
	}

	if err, is := err.(*Error); is {
		write(dlmt + "\n")

		var (
			ml, kk = err.meta.StringKeys()
		)

		sort.Strings(kk)

		for _, key := range kk {
			write("%s:", key)
			write(strings.Repeat(" ", ml-len(key)+1))
			write("%v\n", err.meta[key])
		}

		if len(err.stack) > 0 {
			write(dlmt + "\n")
			write("Call stack:\n")
			for i, f := range err.stack {
				write("%3d. %s:%d\n     %s()\n", len(err.stack)-i, f.File, f.Line, f.Func)
			}
		}

		if we := err.Unwrap(); we != nil {
			write(dlmt + "\n")
			write("Wrapped error:\n\n")
			writeHttpPlain(w, we, mask)
		}
	}
	write(dlmt + "\n")

}

// writeHttpJSON
func writeHttpJSON(ctx context.Context, w io.Writer, err error, mask bool) {
	var (
		wrap = struct {
			Error interface{} `json:"error"`
		}{}
	)

	if se, is := err.(interface{ Safe() bool }); !is || !se.Safe() || mask {
		// trim error details when not debugging or error is not safe or maske
		err = fmt.Errorf(err.Error())
	}

	// if error is translatable, pass in the lambda that returns
	// translated error
	//
	// we're using global locale service directly for now
	if tr, is := err.(translatable); is {
		err = tr.Translate(func(ns string, key string, pairs ...string) string {
			return locale.Global().T(ctx, ns, key, pairs...)
		})
	}

	if c, is := err.(json.Marshaler); is {
		// take advantage of JSON marshaller on error
		wrap.Error = c
	} else {
		wrap.Error = map[string]string{"message": err.Error()}
	}

	if err = json.NewEncoder(w).Encode(wrap); err != nil {
		panic(fmt.Errorf("failed to encode error: %w", err))
	}
}

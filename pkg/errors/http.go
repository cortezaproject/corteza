package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

// ServeHTTP Prepares and encodes given error for HTTP transport
//
// mask arg hides extra/debug info
func ServeHTTP(w http.ResponseWriter, r *http.Request, err error, debug bool) {
	var (
		// Very naive approach on parsing accept headers
		acceptsJson = strings.Contains(r.Header.Get("accept"), "application/json")

		// due to backward compatibility,
		// proper use of HTTP statuses is disabled for now.
		code = http.StatusOK
		//code = http.StatusInternalServerError
	)

	// due to backward compatibility,
	// custom HTTP statuses are disabled for now.
	//if e, is := err.(*Error); is {
	//	code = e.kind.httpStatus()
	//}

	w.WriteHeader(code)

	if !debug {
		// trim error details when not debugging
		err = fmt.Errorf(err.Error())
	}

	if debug && !acceptsJson {
		// Prettify error for plain text debug output
		w.Header().Set("Content-Type", "plain/text")
		writeHttpPlain(w, err)
		fmt.Fprintln(w, "Note: you are seeing this because system is running in development mode")
		fmt.Fprintln(w, "and HTTP request is made without \"Accept: .../json\" headers")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	writeHttpJSON(w, err)
}

func writeHttpPlain(w io.Writer, err error) {
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
			writeHttpPlain(w, we)
		}
	}
	write(dlmt + "\n")

}

func writeHttpJSON(w io.Writer, err error) {
	var (
		wrap = struct {
			Error interface{} `json:"error"`
		}{}
	)

	if c, is := err.(*Error); is {
		wrap.Error = c
	} else {
		wrap.Error = map[string]string{"message": err.Error()}
	}

	if err = json.NewEncoder(w).Encode(wrap); err != nil {
		panic(fmt.Errorf("failed to encode error: %w", err))
	}
}

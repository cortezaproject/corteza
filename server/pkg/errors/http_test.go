package errors

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/require"
)

func Example_writeHttpPlain() {
	writeHttpPlain(os.Stdout, fmt.Errorf("dummy error"), true)

	// Output:
	// Error: dummy error
}

func Example_writeHttpJSON() {
	writeHttpJSON(context.Background(), os.Stdout, fmt.Errorf("dummy error"), true)

	// Output:
	// {"error":{"message":"dummy error"}}
}

func Example_writeHttpJSON_clientAbortedConnectionReset() {
	writeHttpJSON(context.Background(), os.Stdout, syscall.ECONNRESET, true)

	// Output:
}

func Example_writeHttpPlain_clientAbortedConnectionReset() {
	writeHttpPlain(os.Stdout, syscall.ECONNRESET, true)

	// Output:
}

func Example_writeHttpJSON_clientAbortedConnectionPipe() {
	writeHttpJSON(context.Background(), os.Stdout, syscall.EPIPE, true)

	// Output:
}

func Example_writeHttpPlain_clientAbortedConnectionPipe() {
	writeHttpPlain(os.Stdout, syscall.EPIPE, true)

	// Output:
}

func Example_writeHttpPlain_masked() {
	err := New(0, "dummy error", Meta("a", "b"), Meta(&Error{}, "nope"))
	err.stack = nil // will not test the stack as file path & line numbers might change
	writeHttpPlain(os.Stdout, err, true)
	// Output:
	// Error: dummy error
}

func Example_writeHttpPlain_unmasked() {
	err := New(0, "dummy error", Meta("a", "b"), Meta(&Error{}, "nope"))
	err.stack = nil // will not test the stack as file path & line numbers might change
	writeHttpPlain(os.Stdout, err, false)
	// Output:
	// Error: dummy error
	// --------------------------------------------------------------------------------
	// a: b
	// --------------------------------------------------------------------------------
}

func Test_writeHttpJSON(t *testing.T) {
	var (
		err = New(0, "dummy error", Meta("meta", "meta"))
		buf = bytes.NewBuffer(nil)
		req = require.New(t)
	)

	buf.Truncate(0)
	writeHttpJSON(context.Background(), buf, err, false)
	req.Contains(buf.String(), "dummy error")
	req.Contains(buf.String(), "meta")
	req.Contains(buf.String(), "stack")

	// when errors are masked (production env) we do not add meta or stack
	buf.Truncate(0)
	writeHttpJSON(context.Background(), buf, err, true)
	req.Contains(buf.String(), "dummy error")
	req.NotContains(buf.String(), "meta")
	req.NotContains(buf.String(), "stack")
}

// A workflow error step's title is meant for the user, so it survives masking;
// other meta and the stack do not
func Test_writeHttpJSON_automationTitle(t *testing.T) {
	var (
		err = Automation("Order is not ready").Apply(Meta("title", "Order check"), Meta("secret", "s3cr3t"))
		buf = bytes.NewBuffer(nil)
		req = require.New(t)
	)

	writeHttpJSON(context.Background(), buf, err, true)
	req.JSONEq(`{"error":{"message":"Order is not ready","meta":{"title":"Order check"}}}`, buf.String())

	// other kinds do not keep a title
	buf.Truncate(0)
	writeHttpJSON(context.Background(), buf, New(KindInternal, "boom", Meta("title", "Leaked")), true)
	req.JSONEq(`{"error":{"message":"boom"}}`, buf.String())
}

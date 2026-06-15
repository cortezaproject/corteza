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

// Locks the masking matrix for KindAutomation errors so the
// MetaWorkflowErrorSafe opt-in bypass semantics cannot be
// silently changed by future refactors.
//
// Cells:
//  1. KindAutomation + MetaWorkflowErrorSafe=true  + masked → meta preserved
//  2. KindAutomation + no flag                     + masked → meta stripped
//  3. KindInternal  + MetaWorkflowErrorSafe=true  + masked → meta stripped
//     (flag is a no-op on non-automation kinds)
//  4. KindAutomation + MetaWorkflowErrorSafe=true  + unmasked → meta preserved
func Test_writeHttpJSON_workflowErrorSafeBypass(t *testing.T) {
	req := require.New(t)

	// 1. automation + safe flag + masked → meta visible
	safe := Automation("author message").Apply(
		Meta(MetaWorkflowErrorSafe, true),
		Meta("workflow.error.title", "Upstream failure"),
	)
	buf := bytes.NewBuffer(nil)
	writeHttpJSON(context.Background(), buf, safe, true)
	req.Contains(buf.String(), "author message", "message must always be visible")
	req.Contains(buf.String(), "workflow.error.title", "safe flag must preserve meta under mask")
	req.Contains(buf.String(), "Upstream failure")

	// 2. automation + no flag + masked → meta stripped
	plain := Automation("runtime error").Apply(Meta("internal", "secret"))
	buf.Truncate(0)
	writeHttpJSON(context.Background(), buf, plain, true)
	req.Contains(buf.String(), "runtime error", "flat message must still be visible")
	req.NotContains(buf.String(), "internal", "unflagged automation meta must be stripped under mask")
	req.NotContains(buf.String(), "secret")

	// 3. non-automation kind + safe flag + masked → meta stripped
	//    (the flag is only honoured for KindAutomation — it must not
	//    leak meta from other error kinds even if the flag is set)
	nonAuto := Internal("internal boom").Apply(
		Meta(MetaWorkflowErrorSafe, true),
		Meta("leak", "pls no"),
	)
	buf.Truncate(0)
	writeHttpJSON(context.Background(), buf, nonAuto, true)
	req.Contains(buf.String(), "internal boom")
	req.NotContains(buf.String(), "leak", "safe flag on non-automation must be ignored")
	req.NotContains(buf.String(), "pls no")

	// 4. automation + safe flag + unmasked → meta visible (debug path)
	buf.Truncate(0)
	writeHttpJSON(context.Background(), buf, safe, false)
	req.Contains(buf.String(), "author message")
	req.Contains(buf.String(), "Upstream failure")
}

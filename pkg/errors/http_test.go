package errors

import (
	"context"
	"fmt"
	"os"
)

func Example_writeHttpPlain() {
	writeHttpPlain(os.Stdout, fmt.Errorf("dummy error"))

	// Output:
	// Error: dummy error
	// --------------------------------------------------------------------------------
}

func Example_writeHttpJSON() {
	writeHttpJSON(context.Background(), os.Stdout, fmt.Errorf("dummy error"), true)

	// Output:
	// {"error":{"message":"dummy error"}}
}

func Example_writeHttpPlain_2() {
	err := New(0, "dummy error", Meta("a", "b"), Meta(&Error{}, "nope"))
	err.stack = nil // will not test the stack as file path & line numbers might change
	writeHttpPlain(os.Stdout, err)
	// Output:
	// Error: dummy error
	// --------------------------------------------------------------------------------
	// a: b
	// --------------------------------------------------------------------------------
}

func Example_writeHttpJSON_2() {
	err := New(0, "dummy error", Meta("a", "b"), Meta(&Error{}, "nope"))
	err.stack = nil // will not test the stack as file path & line numbers might change
	writeHttpJSON(context.Background(), os.Stdout, err, false)

	// Output:
	// {"error":{"message":"dummy error","meta":{"a":"b"}}}
}

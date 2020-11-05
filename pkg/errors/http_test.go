package errors

import (
	"fmt"
	"os"
)

func ExampleSimpleErrorAsText() {
	writeHttpPlain(os.Stdout, fmt.Errorf("dummy error"))

	// Output:
	// Error: dummy error
	// --------------------------------------------------------------------------------
}

func ExampleSimpleErrorAsJson() {
	writeHttpJSON(os.Stdout, fmt.Errorf("dummy error"), true)

	// Output:
	// {"error":{"message":"dummy error"}}
}

func ExampleErrorAsText() {
	err := New(0, "dummy error", Meta("a", "b"), Meta(&Error{}, "nope"))
	err.stack = nil // will not test the stack as file path & line numbers might change
	writeHttpPlain(os.Stdout, err)
	// Output:
	// Error: dummy error
	// --------------------------------------------------------------------------------
	// a: b
	// --------------------------------------------------------------------------------
}

func ExampleErrorAsJson() {
	err := New(0, "dummy error", Meta("a", "b"), Meta(&Error{}, "nope"))
	err.stack = nil // will not test the stack as file path & line numbers might change
	writeHttpJSON(os.Stdout, err, false)

	// Output:
	// {"error":{"message":"dummy error","meta":{"a":"b"}}}
}

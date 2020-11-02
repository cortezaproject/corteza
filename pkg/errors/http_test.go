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
	// Note: you are seeing this because system is running in development mode
	// and HTTP request is made without "Accept: .../json" headers
}

func ExampleSimpleErrorAsJson() {
	writeHttpJSON(os.Stdout, fmt.Errorf("dummy error"))

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
	// Note: you are seeing this because system is running in development mode
	// and HTTP request is made without "Accept: .../json" headers
}

func ExampleErrorAsJson() {
	err := New(0, "dummy error", Meta("a", "b"), Meta(&Error{}, "nope"))
	err.stack = nil // will not test the stack as file path & line numbers might change
	writeHttpJSON(os.Stdout, err)

	// Output:
	// {"error":{"message":"dummy error","meta":{"a":"b"}}}
}

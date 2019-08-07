package logger

import (
	"fmt"
)

type (
	field struct {
		name  string
		value interface{}
	}
)

func NewField(name string, value interface{}) Field {
	return field{name, value}
}

func (f field) Name() string {
	return f.name
}

func (f field) Value() interface{} {
	return f.value
}

func (f field) String() string {
	return fmt.Sprintf("%s=%v", f.name, f.value)
}

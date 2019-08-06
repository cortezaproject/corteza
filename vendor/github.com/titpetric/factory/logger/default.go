package logger

import (
	"context"
	"log"
)

type (
	Default struct{}
)

func (Default) Log(ctx context.Context, msg string, fields ...Field) {
	args := make([]interface{}, len(fields)+1)
	args[0] = msg
	for key, value := range fields {
		args[key+1] = value
	}
	// fields satisfy fmt.Stringer as well
	log.Println(args...)
}

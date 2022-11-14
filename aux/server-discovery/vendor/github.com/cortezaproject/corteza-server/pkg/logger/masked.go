package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Mask is a passthrough for MaskIf
func Mask(key string, val interface{}) zap.Field {
	return MaskIf(key, val, true)
}

// MaskIf conditionally masks (replaces string with * of equal length input string/stringer
func MaskIf(key string, val interface{}, mask bool) zap.Field {
	var str string

	switch v := val.(type) {
	case fmt.Stringer:
		str = v.String()
	case string:
		str = v
	default:
		// empty string for other types
		str = ""
	}

	if mask {
		str = strings.Repeat("*", len(str))
	}

	return zap.Field{
		Key:    key,
		Type:   zapcore.StringType,
		String: str,
	}
}

package request

import (
	"testing"
)

func Test_parseBool(t *testing.T) {
	truthies := []string{"y", "yes", "TRUE", "true", "t", "1", " T "}
	falsies := []string{"a", "FALSE", "tr", "11111", " FALSE ", "0"}

	for _, truth := range truthies {
		if !parseBool(truth) {
			t.Errorf("Must parse '%s' as boolean value TRUE", truth)
		}
	}

	for _, falsie := range falsies {
		if parseBool(falsie) {
			t.Errorf("Must not parse '%s' as boolean value TRUE", falsie)
		}
	}

}

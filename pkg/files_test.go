package files

import (
	"fmt"
	"testing"

	"github.com/cortezaproject/corteza-server/internal/test"
)

func TestExtractExtFromURL(t *testing.T) {
	tests := [][]string{
		[]string{"https://www.d.tld/x.png", "png"},
		[]string{"https://www.d.tld/x.png?x?y", "png"},
		[]string{"https://www.d.tld/x.png#a#b", "png"},
		[]string{"https://www.d.tld/x.png#a?b#c?d", "png"},
		[]string{"x.png", "png"},
	}

	for _, tc := range tests {
		ext, _ := ExtractExtFromURL(tc[0])
		test.Assert(t, ext == tc[1], fmt.Sprintf("Invalid ext %s, should be %s", ext, tc[1]))
	}

	_, err := ExtractExtFromURL("invalid")
	test.Assert(t, err != nil, "Should return err")
}

func TestExtractNameFromURL(t *testing.T) {
	tests := [][]string{
		[]string{"https://www.d.tld/x.png", "x.png"},
		[]string{"https://www.d.tld/x.png?x?y", "x.png"},
		[]string{"https://www.d.tld/x.png#a#b", "x.png"},
		[]string{"https://www.d.tld/x.png#a?b#c?d", "x.png"},
		[]string{"x.png", "x.png"},
	}

	for _, tc := range tests {
		ext, _ := ExtractNameFromURL(tc[0])
		test.Assert(t, ext == tc[1], fmt.Sprintf("Invalid ext %s, should be %s", ext, tc[1]))
	}

	_, err := ExtractNameFromURL("invalid")
	test.Assert(t, err != nil, "Should return err")
}

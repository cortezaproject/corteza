package envoy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRefIDTerms(t *testing.T) {
	req := require.New(t)

	req.Empty(refIDTerms(nil))
	req.Empty(refIDTerms([]string{"", "0", "abc", "1 OR 1=1"}))
	req.Equal(
		[]string{"recordID=123", "recordID=456"},
		refIDTerms([]string{"", "123", "0", "456"}),
	)
}

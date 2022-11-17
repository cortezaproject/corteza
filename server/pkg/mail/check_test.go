package mail

import (
    "github.com/stretchr/testify/require"
    "testing"
)

func TestHostValidator(t *testing.T) {
    ttc := []struct {
        host string
        ok   bool
    }{
        {"ç$€§az.com", false},
        {"@sendyy.com", false},
        {"qwertyuiop.com", true},
        {"test.foo.bar", true},
        {"10.10.10", true},
        {"192.10.10345", true},
        {"1rg.10ui.10", true},
        {"info .crust tech", false},
        {"info.crust.tech", true},
        {"crust-tech?", false},
        {"crust-tech", true},
        {"crust/tech", false},
    }

    for _, tc := range ttc {
        require.True(t, IsValidHost(tc.host) == tc.ok, "Validation of %s should return %v", tc.host, tc.ok)
    }
}

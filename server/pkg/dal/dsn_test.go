package dal

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestParseDSN(t *testing.T) {
	dsn := "http://user:pass@example.com/api?timeout=10s&max_idle_conns=50&tls_insecure=true"

	out, err := ParseDSN(dsn)
	require.NoError(t, err)

	require.Equal(t, "http", out.Scheme)
	require.Equal(t, "example.com", out.Host)
	require.Equal(t, "80", out.Port)
	require.Equal(t, "basic", out.AuthType)
	require.Equal(t, "user", out.Username)
	require.Equal(t, "pass", out.Password)

	require.Equal(t, 10*time.Second, out.Timeout)
	require.Equal(t, 50, out.MaxIdleConns)
}

func TestParseDSN_headers(t *testing.T) {
	dsn := "http://example.com?header.X-Test=one&header.X-Test=two"

	out, err := ParseDSN(dsn)
	require.NoError(t, err)

	values := out.Headers["X-Test"]
	require.Len(t, values, 2)
	require.Equal(t, "one", values[0])
	require.Equal(t, "two", values[1])
}

func TestParseDSN_Arbitrary(t *testing.T) {
	dsn := `http://example.com?arbitrary={"foo":"bar","n":1}`

	out, err := ParseDSN(dsn)
	require.NoError(t, err)

	require.Equal(t, "bar", out.Arbitrary["foo"])
}

func TestBaseURL(t *testing.T) {
	dsn := "rests://example.com:8443/api?a=b"

	out, err := ParseDSN(dsn)
	require.NoError(t, err)

	require.Equal(t, "https://example.com:8443/api", out.BaseURL())
}

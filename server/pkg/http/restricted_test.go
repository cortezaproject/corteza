package http

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBlockedIP(t *testing.T) {
	var (
		cases = []struct {
			ip        string
			linkLocal bool
			private   bool
		}{
			// public
			{ip: "93.184.216.34"},
			{ip: "2606:2800:220:1:248:1893:25c8:1946"},

			// cloud metadata endpoints
			{ip: "169.254.169.254", linkLocal: true, private: true},
			{ip: "169.254.170.2", linkLocal: true, private: true},
			{ip: "100.100.100.200", linkLocal: true, private: true},
			{ip: "fd00:ec2::254", linkLocal: true, private: true},
			{ip: "fe80::1", linkLocal: true, private: true},

			// link-local hidden inside IPv6
			{ip: "::ffff:169.254.169.254", linkLocal: true, private: true},
			{ip: "64:ff9b::a9fe:a9fe", linkLocal: true, private: true},
			{ip: "2002:a9fe:a9fe::1", linkLocal: true, private: true},

			// internal networks
			{ip: "127.0.0.1", private: true},
			{ip: "::1", private: true},
			{ip: "10.0.0.5", private: true},
			{ip: "172.16.10.1", private: true},
			{ip: "192.168.1.1", private: true},
			{ip: "100.64.0.1", private: true},
			{ip: "fc00::1", private: true},
			{ip: "0.0.0.0", private: true},
			{ip: "::", private: true},
			{ip: "224.0.0.1", linkLocal: true, private: true},
			{ip: "239.1.1.1", private: true},
			{ip: "::ffff:10.0.0.5", private: true},
			{ip: "64:ff9b::a00:5", private: true},
		}
	)

	for _, c := range cases {
		t.Run(c.ip, func(t *testing.T) {
			ip := net.ParseIP(c.ip)
			require.NotNil(t, ip)

			require.False(t, blockedIP(RestrictNone, ip))
			require.Equal(t, c.linkLocal, blockedIP(RestrictLinkLocal, ip), "link-local policy")
			require.Equal(t, c.private, blockedIP(RestrictPrivate, ip), "private policy")
		})
	}
}

func TestParseNetworkRestriction(t *testing.T) {
	for in, exp := range map[string]NetworkRestriction{
		"":           RestrictLinkLocal,
		"link-local": RestrictLinkLocal,
		"LINK-LOCAL": RestrictLinkLocal,
		" private ":  RestrictPrivate,
		"none":       RestrictNone,
	} {
		out, err := ParseNetworkRestriction(in)
		require.NoError(t, err, in)
		require.Equal(t, exp, out, in)
	}

	// unknown values must not silently disable the protection
	out, err := ParseNetworkRestriction("nope")
	require.Error(t, err)
	require.Equal(t, RestrictLinkLocal, out)
}

func TestRestrictedClient(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("internal"))
	}))
	defer srv.Close()

	t.Run("loopback is reachable by default", func(t *testing.T) {
		SetupRestricted(0, false, RestrictLinkLocal)

		rsp, err := RestrictedClient().Get(srv.URL)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rsp.StatusCode)
		_ = rsp.Body.Close()
	})

	t.Run("loopback is blocked with private restriction", func(t *testing.T) {
		SetupRestricted(0, false, RestrictPrivate)
		defer SetupRestricted(0, false, RestrictLinkLocal)

		_, err := RestrictedClient().Get(srv.URL)
		require.ErrorIs(t, err, ErrRestrictedNetwork)
	})

	t.Run("hostname resolving to blocked address", func(t *testing.T) {
		SetupRestricted(0, false, RestrictPrivate)
		defer SetupRestricted(0, false, RestrictLinkLocal)

		_, port, err := net.SplitHostPort(srv.Listener.Addr().String())
		require.NoError(t, err)

		_, err = RestrictedClient().Get("http://localhost:" + port)
		require.ErrorIs(t, err, ErrRestrictedNetwork)
	})

	t.Run("client taken before the setup follows new restriction", func(t *testing.T) {
		SetupRestricted(0, false, RestrictLinkLocal)
		early := RestrictedClient()

		SetupRestricted(0, false, RestrictPrivate)
		defer SetupRestricted(0, false, RestrictLinkLocal)

		_, err := early.Get(srv.URL)
		require.ErrorIs(t, err, ErrRestrictedNetwork)
	})

	t.Run("metadata endpoint is blocked by default", func(t *testing.T) {
		SetupRestricted(0, false, RestrictLinkLocal)

		_, err := RestrictedClient().Get("http://169.254.169.254/latest/meta-data/")
		require.ErrorIs(t, err, ErrRestrictedNetwork)
	})

	t.Run("redirects are not followed", func(t *testing.T) {
		SetupRestricted(0, false, RestrictLinkLocal)

		redir := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "http://169.254.169.254/latest/meta-data/", http.StatusFound)
		}))
		defer redir.Close()

		rsp, err := RestrictedClient().Get(redir.URL)
		require.NoError(t, err)
		require.Equal(t, http.StatusFound, rsp.StatusCode)
		_ = rsp.Body.Close()
	})
}

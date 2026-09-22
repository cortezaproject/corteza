package http

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
)

type (
	// NetworkRestriction controls what networks restricted HTTP client can connect to
	NetworkRestriction string

	// restrictedRoundTripper uses the latest configured transport
	restrictedRoundTripper struct{}
)

const (
	// RestrictNone allows connections to any address
	RestrictNone NetworkRestriction = "none"

	// RestrictLinkLocal blocks link-local addresses and well known cloud metadata endpoints
	RestrictLinkLocal NetworkRestriction = "link-local"

	// RestrictPrivate blocks everything that is not a public address
	RestrictPrivate NetworkRestriction = "private"
)

var (
	ErrRestrictedNetwork = errors.New("connection to restricted network address is not allowed")

	// one client for the lifetime of the process, only transport is reconfigured;
	// keeps restrictions in place for everyone that got the client before the setup
	restrictedClient = &http.Client{
		Transport: restrictedRoundTripper{},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	restrictedTransport atomic.Pointer[http.Transport]

	// cloud metadata endpoints outside of link-local ranges
	metadataNetworks = mustParseCIDRs(
		"100.100.100.200/32",
		"fd00:ec2::/32",
	)

	// non-public networks not covered by net.IP helpers
	privateNetworks = mustParseCIDRs(
		"0.0.0.0/8",
		"100.64.0.0/10",
		"192.0.0.0/24",
		"198.18.0.0/15",
		"240.0.0.0/4",
	)

	// IPv6 prefixes with embedded IPv4 address
	nat64Network = mustParseCIDRs("64:ff9b::/96")[0]
	sixToFour    = mustParseCIDRs("2002::/16")[0]
)

// ParseNetworkRestriction converts string into network restriction
//
// Falls back to link-local restriction on empty and invalid values
func ParseNetworkRestriction(s string) (NetworkRestriction, error) {
	switch r := NetworkRestriction(strings.ToLower(strings.TrimSpace(s))); r {
	case "":
		return RestrictLinkLocal, nil
	case RestrictNone, RestrictLinkLocal, RestrictPrivate:
		return r, nil
	default:
		return RestrictLinkLocal, fmt.Errorf("unknown network restriction %q, expecting one of: %s, %s, %s", s, RestrictNone, RestrictLinkLocal, RestrictPrivate)
	}
}

func init() {
	restrictedTransport.Store(newRestrictedTransport(0, false, RestrictLinkLocal))
}

// RestrictedClient returns HTTP client that should be used
// for all requests to destinations users have control over
func RestrictedClient() *http.Client {
	return restrictedClient
}

// SetupRestricted reconfigures restricted HTTP client, expected to be called when booting
func SetupRestricted(timeout time.Duration, tlsInsecure bool, restriction NetworkRestriction) {
	restrictedTransport.Store(newRestrictedTransport(timeout, tlsInsecure, restriction))
	restrictedClient.Timeout = timeout
}

func (restrictedRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	return restrictedTransport.Load().RoundTrip(req)
}

func newRestrictedTransport(timeout time.Duration, tlsInsecure bool, restriction NetworkRestriction) *http.Transport {
	var (
		transport = http.DefaultTransport.(*http.Transport).Clone()
		dialer    = &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,

			// Address is checked right before connecting and after the host is resolved,
			// this covers hostnames pointing to restricted addresses as well
			Control: func(network, address string, _ syscall.RawConn) error {
				host, _, err := net.SplitHostPort(address)
				if err != nil {
					return err
				}

				if i := strings.IndexByte(host, '%'); i >= 0 {
					// strip IPv6 zone
					host = host[:i]
				}

				if ip := net.ParseIP(host); ip == nil || blockedIP(restriction, ip) {
					return fmt.Errorf("%w: %s", ErrRestrictedNetwork, host)
				}

				return nil
			},
		}
	)

	if timeout > 0 {
		dialer.Timeout = timeout
		transport.TLSHandshakeTimeout = timeout
	}

	if tlsInsecure {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	transport.DialContext = dialer.DialContext

	return transport
}

func blockedIP(restriction NetworkRestriction, ip net.IP) bool {
	if restriction == RestrictNone {
		return false
	}

	if ip4 := ip.To4(); ip4 != nil {
		ip = ip4
	} else if embedded := embeddedIPv4(ip); embedded != nil {
		return blockedIP(restriction, embedded)
	}

	if ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || containsIP(metadataNetworks, ip) {
		return true
	}

	if restriction != RestrictPrivate {
		return false
	}

	return ip.IsLoopback() ||
		ip.IsPrivate() ||
		ip.IsUnspecified() ||
		ip.IsMulticast() ||
		containsIP(privateNetworks, ip)
}

// embeddedIPv4 extracts IPv4 address from NAT64 and 6to4 addresses
func embeddedIPv4(ip net.IP) net.IP {
	ip = ip.To16()
	if ip == nil {
		return nil
	}

	switch {
	case nat64Network.Contains(ip):
		return net.IP(ip[12:16])
	case sixToFour.Contains(ip):
		return net.IP(ip[2:6])
	}

	return nil
}

func containsIP(nn []*net.IPNet, ip net.IP) bool {
	for _, n := range nn {
		if n.Contains(ip) {
			return true
		}
	}

	return false
}

func mustParseCIDRs(cidrs ...string) (out []*net.IPNet) {
	for _, cidr := range cidrs {
		_, n, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(err)
		}

		out = append(out, n)
	}

	return
}

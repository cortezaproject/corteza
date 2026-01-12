package dal

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// DSN represents a parsed REST API Data Source Name
// Format: rest://[user:password@]host[:port][/path][?param=value&...]
type DSN struct {
	Scheme   string
	Host     string
	Port     string
	Path     string
	Username string
	Password string

	// Authentication
	AuthType     string // bearer, basic, apikey, oauth2, none
	Token        string // For bearer tokens
	APIKey       string // For API key auth
	APIKeyHeader string // Header name for API key (default: X-API-Key)
	ClientID     string // For OAuth2
	ClientSecret string // For OAuth2
	TokenURL     string // For OAuth2

	// Connection settings
	Timeout             time.Duration
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeout     time.Duration

	// Headers
	Headers map[string][]string

	// TLS/Security
	TLSInsecure bool

	// Retry settings
	MaxRetries     int
	RetryWaitMin   time.Duration
	RetryWaitMax   time.Duration
	RetryHTTPCodes []int

	// Rate limiting
	RateLimit       int // requests per second
	RateLimitBurst  int
	RateLimitPeriod time.Duration

	// Raw query parameters
	QueryParams url.Values

	Arbitrary map[string]any
}

// ParseDSN parses a REST API DSN string
func ParseDSN(dsn string) (out DSN, err error) {
	if dsn == "" {
		return out, fmt.Errorf("empty DSN")
	}

	u, err := url.Parse(dsn)
	if err != nil {
		return out, fmt.Errorf("invalid DSN: %w", err)
	}

	d := DSN{
		Scheme:      u.Scheme,
		Host:        u.Hostname(),
		Port:        u.Port(),
		Path:        u.Path,
		QueryParams: u.Query(),
		Headers:     make(map[string][]string),
	}

	// Extract basic auth from URL
	if u.User != nil {
		d.Username = u.User.Username()
		d.Password, _ = u.User.Password()
	}

	// Parse query parameters
	if err := d.parseQueryParams(); err != nil {
		return out, err
	}

	// Set defaults
	d.setDefaults()

	return d, nil
}

// parseQueryParams extracts configuration from query parameters
func (d *DSN) parseQueryParams() error {
	q := d.QueryParams

	// Authentication
	d.AuthType = q.Get("auth")
	d.Token = q.Get("token")
	d.APIKey = q.Get("apikey")
	d.APIKeyHeader = q.Get("apikey_header")
	d.ClientID = q.Get("client_id")
	d.ClientSecret = q.Get("client_secret")
	d.TokenURL = q.Get("token_url")

	// Connection timeouts
	if timeout := q.Get("timeout"); timeout != "" {
		dur, err := time.ParseDuration(timeout)
		if err != nil {
			return fmt.Errorf("invalid timeout: %w", err)
		}
		d.Timeout = dur
	}

	if idleTimeout := q.Get("idle_timeout"); idleTimeout != "" {
		dur, err := time.ParseDuration(idleTimeout)
		if err != nil {
			return fmt.Errorf("invalid idle_timeout: %w", err)
		}
		d.IdleConnTimeout = dur
	}

	// Connection pool settings
	if maxIdle := q.Get("max_idle_conns"); maxIdle != "" {
		n, err := strconv.Atoi(maxIdle)
		if err != nil {
			return fmt.Errorf("invalid max_idle_conns: %w", err)
		}
		d.MaxIdleConns = n
	}

	if maxIdlePerHost := q.Get("max_idle_conns_per_host"); maxIdlePerHost != "" {
		n, err := strconv.Atoi(maxIdlePerHost)
		if err != nil {
			return fmt.Errorf("invalid max_idle_conns_per_host: %w", err)
		}
		d.MaxIdleConnsPerHost = n
	}

	// TLS settings
	if insecure := q.Get("tls_insecure"); insecure != "" {
		d.TLSInsecure = insecure == "true" || insecure == "1"
	}

	// Retry settings
	if maxRetries := q.Get("max_retries"); maxRetries != "" {
		n, err := strconv.Atoi(maxRetries)
		if err != nil {
			return fmt.Errorf("invalid max_retries: %w", err)
		}
		d.MaxRetries = n
	}

	if retryWaitMin := q.Get("retry_wait_min"); retryWaitMin != "" {
		dur, err := time.ParseDuration(retryWaitMin)
		if err != nil {
			return fmt.Errorf("invalid retry_wait_min: %w", err)
		}
		d.RetryWaitMin = dur
	}

	if retryWaitMax := q.Get("retry_wait_max"); retryWaitMax != "" {
		dur, err := time.ParseDuration(retryWaitMax)
		if err != nil {
			return fmt.Errorf("invalid retry_wait_max: %w", err)
		}
		d.RetryWaitMax = dur
	}

	if retryCodes := q.Get("retry_codes"); retryCodes != "" {
		codes := strings.Split(retryCodes, ",")
		d.RetryHTTPCodes = make([]int, 0, len(codes))
		for _, code := range codes {
			n, err := strconv.Atoi(strings.TrimSpace(code))
			if err != nil {
				return fmt.Errorf("invalid retry code: %w", err)
			}
			d.RetryHTTPCodes = append(d.RetryHTTPCodes, n)
		}
	}

	// Rate limiting
	if rateLimit := q.Get("rate_limit"); rateLimit != "" {
		n, err := strconv.Atoi(rateLimit)
		if err != nil {
			return fmt.Errorf("invalid rate_limit: %w", err)
		}
		d.RateLimit = n
	}

	if rateLimitBurst := q.Get("rate_limit_burst"); rateLimitBurst != "" {
		n, err := strconv.Atoi(rateLimitBurst)
		if err != nil {
			return fmt.Errorf("invalid rate_limit_burst: %w", err)
		}
		d.RateLimitBurst = n
	}

	if rateLimitPeriod := q.Get("rate_limit_period"); rateLimitPeriod != "" {
		dur, err := time.ParseDuration(rateLimitPeriod)
		if err != nil {
			return fmt.Errorf("invalid rate_limit_period: %w", err)
		}
		d.RateLimitPeriod = dur
	}

	// Custom headers (format: header.Name=Value)
	for key, values := range q {
		if strings.HasPrefix(key, "header.") {
			headerName := strings.TrimPrefix(key, "header.")
			if len(values) > 0 {
				d.Headers[headerName] = values
			}
		}
	}

	// Arbitrary JSON payload
	if arbitrary := q.Get("arbitrary"); arbitrary != "" {
		var m map[string]any
		if err := json.Unmarshal([]byte(arbitrary), &m); err != nil {
			return fmt.Errorf("invalid arbitrary parameter (must be valid JSON): %w", err)
		}
		d.Arbitrary = m
	}

	return nil
}

// setDefaults sets default values for unspecified fields
func (d *DSN) setDefaults() {
	if d.AuthType == "" {
		if d.Token != "" {
			d.AuthType = "bearer"
		} else if d.APIKey != "" {
			d.AuthType = "apikey"
		} else if d.Username != "" && d.Password != "" {
			d.AuthType = "basic"
		} else {
			d.AuthType = "none"
		}
	}

	if d.APIKeyHeader == "" {
		d.APIKeyHeader = "X-API-Key"
	}

	if d.Timeout == 0 {
		d.Timeout = 30 * time.Second
	}

	if d.IdleConnTimeout == 0 {
		d.IdleConnTimeout = 90 * time.Second
	}

	if d.MaxIdleConns == 0 {
		d.MaxIdleConns = 100
	}

	if d.MaxIdleConnsPerHost == 0 {
		d.MaxIdleConnsPerHost = 10
	}

	if d.Port == "" {
		if d.Scheme == "rests" || d.Scheme == "https" {
			d.Port = "443"
		} else {
			d.Port = "80"
		}
	}

	if d.MaxRetries == 0 {
		d.MaxRetries = 3
	}

	if d.RetryWaitMin == 0 {
		d.RetryWaitMin = 1 * time.Second
	}

	if d.RetryWaitMax == 0 {
		d.RetryWaitMax = 30 * time.Second
	}

	if len(d.RetryHTTPCodes) == 0 {
		d.RetryHTTPCodes = []int{429, 500, 502, 503, 504}
	}

	if d.RateLimitPeriod == 0 {
		d.RateLimitPeriod = 1 * time.Second
	}
}

// BaseURL returns the base URL for the API
func (d *DSN) BaseURL() string {
	scheme := d.Scheme
	if scheme == "rest" {
		scheme = "http"
	} else if scheme == "rests" {
		scheme = "https"
	}

	portPart := ""
	if d.Port != "" && d.Port != "80" && d.Port != "443" {
		portPart = ":" + d.Port
	}

	return fmt.Sprintf("%s://%s%s%s", scheme, d.Host, portPart, d.Path)
}

// basicAuth encodes username:password to base64
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64Encode(auth)
}

func base64Encode(s string) string {
	// Simple base64 encoding - in real implementation use encoding/base64
	return fmt.Sprintf("base64(%s)", s) // Placeholder
}

// String returns a string representation of the DSN (without sensitive data)
func (d *DSN) String() string {
	scheme := d.Scheme
	auth := ""

	if d.Username != "" {
		auth = d.Username + ":***@"
	}

	port := ""
	if d.Port != "" {
		port = ":" + d.Port
	}

	return fmt.Sprintf("%s://%s%s%s%s", scheme, auth, d.Host, port, d.Path)
}

func (d *DSN) ToDSN() string {
	var sb strings.Builder

	// Scheme
	sb.WriteString(d.Scheme)
	sb.WriteString("://")

	// Basic auth in URL
	if d.Username != "" {
		sb.WriteString(url.QueryEscape(d.Username))
		if d.Password != "" {
			sb.WriteString(":")
			sb.WriteString(url.QueryEscape(d.Password))
		}
		sb.WriteString("@")
	}

	// Host and port
	sb.WriteString(d.Host)
	if d.Port != "" {
		sb.WriteString(":")
		sb.WriteString(d.Port)
	}

	// Path
	if d.Path != "" {
		if !strings.HasPrefix(d.Path, "/") {
			sb.WriteString("/")
		}
		sb.WriteString(d.Path)
	}

	// Query parameters
	params := url.Values{}

	// Authentication parameters
	if d.AuthType != "" && d.AuthType != "none" {
		params.Set("auth", d.AuthType)
	}
	if d.Token != "" {
		params.Set("token", d.Token)
	}
	if d.APIKey != "" {
		params.Set("apikey", d.APIKey)
	}
	if d.APIKeyHeader != "" && d.APIKeyHeader != "X-API-Key" {
		params.Set("apikey_header", d.APIKeyHeader)
	}
	if d.ClientID != "" {
		params.Set("client_id", d.ClientID)
	}
	if d.ClientSecret != "" {
		params.Set("client_secret", d.ClientSecret)
	}
	if d.TokenURL != "" {
		params.Set("token_url", d.TokenURL)
	}

	// Connection settings
	if d.Timeout != 0 && d.Timeout != 30*time.Second {
		params.Set("timeout", d.Timeout.String())
	}
	if d.IdleConnTimeout != 0 && d.IdleConnTimeout != 90*time.Second {
		params.Set("idle_timeout", d.IdleConnTimeout.String())
	}
	if d.MaxIdleConns != 0 && d.MaxIdleConns != 100 {
		params.Set("max_idle_conns", strconv.Itoa(d.MaxIdleConns))
	}
	if d.MaxIdleConnsPerHost != 0 && d.MaxIdleConnsPerHost != 10 {
		params.Set("max_idle_conns_per_host", strconv.Itoa(d.MaxIdleConnsPerHost))
	}

	// TLS settings
	if d.TLSInsecure {
		params.Set("tls_insecure", "true")
	}

	// Retry settings
	if d.MaxRetries != 0 && d.MaxRetries != 3 {
		params.Set("max_retries", strconv.Itoa(d.MaxRetries))
	}
	if d.RetryWaitMin != 0 && d.RetryWaitMin != 1*time.Second {
		params.Set("retry_wait_min", d.RetryWaitMin.String())
	}
	if d.RetryWaitMax != 0 && d.RetryWaitMax != 30*time.Second {
		params.Set("retry_wait_max", d.RetryWaitMax.String())
	}
	if len(d.RetryHTTPCodes) > 0 {
		// Only add if different from default [429, 500, 502, 503, 504]
		defaultCodes := []int{429, 500, 502, 503, 504}
		if !intSlicesEqual(d.RetryHTTPCodes, defaultCodes) {
			codes := make([]string, len(d.RetryHTTPCodes))
			for i, code := range d.RetryHTTPCodes {
				codes[i] = strconv.Itoa(code)
			}
			params.Set("retry_codes", strings.Join(codes, ","))
		}
	}

	// Rate limiting
	if d.RateLimit != 0 {
		params.Set("rate_limit", strconv.Itoa(d.RateLimit))
	}
	if d.RateLimitBurst != 0 {
		params.Set("rate_limit_burst", strconv.Itoa(d.RateLimitBurst))
	}
	if d.RateLimitPeriod != 0 && d.RateLimitPeriod != 1*time.Second {
		params.Set("rate_limit_period", d.RateLimitPeriod.String())
	}

	// Custom headers
	for key, value := range d.Headers {
		params.Set("header."+key, strings.Join(value, ","))
	}

	// Custom stuff
	if d.Arbitrary != nil {
		bb, err := json.Marshal(d.Arbitrary)
		if err != nil {
			panic(err)
		}
		params.Set("arbitrary", string(bb))
	}

	// Append query string
	if len(params) > 0 {
		sb.WriteString("?")
		sb.WriteString(params.Encode())
	}

	return sb.String()
}

// intSlicesEqual compares two int slices for equality
func intSlicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

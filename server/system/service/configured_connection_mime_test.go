package service

import (
	"encoding/base64"
	"strings"
	"testing"
)

func decode(t *testing.T, encoded string) string {
	t.Helper()
	b, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("failed to base64url-decode output: %v", err)
	}
	return string(b)
}

func TestBuildMIMEEmail_RequiredHeaders(t *testing.T) {
	raw, err := buildMIMEEmail("", "to@example.com", "Hello", "Body text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msg := decode(t, raw)
	if !strings.Contains(msg, "To: to@example.com\r\n") {
		t.Errorf("To header missing or malformed:\n%s", msg)
	}
	if !strings.Contains(msg, "Subject: Hello\r\n") {
		t.Errorf("Subject header missing or malformed:\n%s", msg)
	}
	if !strings.Contains(msg, "Content-Type: text/plain; charset=UTF-8\r\n") {
		t.Errorf("Content-Type header missing or malformed:\n%s", msg)
	}
	if !strings.Contains(msg, "Body text") {
		t.Errorf("body content missing:\n%s", msg)
	}
}

func TestBuildMIMEEmail_WithFrom(t *testing.T) {
	raw, err := buildMIMEEmail("sender@example.com", "to@example.com", "Subject", "Body")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msg := decode(t, raw)
	if !strings.Contains(msg, "From: sender@example.com\r\n") {
		t.Errorf("From header missing or malformed:\n%s", msg)
	}
}

func TestBuildMIMEEmail_EmptyFromOmitted(t *testing.T) {
	raw, err := buildMIMEEmail("", "to@example.com", "Subject", "Body")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msg := decode(t, raw)
	if strings.Contains(msg, "From:") {
		t.Errorf("From header should be omitted when empty, got:\n%s", msg)
	}
}

func TestBuildMIMEEmail_ValidBase64URL(t *testing.T) {
	raw, err := buildMIMEEmail("", "to@example.com", "Subject", "Body")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Must not contain standard base64 characters that are replaced in URL encoding
	if strings.ContainsAny(raw, "+/") {
		t.Errorf("output contains non-URL-safe base64 characters: %s", raw)
	}
}

func TestBuildMIMEEmail_HeaderBodySeparator(t *testing.T) {
	raw, err := buildMIMEEmail("", "to@example.com", "Subj", "Hello World")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	msg := decode(t, raw)
	// RFC 2822: blank line separates headers from body
	if !strings.Contains(msg, "\r\n\r\n") {
		t.Errorf("missing blank line (CRLF CRLF) between headers and body:\n%s", msg)
	}
	parts := strings.SplitN(msg, "\r\n\r\n", 2)
	if len(parts) != 2 {
		t.Fatalf("unexpected split result: %d parts", len(parts))
	}
	if parts[1] != "Hello World" {
		t.Errorf("body mismatch: got %q", parts[1])
	}
}

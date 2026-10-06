package xss

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	nethtml "golang.org/x/net/html"
)

// A value is sanitised on every write and every read, so each check runs it
// through RichText repeatedly, the way a stored value is.
func passes(in string, n int) []string {
	out := make([]string, 0, n)
	for i := 0; i < n; i++ {
		in = RichText(in)
		out = append(out, in)
	}
	return out
}

// tags lists the start tags in s, with the names of their attributes.
func tags(s string) (tt []string) {
	z := nethtml.NewTokenizer(strings.NewReader(s))
	for {
		switch z.Next() {
		case nethtml.ErrorToken:
			return
		case nethtml.StartTagToken, nethtml.SelfClosingTagToken:
			name, more := z.TagName()
			tag := string(name)
			for more {
				var key []byte
				key, _, more = z.TagAttr()
				tag += " " + string(key)
			}
			tt = append(tt, tag)
		}
	}
}

func TestRichTextNeverTurnsTextIntoMarkup(t *testing.T) {
	payload := `<img src=x onerror=alert(1)>`
	encoded := payload
	for depth := 0; depth < 6; depth++ {
		encoded = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(encoded)
		for i, out := range passes(encoded, 8) {
			require.Empty(t, tags(out), "encoded %d times, pass %d: %q", depth+1, i+1, out)
		}
	}
}

func TestRichTextKeepsAttributesClosed(t *testing.T) {
	for i, out := range passes(`<a href="https://example.com/&#34; onmouseover=alert(1)">a</a>`, 4) {
		for _, tag := range tags(out) {
			require.NotContains(t, tag, "onmouseover", "pass %d: %q", i+1, out)
		}
	}
}

func TestRichTextKeepsText(t *testing.T) {
	for _, tc := range []struct{ name, in, out string }{
		{"email the editor escaped", `<p>John &lt;john@example.com&gt;</p>`, `<p>John &lt;john@example.com></p>`},
		{"ampersand and quotes", `a & b, it's "x"`, `a & b, it's "x"`},
		{"arrow", `before->after, a > b`, `before->after, a > b`},
		{"less-than", `a < b`, `a &lt; b`},
		{"entities in text", `a &amp; b, it&#39;s &#34;x&#34;`, `a & b, it's "x"`},
		{"allowed markup", `<p><strong>bold</strong> <a href="mailto:j@example.com">j</a></p>`, `<p><strong>bold</strong> <a href="mailto:j@example.com" rel="nofollow">j</a></p>`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for i, out := range passes(tc.in, 3) {
				require.Equal(t, tc.out, out, "pass %d", i+1)
			}
		})
	}
}

package locale

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/crusttech/human/server/pkg/options"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

// underscores bypass the separator guard in x/text and hit its quadratic path
var oversizedTag = "en" + strings.Repeat("_aaaaa", 160000)

func TestParseTagOversized(t *testing.T) {
	start := time.Now()

	_, err := ParseTag(oversizedTag)
	require.ErrorIs(t, err, ErrTagTooLong)
	require.Equal(t, language.Und, MakeTag(oversizedTag))

	_, _, err = ParseAcceptLanguage(oversizedTag)
	require.ErrorIs(t, err, ErrTagTooLong)

	require.Less(t, time.Since(start), 50*time.Millisecond)

	tag, err := ParseTag("de-AT")
	require.NoError(t, err)
	require.Equal(t, language.MustParse("de-AT"), tag)

	tags, _, err := ParseAcceptLanguage("en-US,en;q=0.9,de;q=0.8")
	require.NoError(t, err)
	require.Len(t, tags, 3)
}

func TestResolveAcceptLanguageOversized(t *testing.T) {
	ll := &service{
		opt:  options.LocaleOpt{QueryStringParam: "lng"},
		tags: []language.Tag{language.English, language.German},
		def:  &Language{Tag: language.English},
	}

	for _, via := range []string{"header", "query"} {
		r := httptest.NewRequest("GET", "/", nil)
		if via == "header" {
			r.Header.Set(AcceptLanguageHeader, oversizedTag)
		} else {
			r.URL.RawQuery = "lng=" + oversizedTag
		}

		start := time.Now()
		require.Equal(t, language.English, resolveAcceptLanguageHeaders(r, ll), via)
		require.Less(t, time.Since(start), 50*time.Millisecond, via)
	}

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set(AcceptLanguageHeader, "de-AT,de;q=0.9")
	require.Equal(t, language.German, resolveAcceptLanguageHeaders(r, ll))
}

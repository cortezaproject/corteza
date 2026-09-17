package helpers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRecursiveDotEnvLoad_localePath(t *testing.T) {
	// <tmp>/server/.env beside <tmp>/locale, loaded from <tmp>/server/tests/suite
	// the way an integration suite loads it.
	setup := func(t *testing.T, env string) (locale string) {
		old, had := os.LookupEnv("LOCALE_PATH")
		os.Unsetenv("LOCALE_PATH")
		t.Cleanup(func() {
			if had {
				os.Setenv("LOCALE_PATH", old)
			} else {
				os.Unsetenv("LOCALE_PATH")
			}
		})

		tmp, err := filepath.EvalSymlinks(t.TempDir())
		require.NoError(t, err)

		locale = filepath.Join(tmp, "locale")
		suite := filepath.Join(tmp, "server", "tests", "suite")
		require.NoError(t, os.MkdirAll(locale, 0o755))
		require.NoError(t, os.MkdirAll(suite, 0o755))
		require.NoError(t, os.WriteFile(filepath.Join(tmp, "server", ".env"), []byte(env), 0o644))
		t.Chdir(suite)

		return locale
	}

	t.Run("relative entries resolve against the .env directory", func(t *testing.T) {
		locale := setup(t, "LOCALE_PATH=../locale:/opt/locale\n")

		RecursiveDotEnvLoad()

		first, rest, _ := strings.Cut(os.Getenv("LOCALE_PATH"), ":")
		first, err := filepath.EvalSymlinks(first)
		require.NoError(t, err)
		require.Equal(t, locale, first)
		require.Equal(t, "/opt/locale", rest)
	})

	t.Run("a LOCALE_PATH set before loading is left alone", func(t *testing.T) {
		setup(t, "LOCALE_PATH=../locale\n")
		os.Setenv("LOCALE_PATH", "../mine")

		RecursiveDotEnvLoad()

		require.Equal(t, "../mine", os.Getenv("LOCALE_PATH"))
	})
}

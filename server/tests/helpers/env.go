package helpers

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

// RecursiveDotEnvLoad loads the nearest .env from the working directory or up
// to two directories above it. A relative LOCALE_PATH read from that file is
// resolved against the file's directory, where the server runs.
func RecursiveDotEnvLoad() {
	for _, loc := range []string{".env", "../.env", "../../.env"} {
		if _, err := os.Stat(loc); err == nil {
			_, preset := os.LookupEnv("LOCALE_PATH")
			godotenv.Load(loc)
			if !preset {
				resolveLocalePath(filepath.Dir(loc))
			}
			break
		}
	}
}

// resolveLocalePath makes every relative entry of the colon-separated
// LOCALE_PATH absolute, relative to dir.
func resolveLocalePath(dir string) {
	val := os.Getenv("LOCALE_PATH")
	if val == "" {
		return
	}

	pp := strings.Split(val, ":")
	for i, p := range pp {
		if p == "" || filepath.IsAbs(p) {
			continue
		}

		if abs, err := filepath.Abs(filepath.Join(dir, p)); err == nil {
			pp[i] = abs
		}
	}

	os.Setenv("LOCALE_PATH", strings.Join(pp, ":"))
}

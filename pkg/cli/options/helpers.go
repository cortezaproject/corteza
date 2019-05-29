package options

import (
	"os"
	"strings"
	"time"

	"github.com/spf13/cast"
)

func makeEnvKeys(pfix, name string) []string {
	return []string{
		strings.ToUpper(strings.Trim(pfix, "_") + "_" + name),
		strings.ToUpper(name),
	}
}

func EnvString(pfix, key string, def string) string {
	for _, key = range makeEnvKeys(pfix, key) {
		if val, has := os.LookupEnv(key); has {
			return val
		}
	}
	return def
}

func EnvBool(pfix, key string, def bool) bool {
	for _, key = range makeEnvKeys(pfix, key) {
		if val, has := os.LookupEnv(key); has {
			if b, err := cast.ToBoolE(val); err == nil {
				return b
			}
		}
	}
	return def
}

func EnvInt(pfix, key string, def int) int {
	for _, key = range makeEnvKeys(pfix, key) {
		if val, has := os.LookupEnv(key); has {
			if i, err := cast.ToIntE(val); err == nil {
				return i
			}
		}
	}
	return def
}

func EnvDuration(pfix, key string, def time.Duration) time.Duration {
	for _, key = range makeEnvKeys(pfix, key) {
		if val, has := os.LookupEnv(key); has {
			if d, err := cast.ToDurationE(val); err == nil {
				return d
			}
		}
	}
	return def
}

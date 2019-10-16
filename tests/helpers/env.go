package helpers

import (
	"os"

	"github.com/joho/godotenv"
)

// RecursiveDotEnvLoad loads ENV variables from .evn files 3 levels down
func RecursiveDotEnvLoad() {
	for _, loc := range []string{".env", "../.env", "../../.env"} {
		if _, err := os.Stat(loc); err == nil {
			godotenv.Load(loc)
		}
	}
}

package helpers

import (
	"github.com/joho/godotenv"
	"os"
)

// RecursiveDotEnvLoad loads ENV variables from .evn files 3 levels down
func RecursiveDotEnvLoad() {
	for _, loc := range []string{".env", "../.env", "../../.env"} {
		if _, err := os.Stat(loc); err == nil {
			godotenv.Load(loc)
			break
		}
	}
}

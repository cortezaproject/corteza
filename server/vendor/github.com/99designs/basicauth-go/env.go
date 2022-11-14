package basicauth

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// NewFromEnv reads a set of credentials in from environment variables in
// the format {PREFIX}{USERNAME|tolower}=password1,password2 and returns
// middleware that will validate incoming requests.
func NewFromEnv(realm, prefix string) func(http.Handler) http.Handler {
	credentials := map[string][]string{}

	re := regexp.MustCompile(fmt.Sprintf("^%s(?P<username>.*)$", strings.ToUpper(prefix)))
	for _, envVar := range os.Environ() {
		name, value := split2(envVar, "=")

		if res := re.FindStringSubmatch(name); res != nil {
			username := strings.ToLower(res[1])
			credentials[username] = strings.Split(value, ",")
		}
	}

	return New(realm, credentials)
}

func split2(s, sep string) (string, string) {
	res := strings.SplitN(s, sep, 2)

	return res[0], res[1]
}

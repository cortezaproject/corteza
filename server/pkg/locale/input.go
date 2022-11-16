package locale

import (
	"github.com/cortezaproject/corteza/server/pkg/xss"
)

func SanitizeMessage(in string) string {
	return xss.RichText(in)
}

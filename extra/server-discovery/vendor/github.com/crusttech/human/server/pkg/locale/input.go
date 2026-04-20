package locale

import (
	"github.com/crusttech/human/server/pkg/xss"
)

func SanitizeMessage(in string) string {
	return xss.RichText(in)
}

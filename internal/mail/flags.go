package mail

import (
	"github.com/crusttech/crust/internal/config"
)

var flags *config.SMTP

func Flags(prefix ...string) {
	flags = new(config.SMTP).Init(prefix...)
}

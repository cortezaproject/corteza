package mail

import (
	"fmt"

	"github.com/crusttech/crust/internal/config"
)

var flags *config.SMTP

func Flags(prefix ...string) {
	flags = new(config.SMTP).Init(prefix...)
}

func Debug() {
	fmt.Printf("Debug SMTP flags: %#v", *flags)
}

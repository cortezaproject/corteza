package rbac

import (
	"fmt"

	"github.com/crusttech/crust/internal/config"
)

var flags *config.RBAC

func Flags(prefix ...string) {
	flags = new(config.RBAC).Init(prefix...)
}

func Debug() {
	fmt.Printf("Debug RBAC flags: %#v", *flags)
}

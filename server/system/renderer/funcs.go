package renderer

import "github.com/crusttech/human/server/pkg/valuestore"

func envGetter() func(k string) any {
	return func(k string) any {
		return valuestore.Global().Env(k)
	}
}

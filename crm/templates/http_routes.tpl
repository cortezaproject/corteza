package {package}

{load warning.tpl}

import (
	"fmt"
	"runtime"
	"reflect"
	"github.com/go-chi/chi"
)

func MountRoutes(r chi.Router) {
{foreach $apis as $api}
	{api.interface|strtolower} := {api.interface|capitalize}Handlers{}.new()
{/foreach}
{foreach $apis as $api}
	r.Route("{api.path}", func(r chi.Router) {
{foreach $api.apis as $call}
		r.{eval echo capitalize(strtolower($call.method))}("{call.path}", {api.interface|strtolower}.{call.name|capitalize})
{/foreach}
	})
{/foreach}

	var printRoutes func(chi.Routes, string, string)
	printRoutes = func(r chi.Routes, indent string, prefix string) {
		routes := r.Routes()
		for _, route := range routes {
			if route.SubRoutes != nil && len(route.SubRoutes.Routes()) > 0 {
				fmt.Printf(indent+"%s - with %d handlers, %d subroutes\n", route.Pattern, len(route.Handlers), len(route.SubRoutes.Routes()))
				printRoutes(route.SubRoutes, indent+"\t", prefix+route.Pattern[:len(route.Pattern)-2])
			} else {
				for key, fn := range route.Handlers {
					fmt.Printf("%s%s\t%s -> %s\n", indent, key, prefix+route.Pattern, runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name())
				}
			}
		}
	}
	printRoutes(r, "", "")
}
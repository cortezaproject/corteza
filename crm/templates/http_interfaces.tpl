package {package}

{load warning.tpl}

import (
	"net/http"
)

// HTTP handlers are a superset of internal APIs
type {name}Handlers struct {
	*{name}
}

func ({name}Handlers) new() *{name}Handlers {
	return &{name}Handlers{
		{name}{}.New(),
	}
}

// Internal API interface
type {name}API interface {
{foreach $calls as $call}
	{call.name|capitalize}(*{name|lcfirst}{call.name|capitalize}Request) (interface{}, error)
{/foreach}
}

// HTTP API interface
type {name}HandlersAPI interface {
{foreach $calls as $call}
	{call.name|capitalize}(http.ResponseWriter, *http.Request)
{/foreach}

	// Authenticate API requests
	Authenticator() func(http.Handler) http.Handler
}

// Compile time check to see if we implement the interfaces
var _ {name}HandlersAPI = &{name}Handlers{}
var _ {name}API = &{name}{}
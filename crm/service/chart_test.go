package service

import (
	"testing"
	"context"

	"github.com/crusttech/crust/crm/rest/request"
)

func TestModuleGraph(t *testing.T) {
	repository := Module().With(context.Background())

	params := &request.ModuleChart{
		Kind: "line",
	}

	{
		_, err := repository.Chart(params)
		assert(t, err == nil, "Error when getting module graph: %+v", err)
	}

	{
		params.Kind = "404"
		_, err := repository.Chart(params)
		assert(t, err != nil, "Expected error when getting chart for 404")
	}
}

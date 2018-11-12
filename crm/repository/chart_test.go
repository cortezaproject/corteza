package repository

import (
	"context"
	"testing"

	"github.com/crusttech/crust/crm/rest/request"
)

func TestModuleGraph(t *testing.T) {
	repository := Module(context.TODO(), nil).With(context.Background(), nil)

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

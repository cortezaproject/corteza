package helpers

import (
	"context"
	"github.com/cortezaproject/corteza-server/pkg/label/types"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func SetLabelsViaAPI(api *apitest.APITest, t *testing.T, endpoint string, in, out interface{}) {
	var (
		payload = struct {
			Response interface{}
		}{}
		req *apitest.Request
	)

	payload.Response = out

	if strings.HasPrefix(endpoint, "PUT ") {
		// a little workaround for our inconsistencies...
		req = api.Put(endpoint[4:])
	} else {
		req = api.Post(endpoint)
	}

	req.JSON(JSON(in)).
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(AssertNoErrors).
		End().
		JSON(&payload)
}

func SearchWithLabelsViaAPI(api *apitest.APITest, t *testing.T, endpoint string, res interface{}, labels url.Values) {
	var payload = struct{ Response struct{ Set interface{} } }{}

	payload.Response.Set = res

	api.Get(endpoint).
		QueryCollection(labels).
		Expect(t).
		Status(http.StatusOK).
		Assert(AssertNoErrors).
		End().
		JSON(&payload)
}

func LoadLabelsFromStore(t *testing.T, s store.Storer, resKind string, id uint64) (sl map[string]string) {
	ll, _, err := s.SearchLabels(context.Background(), types.LabelFilter{Kind: resKind})
	require.NoError(t, err)
	if len(ll) == 0 {
		// small adjustment
		return nil
	}

	return ll.FilterByResource(resKind, id)
}

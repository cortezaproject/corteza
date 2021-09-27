package apigw

import (
	"net/http"
	"testing"
)

func Test_prefilter_query_failing(t *testing.T) {
	var (
		_, h, _ = setupScenario(t)
	)

	h.apiInit().
		Get("/test").
		Query("token", "brute-force-guess").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusBadRequest).
		End()
}

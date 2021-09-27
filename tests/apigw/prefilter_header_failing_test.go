package apigw

import (
	"net/http"
	"testing"
)

func Test_prefilter_header_failing(t *testing.T) {
	var (
		_, h, _ = setupScenario(t)
	)

	h.apiInit().
		Get("/test").
		Header("Accept", "application/json").
		Header("Token", "brute-force-guess").
		Expect(t).
		Status(http.StatusBadRequest).
		End()
}

package apigw

import (
	"net/http"
	"testing"
)

func Test_prefilter_header_passing(t *testing.T) {
	var (
		_, h, _ = setupScenario(t)
	)

	h.apiInit().
		Get("/test").
		Header("Accept", "application/json").
		Header("Token", "super-secret-token").
		Expect(t).
		Status(http.StatusOK).
		End()

	h.apiInit().
		Get("/test-hyphen").
		Header("Accept-Language", "fr-CH, fr;q=0.9").
		Header("Token", "super-secret-token").
		Expect(t).
		Status(http.StatusOK).
		End()
}

package apigw

import (
	"net/http"
	"testing"
)

func Test_prefilter_header_n_routes(t *testing.T) {
	var (
		_, h, _ = setupScenario(t)
	)

	// First (a) route query validation
	h.apiInit().
		Get("/a").
		Header("P", "a").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		End()

	h.apiInit().
		Get("/a").
		Header("P", "b").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Second (b) route query validation
	h.apiInit().
		Get("/b").
		Header("P", "b").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		End()

	h.apiInit().
		Get("/b").
		Header("P", "a").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusBadRequest).
		End()
}

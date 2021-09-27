package apigw

import (
	"net/http"
	"testing"
)

func Test_prefilter_query_n_routes(t *testing.T) {
	var (
		_, h, _ = setupScenario(t)
	)

	// First (a) route query validation
	h.apiInit().
		Get("/a").
		Query("p", "a").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		End()

	h.apiInit().
		Get("/a").
		Query("p", "b").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusBadRequest).
		End()

	// Second (b) route query validation
	h.apiInit().
		Get("/b").
		Query("p", "b").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		End()

	h.apiInit().
		Get("/b").
		Query("p", "a").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusBadRequest).
		End()
}

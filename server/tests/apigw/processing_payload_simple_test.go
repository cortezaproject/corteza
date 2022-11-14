package apigw

import (
	"net/http"
	"testing"
)

func Test_processor_payload_simple(t *testing.T) {
	var (
		_, h, _ = setupScenario(t)
	)

	h.apiInit().
		Get("/test").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Body("60").
		End()

}

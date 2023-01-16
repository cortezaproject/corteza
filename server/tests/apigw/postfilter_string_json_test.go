package apigw

import (
	"net/http"
	"testing"
)

func Test_postfilter_string_json(t *testing.T) {
	var (
		_, h, _ = setupScenario(t)
	)

	h.apiInit().
		Get("/json/string").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Body("{\"baz\":{\"@value\":{\"1\":{\"@value\":1,\"@type\":\"Float\"},\"2\":{\"@value\":2,\"@type\":\"Float\"}},\"@type\":\"Vars\"},\"foo\":{\"@value\":\"bar\",\"@type\":\"String\"}}").
		End()
}

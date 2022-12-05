package apigw

import (
	"net/http"
	"testing"
)

func Test_postfilter_string_kv(t *testing.T) {
	var (
		_, h, _ = setupScenario(t)
	)

	h.apiInit().
		Get("/json/kv").
		Header("Accept", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Body("{\"baz\":\"123\",\"foo\":\"bar\"}").
		End()
}

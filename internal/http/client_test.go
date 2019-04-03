// +build external

package http

import (
	"testing"

	"github.com/crusttech/crust/internal/config"
	"github.com/crusttech/crust/internal/test"
)

func TestHTTPClient(t *testing.T) {
	client, err := New(&config.HTTPClient{
		Timeout: 10,
	})
	test.Assert(t, err == nil, "%+v", err)
	client.Debug(FULL)

	req, err := client.Get("https://api.scene-si.org/fortune.php")
	test.Assert(t, err == nil, "%+v", err)

	resp, err := client.Do(req)
	test.Assert(t, err == nil, "%+v", err)

	err = func() error {
		defer resp.Body.Close()
		switch resp.StatusCode {
		case 200:
			return nil
		default:
			return ToError(resp)
		}
	}()

	test.Assert(t, err == nil, "Invalid response: %+v", err)
}

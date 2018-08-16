package rbac

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

func toError(resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if body == nil || err != nil {
		return errors.Errorf("unexpected response (%d, %s)", resp.StatusCode, err)
	}
	return errors.New(string(body))
}

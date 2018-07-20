package rbac

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

func toError(resp *http.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Errorf("unexpected response (%d)", resp.StatusCode)
	}
	return errors.New(string(body))
}

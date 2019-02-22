package permit

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type (
	httpClient interface {
		Do(req *http.Request) (*http.Response, error)
	}
)

const (
	permitCheckEndpoint = "https://permit.crust.tech/check"
)

func Check(ctx context.Context, p Permit) (*Permit, error) {
	return CheckWithClient(ctx, http.DefaultClient, p)
}

func CheckWithClient(ctx context.Context, client httpClient, p Permit) (*Permit, error) {
	if len(p.Key) == 0 {
		return nil, errors.New("key not set")
	} else if len(p.Key) != KeyLength {
		return nil, errors.Errorf("invalid key length (%d chars)", len(p.Key))
	}

	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(p); err != nil {
		return nil, errors.Wrap(err, "permit encoding failed")
	}

	if req, err := http.NewRequest("POST", permitCheckEndpoint, buf); err != nil {
		return nil, errors.Wrap(err, "unable to create request")
	} else {
		return CheckWithRequest(client, req.WithContext(ctx))
	}
}

func CheckWithRequest(client httpClient, request *http.Request) (p *Permit, err error) {
	var rsp *http.Response
	if rsp, err = client.Do(request); err != nil {
		return nil, errors.Wrap(err, "unable to fetch permit")
	}

	defer rsp.Body.Close()

	switch rsp.StatusCode {
	case http.StatusBadRequest:
		return nil, errors.New("bad request")
	case http.StatusNotFound:
		return nil, errors.New("subscription key not found")
	case http.StatusInternalServerError:
		return nil, errors.New("subscription server error")
	case http.StatusUnauthorized:
		return nil, errors.New("subscription key invalid")
	}

	p = &Permit{}
	if err = json.NewDecoder(rsp.Body).Decode(&p); err != nil {
		return nil, errors.Wrap(err, "unable to decode response into permit")
	}

	return
}

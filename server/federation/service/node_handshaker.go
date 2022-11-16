package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cortezaproject/corteza/server/federation/types"
)

type (
	// HttpHandshake performs HTTP handshake with remote node
	httpNodeHandshake struct {
		httpClient httpNodeHandshaker
	}

	httpNodeHandshaker interface {
		Do(*http.Request) (*http.Response, error)
	}
)

func HttpHandshake(client httpNodeHandshaker) *httpNodeHandshake {
	return &httpNodeHandshake{httpClient: client}
}

// Init calls handshake init on a remote node via HTTP
//
// Called from node.Pair() and sends POST request to remote Corteza
// There, it is handled with node.HandshakeInit
//
// Fn is not part of the node struct to allow injection
func (h httpNodeHandshake) Init(ctx context.Context, n *types.Node, authToken string) error {
	var (
		endpoint = fmt.Sprintf("%s/nodes/%d/handshake", n.BaseURL, n.SharedNodeID)
		payload  = map[string]string{
			// Share pairing token, so we can authenticate ourselves
			"pairToken": n.PairToken,

			// Share auth token for the federation service user for this node
			"authToken":    authToken,
			"sharedNodeID": strconv.FormatUint(n.SharedNodeID, 10),
		}
	)

	return h.send(ctx, endpoint, payload)
}

// Confirm calls handshake init on a remote node via HTTP
//
// Called from node.Confirm and sends POST request to remote Corteza
// There, it's handled with node.HandshakeComplete
func (h httpNodeHandshake) Complete(ctx context.Context, n *types.Node, authToken string) error {
	var (
		endpoint = fmt.Sprintf("%s/nodes/%d/handshake-complete", n.BaseURL, n.SharedNodeID)
		payload  = map[string]string{
			// Share auth token for the federation service user for this node
			"authToken": authToken,
		}
	)

	// use n.AuthToken in context
	ctx = context.WithValue(ctx, "authToken", n.AuthToken)

	return h.send(ctx, endpoint, payload)
}

func (h httpNodeHandshake) send(ctx context.Context, endpoint string, payload map[string]string) error {
	body := bytes.NewBuffer(nil)
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	if authToken := ctx.Value("authToken"); authToken != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	}

	rsp, err := h.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}

	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("pairing failed, remote host retuned non-200 response")
	}

	return nil
}

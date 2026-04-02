package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/system/service"
)

type RagClient struct {
	discoveryBaseURL string
	httpClient       *http.Client
}

func NewRagClient(discoveryBaseURL string) *RagClient {
	return &RagClient{
		discoveryBaseURL: strings.TrimRight(discoveryBaseURL, "/"),
		httpClient:       &http.Client{Timeout: 30 * time.Second},
	}
}

type RagRequest struct {
	Query string `json:"query"`
}

type RagResponse struct {
	Answer string `json:"answer"`
	// add other fields as needed
}

func (c *RagClient) Query(ctx context.Context, userID uint64, roles []uint64, query string) (*RagResponse, error) {
	// Mint a short-lived service token on behalf of the requesting user
	// auth.TokenIssuer is already initialized by the server's initAuth()
	//

	user, err := service.DefaultUser.FindByAny(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("could not find the user: %w", err)
	}

	token, err := auth.TokenIssuer.Issue(ctx,
		auth.WithIdentity(user),
		auth.WithExpiration(5*time.Minute),
		auth.WithScope("discovery"),
	)
	if err != nil {
		return nil, fmt.Errorf("could not issue discovery token: %w", err)
	}

	body, _ := json.Marshal(RagRequest{Query: query})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.discoveryBaseURL+"/rag",
		strings.NewReader(string(body)),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+string(token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("discovery request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("discovery returned %d", resp.StatusCode)
	}

	var result RagResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

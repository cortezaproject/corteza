package service

import (
	"fmt"
	"github.com/cortezaproject/corteza/server/federation/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNode_decodePairingURI(t *testing.T) {
	var (
		tcc = []struct {
			name string
			uri  string
			node *types.Node
			err  error
		}{
			{
				"happy path",
				"https://42:secret-secret-secret-secret-1234@example.tld/federation?name=Noddy",
				&types.Node{SharedNodeID: 42, Name: "Noddy", BaseURL: "https://example.tld/federation", PairToken: "secret-secret-secret-secret-1234"},
				nil,
			},
			{
				"no name",
				"https://42:secret-secret-secret-secret-1234@example.tld",
				&types.Node{SharedNodeID: 42, BaseURL: "https://example.tld/", PairToken: "secret-secret-secret-secret-1234"},
				nil,
			},
			{
				"no token",
				"https://42@example.tld",
				nil,
				NodeErrPairingURITokenInvalid(),
			},
			{
				"no node id",
				"https://:secret-secret-secret-secret-1234@example.tld",
				nil,
				NodeErrPairingURISourceIDInvalid(),
			},
			{
				"invalid URL",
				"https://this is not a valid url",
				nil,
				NodeErrPairingURIInvalid().Wrap(fmt.Errorf(`invalid character " " in host name`)),
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req      = require.New(t)
				out, err = (&node{}).decodePairingURI(tc.uri)
			)

			if tc.err != nil {
				req.Equal(tc.err.Error(), err.Error())
			} else {
				req.NoError(err)
				req.Equal(tc.node, out)
			}
		})
	}
}

func TestNode_makePairingURI(t *testing.T) {
	var (
		svc = &node{host: "example.tld", baseURL: "federation", name: "Noddy"}
		req = require.New(t)
		n   = &types.Node{ID: 42, PairToken: "secret"}
		uri = svc.makePairingURI(n)
	)

	req.Equal("https://42:secret@example.tld/federation?name=Noddy", uri)
}

package util

import (
	"fmt"
	"net/url"
	"strconv"
)

var (
	ErrorURIMissingToken = fmt.Errorf("uri: token missing")
)

const ()

type (
	NodeURIParams struct {
		Name string
	}

	DecodedURI struct {
		Domain string
		Token  string
		NodeID uint64
		Params NodeURIParams
	}
)

func DecodeURI(i string) (du *DecodedURI, err error) {
	du = &DecodedURI{
		Params: NodeURIParams{},
	}

	u, err := url.Parse(i)
	if err != nil {
		return
	}

	du.Domain = u.Host
	nodeS := u.User.Username()
	if nodeS == "" {
		return nil, ErrorURIMissingToken
	}
	nodeID, err := strconv.ParseUint(nodeS, 10, 64)
	if err != nil {
		return nil, err
	}
	du.NodeID = nodeID

	token, has := u.User.Password()
	if !has {
		return nil, ErrorURIMissingToken
	}
	du.Token = token

	params := u.Query()
	if name, has := params["name"]; has && len(name) > 0 {
		du.Params.Name = name[0]
	}

	return
}

func EncodeURI(ott, domain string, nodeID uint64) string {
	u := url.URL{
		Scheme: "corteza",
		User:   url.UserPassword(strconv.FormatUint(nodeID, 10), ott),
		Host:   domain,
	}

	return u.String()
}

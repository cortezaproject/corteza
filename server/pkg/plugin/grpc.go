package plugin

import (
	hcp "github.com/hashicorp/go-plugin"
)

var (
	handshakeConfig = hcp.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "BASIC_PLUGIN",
		MagicCookieValue: "hello",
	}

	allowedProtocols = []hcp.Protocol{hcp.ProtocolGRPC}
)

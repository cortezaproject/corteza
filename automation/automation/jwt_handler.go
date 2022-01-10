package automation

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type (
	jwtHandler struct {
		reg jwtHandlerRegistry
	}
)

func JwtHandler(reg jwtHandlerRegistry) *jwtHandler {
	h := &jwtHandler{
		reg: reg,
	}

	h.register()
	return h
}

func (h jwtHandler) generate(ctx context.Context, args *jwtGenerateArgs) (res *jwtGenerateResults, err error) {
	var (
		secret interface{}

		auxp = make(map[string]interface{})
		auxh = make(map[string]interface{})
	)

	if !args.hasPayload {
		err = fmt.Errorf("could not generate JWT, payload missing")
		return
	}

	if !args.hasSecret {
		err = fmt.Errorf("could not generate JWT, secret or cert missing")
		return
	}

	for k, v := range args.headerVars {
		auxh[k] = v.Get()
	}

	for k, v := range args.payloadVars {
		auxp[k] = v.Get()
	}

	if args.payloadString != "" {
		if err = json.Unmarshal([]byte(args.payloadString), &auxp); err != nil {
			return
		}
	}

	if args.headerString != "" {
		if err = json.Unmarshal([]byte(args.headerString), &auxh); err != nil {
			return
		}
	}

	// check for delimiters
	auxp["scope"] = strings.FieldsFunc(args.Scope, func(r rune) bool {
		return r == ' ' || r == ','
	})

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(auxp))

	// merge header with user-provided header
	for k, v := range auxh {
		token.Header[k] = v
	}

	// check if we use cert
	{
		pemBlock, _ := pem.Decode([]byte(args.secretString))

		if pemBlock != nil {
			if secret, err = x509.ParsePKCS8PrivateKey(pemBlock.Bytes); err != nil {
				return
			}
		} else {
			secret = []byte(args.secretString)
		}
	}

	res = &jwtGenerateResults{}
	res.Token, err = token.SignedString(secret)

	return
}

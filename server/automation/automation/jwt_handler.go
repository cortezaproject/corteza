package automation

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jws"
	"github.com/lestrrat-go/jwx/jwt"
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

	var (
		key interface{}

		tkn        = jwt.New()
		tokenBytes []byte
		headers    = jws.NewHeaders()
	)

	for k, v := range auxh {
		_ = headers.Set(k, v)
	}

	for k, v := range auxp {
		if err = tkn.Set(k, v); err != nil {
			return
		}
	}

	if decodedKey, _ := pem.Decode([]byte(args.secretString)); decodedKey != nil {
		if key, err = x509.ParsePKCS8PrivateKey(decodedKey.Bytes); err != nil {
			return nil, err
		}

		tokenBytes, err = jwt.Sign(tkn, jwa.RS256, key, jwt.WithHeaders(headers))
	} else {
		key = []byte(args.secretString)
		tokenBytes, err = jwt.Sign(tkn, jwa.HS256, key, jwt.WithHeaders(headers))
	}

	if err != nil {
		return nil, fmt.Errorf("could not sign token: %w", err)
	}

	return &jwtGenerateResults{Token: string(tokenBytes)}, nil
}

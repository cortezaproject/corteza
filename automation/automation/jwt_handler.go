package automation

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
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
		tkn        = jwt.New()
		keySet     jwk.Set
		tokenBytes []byte
	)

	for k, v := range auxp {
		if err = tkn.Set(k, v); err != nil {
			return
		}
	}

	//< HEAD
	//	// check if we use cert
	//	{
	//		pemBlock, _ := pem.Decode([]byte(args.secretString))
	//
	//		if pemBlock != nil {
	//			if secret, err = x509.ParsePKCS8PrivateKey(pemBlock.Bytes); err != nil {
	//				return
	//			}
	//		} else {
	//			secret = []byte(args.secretString)
	//		}
	//=
	// @todo check if jwk.Parse provides the same logic as before with pem.Decode and x59.ParsePkC8PrivateKey
	if keySet, err = jwk.Parse([]byte(args.secretString)); err != nil {
		return
		//> e3a304d5... Replacing dgrijalva/jwt-go with lestrrat-go/jwx
	}

	if tokenBytes, err = jwt.Sign(tkn, jwa.HS512, keySet); err != nil {
		return
	}

	return &jwtGenerateResults{Token: string(tokenBytes)}, nil
}

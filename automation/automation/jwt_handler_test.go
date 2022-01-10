package automation

import (
	"context"
	"errors"
	"testing"

	"github.com/cortezaproject/corteza-server/pkg/expr"
	"github.com/stretchr/testify/require"
)

func TestJwtHandler(t *testing.T) {
	type (
		tf struct {
			name   string
			exp    string
			err    error
			params *jwtGenerateArgs
		}
	)

	var (
		handler = &jwtHandler{}
		tcc     = []tf{
			{
				name: "proxy processer with auth headers",
				exp:  "eyJhbGciOiJSUzI1NiIsImtleSI6InZhbHVlIiwidHlwIjoiSldUIn0.eyJrZXkiOiJ2YWx1ZSIsInNjb3BlIjpbInNjb3BlIiwic2NvcGUyIl19.OnnLBsJdrZwAPeMjYIJPlZPW7bZdp-JbgNKFrRPnPXcuQVAUxK-R69-kDZbdsRZFOaU-AC52Tz4Ft3-SzADgFroSpoNjpEdJwbaKANlG_pm2b8-1pXSc0YzhY7DqBs4iae2pI0FGpeT_dza6kp9TL3NQfgqjx05Q3Gz5-3Kk32MD3zIvpUXgwkPbb4XLJxiY2Ra1dRWVI5Guk4GyLA19b7Z-DrHg1GE9mDy_NwZZD994Iri9e5zmcAikRtfHO7guPtBKVwhvt3u37wXtRgEYMyFQAn2ZSZaTytK8161Y-TOcdLVlqy4OfasaVt1pP0aNI9GGz5R-OVCOghW7TqZ6YQ",
				err:  nil,
				params: &jwtGenerateArgs{
					Scope:        `scope scope2`,
					headerVars:   must(expr.CastToVars(map[string]interface{}{"key": "value"})),
					payloadVars:  must(expr.CastToVars(map[string]interface{}{"key": "value"})),
					secretString: prKey,
					hasHeader:    true,
					hasPayload:   true,
					hasSecret:    true},
			},
			{
				name: "proxy processer with auth headers",
				exp:  "",
				err:  errors.New("could not generate JWT, payload missing"),
				params: &jwtGenerateArgs{
					hasHeader:  true,
					hasPayload: false,
					hasSecret:  true},
			},
			{
				name: "proxy processer with auth headers",
				exp:  "",
				err:  errors.New("could not generate JWT, secret or cert missing"),
				params: &jwtGenerateArgs{
					hasHeader:  true,
					hasPayload: true,
					hasSecret:  false},
			},
		}
	)

	for _, tc := range tcc {
		t.Run(tc.name, func(t *testing.T) {
			var (
				req = require.New(t)
				ctx = context.Background()
			)

			res, err := handler.generate(ctx, tc.params)

			if tc.err == nil {
				req.NoError(err)
			} else {
				req.EqualError(err, tc.err.Error())
			}

			if res != nil {
				req.Equal(tc.exp, res.Token)
			}
		})
	}
}

const prKey string = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCn2MSVSyMzdZ5Q
cD0KMZAuhAmWE06qS/agpkG9OmvxsBAXjS8/2oaLCtmpF7Sx7XEJOaVUz8caty0L
Qjpovv5ZNqY0S/A+KahayTVXvNhjtwVb4kSYv8Tfv1H72CJAsHEi++DuFJ3uJyW5
zAblXpY96dCLKd4R6BOwQSAqpqAgCBJX3NiM1CtB6zvdYcGPxUJmFOxUSdU779WY
IloTDCi9lY/SWgt5tHb38jiCv8uZ9B+sgD6PXXqjnOnji+XXJb3zGrmyeSBvhgF+
bqlMcOXaufikzdq8v6wrW6zPsdDCYqtbpyku0pVz5JuN/czItS4fFGaQ+hjYnnq4
nU5YKGhjAgMBAAECggEAPlzd/ZJjS9VhswVgyI7NwVqxrR8TVVbQFbRwLHyuaqg9
8mI0sgbhgnvPj3INYyaTnxfaA/8HPTfd9pbu2MhN/Ju/eSLV6mLT+JdVyHmT9Mil
pxQU5KQr4+5T6bzOTTbBcnwfgJYMb9X/wF68GTDhpbNgFrTBm+mclxo7d11dlUiQ
bdA2KBhIM54kKcWEO90AP2FiC/cRSbv8afHVb0Fxy1SBAGVDSscUaEwgyHkRf6H7
Mmt4xhfBwLZZk5pVa97pZsas7L0sm4BtyYSs4elJesuo3DTpezZSzWHNreJcCL00
v7pm9ZWv9PYeUcQqSw6Ws2fhsydiN93o8s/fPpF0IQKBgQDQqvKPjGo3HLRnBmpc
nhUVQgjdh/QobJCGBjdxR/C3NkcatNLONddhefqovfb2jITdHmRKW+7XgRq/q457
mZ2JPvvryiWcVqf6D91tSoFYrVg5nU1zr6nTga1BpTH00/eNEXxjMAAguPj3m7Sw
3BP7L8yZZI1pOOjzVzuISEDlGwKBgQDN62ZXzRqorn5t+mDywTAJ2qAyZzoDAtUP
CjANOKBCECqE9cOW1qBDJcaEF7G+xyB22CMciOELjGXWw6ZwFWjc23tDCEkoHoia
5IWvwNs5ncdPCuOm3aay9w2dx9Idy+Eq5Vpe4XDlnZvEfzjwXX3gbQ5YS9kbNE2r
XtmgD/BmWQKBgBNKqrBA0BUWT0tzGWRErThQ6ZbpmdYe62Gos3mCqCuYFgzPCOpN
qgL2DwmIvotexG3ZAHardzJvWjS8PKkKs7jbnNjY0I9ap58D1nnjOIAlTpHNDDsU
04OdapI2Hp8+9ZUSN8jHyEs+Lq5ds9/iCOrhKW5JEJXY0BinSPa5j15fAoGBAM2I
dqCQollXwe34CaiD13UeeOOWUTsMKqlWW9v2d085X5dSzyTRmSksnVbfZ5SqoOa+
mV0z6pxiSIvywUACvqYjlIa10H9w6pzgF+fzMV3y9CsbDVtSxb7ABSFFf54qD9eH
EYq+rrchd4bMDYMtbiUB9V2AZ3VV4Wh5xfKTtjoRAoGAec/nBmX/oHiG/3T2LEmo
bTT+/eXc4LUzgHvFG6HW40PI8T2TfMTUGQH+90zkraryUF5PNoUF58tooDAcMWXm
DPEHFHeWn8T4obfuUjTw1mwj7Xjzr40HDitjIQa0bvBVwdYyEgiEw1CaeaUpq9Z8
ElGWoRHT87pfoZjbdPz7a4Y=
-----END PRIVATE KEY-----`

func must(v map[string]expr.TypedValue, err error) map[string]expr.TypedValue {
	if err != nil {
		panic(err)
	}
	return v
}

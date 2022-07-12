package app

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/options"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPrepareSignatureFnParams(t *testing.T) {
	type (
		test struct {
			name string
			opt  options.AuthOpt
			err  error
		}
	)

	var (
		privateKey = string(genKey())
	)

	tests := []test{
		{
			name: "empty algo",
			opt:  options.AuthOpt{JwtKey: "foobar"},
			err:  fmt.Errorf("token signature algorithm empty or missing"),
		},
		{
			name: "unknown algo",
			opt:  options.AuthOpt{JwtAlgorithm: "foobar", JwtKey: "foobar"},
			err:  fmt.Errorf("token signature algorithm \"foobar\" not supported"),
		},
		{
			name: "empty key",
			opt:  options.AuthOpt{JwtAlgorithm: jwa.HS256.String()},
			err:  fmt.Errorf("token secret missing"),
		},
		{
			name: "empty key",
			opt:  options.AuthOpt{JwtAlgorithm: jwa.PS256.String()},
			err:  fmt.Errorf("token key missing"),
		},
		{
			// "shared secret" string
			name: "HS256",
			opt:  options.AuthOpt{JwtAlgorithm: jwa.HS256.String(), Secret: "test key"},
		},
		{
			// "shared secret" string
			name: "HS384",
			opt:  options.AuthOpt{JwtAlgorithm: jwa.HS384.String(), Secret: "test key"},
		},
		{
			// "shared secret" string
			name: "HS512",
			opt:  options.AuthOpt{JwtAlgorithm: jwa.HS512.String(), Secret: "test key"},
		},
		{
			// requires private key
			name: "PS256",
			opt:  options.AuthOpt{JwtAlgorithm: jwa.PS256.String(), JwtKey: privateKey},
		},
		{
			// requires private key
			name: "PS384",
			opt:  options.AuthOpt{JwtAlgorithm: jwa.PS384.String(), JwtKey: privateKey},
		},
		{
			// requires private key
			name: "PS512",
			opt:  options.AuthOpt{JwtAlgorithm: jwa.PS512.String(), JwtKey: privateKey},
		},
		{
			// requires private key
			name: "RS256",
			opt:  options.AuthOpt{JwtAlgorithm: jwa.RS256.String(), JwtKey: privateKey},
		},
		{
			// requires private key
			name: "RS384",
			opt:  options.AuthOpt{JwtAlgorithm: jwa.RS384.String(), JwtKey: privateKey},
		},
		{
			// requires private key
			name: "RS512",
			opt:  options.AuthOpt{JwtAlgorithm: jwa.RS512.String(), JwtKey: privateKey},
		},
		{
			// requires private key
			name: "RS512 from a file",
			opt:  options.AuthOpt{JwtAlgorithm: jwa.RS512.String(), JwtKey: "test_files/key.pem"},
		},
		{
			// requires private key
			name: "RS512 from a non-existing file",
			opt:  options.AuthOpt{JwtAlgorithm: jwa.RS512.String(), JwtKey: "test_files/not-here.pem"},
			err:  fmt.Errorf("could not read key file: open test_files/not-here.pem: no such file or directory"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				req = require.New(t)
			)

			alg, key, err := prepareSignatureFnParams(tt.opt)
			if tt.err != nil {
				req.EqualError(err, tt.err.Error())
				return
			}

			req.NoError(err)
			req.NotEmpty(alg)
			req.NotEmpty(key)
		})
	}
}

func genKey() []byte {
	bitSize := 4096

	// Generate RSA key.
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		panic(err)
	}

	// Encode private key to PKCS#1 ASN.1 PEM.
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

}

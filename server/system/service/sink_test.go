package service

import (
	"bytes"
	"github.com/cortezaproject/corteza/server/pkg/errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	internalAuth "github.com/cortezaproject/corteza/server/pkg/auth"
)

func Test_sink_SignURL(t *testing.T) {
	var (
		signer = internalAuth.HmacSigner("test")

		tests = []struct {
			name          string
			surp          SinkRequestUrlParams
			wantSignedURL string
			wantErr       bool
		}{
			{
				name: "basic",
				surp: SinkRequestUrlParams{
					Method:      "POST",
					Origin:      "test",
					Expires:     nil,
					MaxBodySize: 1024,
					ContentType: "plain/text",
				},
				wantSignedURL: "/sink?__sign=d8a8c5591acb0f5f6695ab6aa4a205a7066b3bf4_eyJtdGQiOiJQT1NUIiwib3JpZ2luIjoidGVzdCIsIm1icyI6MTAyNCwiY3QiOiJwbGFpbi90ZXh0In0%3D",
			},
			{
				name: "signature in a path",
				surp: SinkRequestUrlParams{
					Method:          "POST",
					Origin:          "test",
					Expires:         nil,
					MaxBodySize:     1024,
					ContentType:     "plain/text",
					SignatureInPath: true,
				},
				wantSignedURL: "/sink/__sign=a4a5652c66159ed0142c01a9aa5d90b3e9f76241_eyJtdGQiOiJQT1NUIiwib3JpZ2luIjoidGVzdCIsIm1icyI6MTAyNCwiY3QiOiJwbGFpbi90ZXh0Iiwic2lwIjp0cnVlfQ==",
			},
			{
				name: "required path",
				surp: SinkRequestUrlParams{
					Path: "/foo/bar",
				},
				wantSignedURL: "/sink/foo/bar?__sign=910436a3944ad1ee010db07b936209f198be481d_eyJwdCI6Ii9mb28vYmFyIn0%3D",
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := sink{
				signer: signer,
			}

			gotSignedURL, _, err := svc.SignURL(tt.surp)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignURL() \n"+
					"  error: %v, \n"+
					"wantErr: %v", err, tt.wantErr)
				return
			}
			if gotSignedURL.String() != tt.wantSignedURL {
				t.Errorf("SignURL() \n"+
					"gotSignedURL: %v\n"+
					"        want: %v", gotSignedURL, tt.wantSignedURL)
			}
		})
	}
}

func Test_sink_handleRequest(t *testing.T) {
	var (
		signer = internalAuth.HmacSigner("test")
		svc    = sink{signer: signer}

		signParams = SinkRequestUrlParams{
			Method:      "POST",
			Origin:      "test",
			Expires:     nil,
			MaxBodySize: 1024,
			ContentType: "plain/text",
		}
		signedUrl, _, _ = svc.SignURL(signParams)

		signParamsNoPref      = SinkRequestUrlParams{}
		signedUrlNoPref, _, _ = svc.SignURL(signParamsNoPref)

		signParamsExp      = SinkRequestUrlParams{Expires: &time.Time{}}
		signedUrlExp, _, _ = svc.SignURL(signParamsExp)

		signParamsInPath      = SinkRequestUrlParams{SignatureInPath: true, Path: "/foo"}
		signedUrlInPath, _, _ = svc.SignURL(signParamsInPath)

		signParamsInPathNoPath      = SinkRequestUrlParams{SignatureInPath: true}
		signedUrlInPathNoPath, _, _ = svc.SignURL(signParamsInPathNoPath)

		signParamsFixedPath      = SinkRequestUrlParams{Path: "/foo/bar"}
		signedUrlFixedPath, _, _ = svc.SignURL(signParamsFixedPath)
	)

	var (
		tests = []struct {
			name       string
			withMethod string
			withURL    string
			withBody   io.Reader
			withHeader http.Header
			wantParams *SinkRequestUrlParams
			wantErr    error
		}{
			{
				name:    "missing signature",
				withURL: "/sink",
				wantErr: SinkErrMissingSignature(),
			},
			{
				name:    "invalid signature",
				withURL: "/sink?" + SinkSignUrlParamName + "=foo",
				wantErr: SinkErrInvalidSignatureParam(),
			},
			{
				name:    "invalid signature",
				withURL: "/sink?__sign=foo_bar",
				wantErr: SinkErrBadSinkParamEncoding(),
			},
			{
				name:    "invalid signature",
				withURL: "/sink?__sign=foo_eyJtdGQiOiJQT1NUIiwib3JpZ2luIjoidGVzdCIsIm1icyI6MTAyNCwiY3QiOiJwbGFpbi90ZXh0In0%3D",
				wantErr: SinkErrInvalidSignature(),
			},
			{
				name:       "any HTTP method (no pref. method)",
				withMethod: "DELETE",
				withURL:    signedUrlNoPref.String(),
				wantParams: &signParamsNoPref,
				wantErr:    nil,
			},
			{
				name:       "invalid HTTP method (POST only)",
				withMethod: "GET",
				withURL:    signedUrl.String(),
				wantErr:    SinkErrInvalidHttpMethod(),
			},
			{
				name:       "invalid content type",
				withMethod: "POST",
				withHeader: map[string][]string{"content-type": {"foo/bar"}},
				withURL:    signedUrl.String(),
				wantErr:    SinkErrInvalidContentType(),
			},
			{
				name:       "valid content type",
				withMethod: "POST",
				withHeader: map[string][]string{"content-type": {"plain/text"}},
				withURL:    signedUrl.String(),
				wantErr:    nil,
				wantParams: &signParams,
			},
			{
				name:       "expired",
				withMethod: "POST",
				withURL:    signedUrlExp.String(),
				wantErr:    SinkErrSignatureExpired(),
			},
			{
				name:       "content length exceeds",
				withMethod: "POST",
				withHeader: map[string][]string{"content-type": {"plain/text"}},
				withBody:   bytes.NewBufferString(strings.Repeat(".", 1025)),
				withURL:    signedUrl.String(),
				wantErr:    SinkErrContentLengthExceedsMaxAllowedSize(),
			},
			{
				name:       "signature in a path (no path constraint)",
				withURL:    signedUrlInPathNoPath.String(),
				wantParams: &signParamsInPathNoPath,
			},
			{
				name:       "signature in a path (with path constraint)",
				withURL:    signedUrlInPath.String(),
				wantParams: &signParamsInPath,
			},
			{
				name:       "signed fixed path",
				withURL:    signedUrlFixedPath.String(),
				wantParams: &signParamsFixedPath,
			},
		}
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.withMethod == "" {
				tt.withMethod = "GET"
			}

			req, err := http.NewRequest(tt.withMethod, tt.withURL, tt.withBody)
			if err != nil {
				panic(err)
			}

			for n, vv := range tt.withHeader {
				for _, v := range vv {
					req.Header.Add(n, v)
				}
			}

			got, err := svc.handleRequest(req)
			if tt.wantErr != nil && err != nil && !errors.Is(err, tt.wantErr) {
				t.Errorf("handleRequest()\n"+
					"  error: %v\n"+
					"wantErr: %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.wantParams) {
				t.Errorf("handleRequest()\n"+
					"params: %v\n"+
					"  want: %v", got, tt.wantParams)
			}
		})
	}
}

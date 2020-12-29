package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza-server/pkg/actionlog"
	"github.com/cortezaproject/corteza-server/pkg/api"
	internalAuth "github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/eventbus"
	"github.com/cortezaproject/corteza-server/system/service/event"
	"github.com/cortezaproject/corteza-server/system/types"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type (
	sink struct {
		signer     internalAuth.Signer
		actionlog  actionlog.Recorder
		eventbus   sinkEventDispatcher
		isMonolith bool
	}

	SinkRequestUrlParams struct {
		// Expect sink request to be of this method
		Method string `json:"mtd,omitempty"`

		// Origin is as an identifier, no validation of request params
		Origin string `json:"origin,omitempty"`

		// Optional, signature expiration
		Expires *time.Time `json:"exp,omitempty"`

		// When set it enables body processing (but limits it to that size!)
		MaxBodySize int64 `json:"mbs,omitempty"`

		// Acceptable content type
		ContentType string `json:"ct,omitempty"`

		Path string `json:"pt,omitempty"`

		// Should we put signature in the path (true)
		// or in query string (false, default)
		SignatureInPath bool `json:"sip,omitempty"`
	}

	sinkEventDispatcher interface {
		WaitFor(ctx context.Context, ev eventbus.Event) (err error)
	}
)

const (
	SinkContentTypeMail = "message/rfc822"

	// base url
	// we're using this for router, signature...
	SinkBaseURL = "/sink"

	// name of the parameter used for sink request signature
	SinkSignUrlParamName = "__sign"

	// delimiter between signature and payload
	SinkSignUrlParamDelimiter = "_"
)

func Sink() *sink {
	return &sink{
		actionlog:  DefaultActionlog,
		signer:     internalAuth.DefaultSigner,
		eventbus:   eventbus.Service(),
		isMonolith: true,
	}
}

// SignURL takes sink request parameters and generates signed URL
//
// With signed URL, external systems can make requests to sink subsystem
// and trigger scripts
func (svc sink) SignURL(srup SinkRequestUrlParams) (signedURL *url.URL, out SinkRequestUrlParams, err error) {
	var (
		params []byte
		sap    = &sinkActionProps{sinkParams: &srup}
		qs     = url.Values{}
	)

	err = func() error {
		// Append normalized path to the base URL
		srup.Path = svc.pathCleanup(srup.Path)
		path := svc.GetPath() + srup.Path

		srup.Method = strings.ToUpper(srup.Method)

		params, err = json.Marshal(srup)
		if err != nil {
			return SinkErrFailedToSign(sap).Wrap(err)
		}

		signature := svc.signer.Sign(0, params) + SinkSignUrlParamDelimiter + base64.StdEncoding.EncodeToString(params)

		if srup.SignatureInPath {
			// Optional, use path for sink signature
			path = fmt.Sprintf("%s/%s=%s", path, SinkSignUrlParamName, signature)
		} else {
			// By default put signature in query string
			qs.Set(SinkSignUrlParamName, signature)
		}

		signedURL = &url.URL{RawQuery: qs.Encode(), Path: path}
		return nil
	}()

	return signedURL, srup, svc.recordAction(context.Background(), sap, SinkActionSign, err)
}

func (svc sink) GetPath() string {
	path := ""

	if svc.isMonolith {
		path = "/system"
	}

	return path + SinkBaseURL
}

// pathCleanup removes base URL prefix and adds leading slash
func (svc sink) pathCleanup(p string) string {
	if len(p) > 0 {
		if pos := strings.Index(p, SinkBaseURL); pos > -1 {
			p = p[pos+len(SinkBaseURL):]
		}

		return "/" + strings.Trim(p, "/")
	}

	return ""
}

// ProcessRequest function is used directly in the HTTP controller
func (svc *sink) ProcessRequest(w http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		sap = &sinkActionProps{}
	)

	// capture error from request handling and process functions
	err := func() error {
		defer r.Body.Close()
		srup, err := svc.handleRequest(r)
		if err != nil {
			return err
		}

		var body io.Reader
		if srup.MaxBodySize > 0 {
			// Utilize body only when max-body-size limit is set
			body = http.MaxBytesReader(w, r.Body, srup.MaxBodySize)
		} else {
			body = http.MaxBytesReader(w, r.Body, 32<<10) // 32k limit
		}

		if err := svc.process(srup, w, r, body); err != nil {
			return SinkErrProcessingError(sap).Wrap(err)
		}

		return nil
	}()

	_ = svc.recordAction(ctx, sap, SinkActionRequest, err)
	if err != nil {
		// use standard facility for encoding errors for HTTP
		api.Send(w, r, err)
	}
}

// Verifies and extracts sink request params
func (svc sink) handleRequest(r *http.Request) (*SinkRequestUrlParams, error) {
	var (
		srup = &SinkRequestUrlParams{}
		sap  = &sinkActionProps{}
		qs   = r.URL.Query()

		signatureFoundInPath bool

		param string

		// this value is modified if signature is found in a path
		reqPath = r.URL.Path
	)

	// try to find a signature
	if _, has := qs[SinkSignUrlParamName]; has {
		// first, in a query string
		param = r.URL.Query().Get(SinkSignUrlParamName)
	} else if i := strings.Index(reqPath, SinkSignUrlParamName); i > -1 {
		// fallback to path, expecting signature to be at the end
		// offset string index by start of signature param name, length of param name, and = char
		param = reqPath[i+len(SinkSignUrlParamName)+1:]
		reqPath = reqPath[:i]

		// this is more for consistency and cleaner tests
		signatureFoundInPath = true
	}

	if len(param) == 0 {
		return nil, SinkErrMissingSignature(sap)
	}

	split := strings.SplitN(param, SinkSignUrlParamDelimiter, 2)
	if len(split) < 2 {
		return nil, SinkErrInvalidSignatureParam(sap)
	}

	params, err := base64.StdEncoding.DecodeString(split[1])
	if err != nil {
		return nil, SinkErrBadSinkParamEncoding(sap)
	}

	if !svc.signer.Verify(split[0], 0, params) {
		return nil, SinkErrInvalidSignature(sap)
	}

	if err = json.Unmarshal(params, srup); err != nil {
		// Impossible scenario :)
		// How can we have verified signature of an invalid JSON ?!
		return nil, SinkErrInvalidSinkRequestUrlParams(sap)
	}

	sap.setSinkParams(srup)

	if srup.SignatureInPath != signatureFoundInPath {
		return nil, SinkErrMisplacedSignature(sap)
	}

	if srup.Method != "" && srup.Method != r.Method {
		return nil, SinkErrInvalidHttpMethod(sap)
	}

	contentType := strings.ToLower(r.Header.Get("content-type"))
	if i := strings.Index(contentType, ";"); i > 0 {
		contentType = contentType[0 : i-1]
	}

	if srup.ContentType != "" {
		if strings.ToLower(srup.ContentType) != contentType {
			return nil, SinkErrInvalidContentType(sap)
		}
	}

	if srup.Path != "" {
		if srup.Path != svc.pathCleanup(reqPath) {
			return nil, SinkErrInvalidPath(sap)
		}
	}

	if srup.Expires != nil && srup.Expires.Before(time.Now()) {
		return nil, SinkErrSignatureExpired(sap)
	}

	if srup.MaxBodySize > 0 {
		// See if there is content length param and reject it right away
		if r.ContentLength > srup.MaxBodySize {
			return nil, SinkErrContentLengthExceedsMaxAllowedSize(sap)
		}
	}

	return srup, nil
}

// Processes sink request, casts it and forwards it to processor (depending on content type)
//
// Main reason for content-type & body to be passed separately (and not extracted from r param) is
// that:
// a) content type might be forced via sink params
// This is useful to enforce mail processing
// b) Max-body-size check might be limited via sink params
// and io.Reader that is passed is limited w/ io.LimitReader
//
func (svc *sink) process(srup *SinkRequestUrlParams, w http.ResponseWriter, r *http.Request, body io.Reader) error {
	var (
		err         error
		ctx         = r.Context()
		contentType = srup.ContentType
		sap         = &sinkActionProps{
			contentType: contentType,
		}
	)

	switch strings.ToLower(contentType) {
	case SinkContentTypeMail, "rfc822", "email", "mail":
		// this is handled by dedicated event that parses raw payload from HTTP request
		// as rfc882 message.
		var msg *types.MailMessage
		msg, err = types.NewMailMessage(body)
		if err != nil {
			return SinkErrFailedToCreateEvent(sap).Wrap(err)
		}

		sap.setMailHeader(&msg.Header)

		err = svc.eventbus.WaitFor(ctx, event.MailOnReceive(msg))
		if err != nil {
			return SinkErrFailedToProcess(sap).Wrap(err)
		}

	default:
		var (
			sr *types.SinkRequest

			// Predefine default response
			rsp = &types.SinkResponse{
				Status: http.StatusOK,
			}
		)

		// Sanitize URL
		sanitizedURL := r.URL

		// Step 1: removing sink sign url param
		sanitizedQuery := r.URL.Query()
		sanitizedQuery.Del(SinkSignUrlParamName)
		sanitizedURL.RawQuery = sanitizedQuery.Encode()

		// Step 2: remove prefix
		if i := strings.Index(sanitizedURL.Path, SinkBaseURL); i > -1 {
			sanitizedURL.Path = sanitizedURL.Path[i+len(SinkBaseURL):]
		}

		// Step 3: remove sink suffix if in path
		if srup.SignatureInPath {
			i := strings.Index(sanitizedURL.Path, SinkSignUrlParamName)
			if i > 0 {
				sanitizedURL.Path = sanitizedURL.Path[0 : i-1]
			}
		}

		r.URL = sanitizedURL
		r.RequestURI = sanitizedURL.String()

		sr, err = types.NewSinkRequest(r, body)
		if err != nil {
			return SinkErrFailedToCreateEvent(sap).Wrap(err)

		}

		sap.setUrl(sanitizedURL.String())

		err = svc.eventbus.WaitFor(ctx, event.SinkOnRequest(rsp, sr))
		if err != nil {
			return SinkErrFailedToProcess(sap).Wrap(err)
		}

		sap.setResponseStatus(rsp.Status)

		// Now write everything we've received from the script
		for k, vv := range rsp.Header {
			for _, v := range vv {
				w.Header().Add(k, v)
			}
		}

		w.WriteHeader(rsp.Status)

		var output []byte
		if bb, ok := rsp.Body.([]byte); ok {
			// Ok, handled
			output = bb
		} else if s, ok := rsp.Body.(string); ok {
			output = []byte(s)
		}

		if _, err = w.Write(output); err != nil {
			return SinkErrFailedToRespond(sap).Wrap(err)
		}
	}

	return nil
}

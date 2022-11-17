package request

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
//

import (
	"encoding/json"
	"fmt"
	"github.com/cortezaproject/corteza/server/pkg/payload"
	"github.com/go-chi/chi/v5"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
)

// dummy vars to prevent
// unused imports complain
var (
	_ = chi.URLParam
	_ = multipart.ErrMessageTooLarge
	_ = payload.ParseUint64s
	_ = strings.ToLower
	_ = io.EOF
	_ = fmt.Errorf
	_ = json.NewEncoder
)

type (
	// Internal API interface
	SmtpConfigurationCheckerCheck struct {
		// Host POST parameter
		//
		// SMTP server host name
		Host string

		// Port POST parameter
		//
		// SMTP server port
		Port uint

		// Recipients POST parameter
		//
		// List of recipients email addresses that should recieve test email
		Recipients []string

		// Username POST parameter
		//
		// SMTP server authentication username
		Username string

		// Password POST parameter
		//
		// SMTP server authentication password
		Password string

		// TlsInsecure POST parameter
		//
		// TLS mode
		TlsInsecure bool

		// TlsServerName POST parameter
		//
		// TLS server name
		TlsServerName string
	}
)

// NewSmtpConfigurationCheckerCheck request
func NewSmtpConfigurationCheckerCheck() *SmtpConfigurationCheckerCheck {
	return &SmtpConfigurationCheckerCheck{}
}

// Auditable returns all auditable/loggable parameters
func (r SmtpConfigurationCheckerCheck) Auditable() map[string]interface{} {
	return map[string]interface{}{
		"host":          r.Host,
		"port":          r.Port,
		"recipients":    r.Recipients,
		"username":      r.Username,
		"password":      r.Password,
		"tlsInsecure":   r.TlsInsecure,
		"tlsServerName": r.TlsServerName,
	}
}

// Auditable returns all auditable/loggable parameters
func (r SmtpConfigurationCheckerCheck) GetHost() string {
	return r.Host
}

// Auditable returns all auditable/loggable parameters
func (r SmtpConfigurationCheckerCheck) GetPort() uint {
	return r.Port
}

// Auditable returns all auditable/loggable parameters
func (r SmtpConfigurationCheckerCheck) GetRecipients() []string {
	return r.Recipients
}

// Auditable returns all auditable/loggable parameters
func (r SmtpConfigurationCheckerCheck) GetUsername() string {
	return r.Username
}

// Auditable returns all auditable/loggable parameters
func (r SmtpConfigurationCheckerCheck) GetPassword() string {
	return r.Password
}

// Auditable returns all auditable/loggable parameters
func (r SmtpConfigurationCheckerCheck) GetTlsInsecure() bool {
	return r.TlsInsecure
}

// Auditable returns all auditable/loggable parameters
func (r SmtpConfigurationCheckerCheck) GetTlsServerName() string {
	return r.TlsServerName
}

// Fill processes request and fills internal variables
func (r *SmtpConfigurationCheckerCheck) Fill(req *http.Request) (err error) {

	if strings.HasPrefix(strings.ToLower(req.Header.Get("content-type")), "application/json") {
		err = json.NewDecoder(req.Body).Decode(r)

		switch {
		case err == io.EOF:
			err = nil
		case err != nil:
			return fmt.Errorf("error parsing http request body: %w", err)
		}
	}

	{
		// Caching 32MB to memory, the rest to disk
		if err = req.ParseMultipartForm(32 << 20); err != nil && err != http.ErrNotMultipart {
			return err
		} else if err == nil {
			// Multipart params

			if val, ok := req.MultipartForm.Value["host"]; ok && len(val) > 0 {
				r.Host, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["port"]; ok && len(val) > 0 {
				r.Port, err = payload.ParseUint(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["username"]; ok && len(val) > 0 {
				r.Username, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["password"]; ok && len(val) > 0 {
				r.Password, err = val[0], nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["tlsInsecure"]; ok && len(val) > 0 {
				r.TlsInsecure, err = payload.ParseBool(val[0]), nil
				if err != nil {
					return err
				}
			}

			if val, ok := req.MultipartForm.Value["tlsServerName"]; ok && len(val) > 0 {
				r.TlsServerName, err = val[0], nil
				if err != nil {
					return err
				}
			}
		}
	}

	{
		if err = req.ParseForm(); err != nil {
			return err
		}

		// POST params

		if val, ok := req.Form["host"]; ok && len(val) > 0 {
			r.Host, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["port"]; ok && len(val) > 0 {
			r.Port, err = payload.ParseUint(val[0]), nil
			if err != nil {
				return err
			}
		}

		//if val, ok := req.Form["recipients[]"]; ok && len(val) > 0  {
		//    r.Recipients, err = val, nil
		//    if err != nil {
		//        return err
		//    }
		//}

		if val, ok := req.Form["username"]; ok && len(val) > 0 {
			r.Username, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["password"]; ok && len(val) > 0 {
			r.Password, err = val[0], nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["tlsInsecure"]; ok && len(val) > 0 {
			r.TlsInsecure, err = payload.ParseBool(val[0]), nil
			if err != nil {
				return err
			}
		}

		if val, ok := req.Form["tlsServerName"]; ok && len(val) > 0 {
			r.TlsServerName, err = val[0], nil
			if err != nil {
				return err
			}
		}
	}

	return err
}

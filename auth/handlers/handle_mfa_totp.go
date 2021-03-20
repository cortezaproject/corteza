package handlers

import (
	"encoding/base32"
	"fmt"
	"github.com/cortezaproject/corteza-server/auth/request"
	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/cortezaproject/corteza-server/system/types"
	"go.uber.org/zap"
	"math/rand"
	"net/url"
	"rsc.io/qr"
)

const (
	// session key where the secret is kept between requests
	totpSecretKey = "totpSecret"
)

// Handles MFA TOTP configuration form
//
// Where the TOTP QR & code are displayed and where
func (h AuthHandlers) mfaTotpConfigForm(req *request.AuthReq) (err error) {
	var (
		rawSecret [10]byte
		secret    string

		// this is more for debugging & development purposes
		// but it does not hurt to keep it here
		_, fresh = req.Request.URL.Query()["fresh"]
	)

	if s, has := req.Session.Values[totpSecretKey]; has && !fresh {
		// secret is already in the session and
		// there's no explicit request to change it
		secret = s.(string)
	} else {
		rand.Read(rawSecret[:])
		secret = base32.StdEncoding.EncodeToString(rawSecret[:])
		req.Session.Values[totpSecretKey] = secret
	}

	req.Data["secret"] = secret
	req.Data["enforced"] = h.Settings.MultiFactor.TOTP.Enforced
	req.Data["form"] = req.PopKV()
	req.Template = TmplMfaTotp
	return nil
}

// Handles MFA OTP form processing
func (h AuthHandlers) mfaTotpConfigProc(req *request.AuthReq) (err error) {
	req.RedirectTo = GetLinks().MfaTotpNewSecret
	req.SetKV(nil)

	var (
		user        *types.User
		secret, has = req.Session.Values[totpSecretKey]
	)

	if !has {
		return fmt.Errorf("no TOTP secret in session")
	}

	// Here is where code validation is done and where the secret is stored
	user, err = h.AuthService.ConfigureTOTP(
		auth.SetIdentityToContext(req.Context(), req.AuthUser.User),
		secret.(string),
		req.Request.PostFormValue("code"),
	)

	if err == nil {
		req.NewAlerts = append(req.NewAlerts, request.Alert{
			Type: "primary",
			Text: "Two factor authentication with TOTP enabled",
		})

		// Make sure we update User's data in the session
		req.AuthUser.User = user
		req.AuthUser.CompleteTOTP()
		req.AuthUser.Save(req.Session)

		h.Log.Info("TOTP code verified")
		req.RedirectTo = GetLinks().Security
		delete(req.Session.Values, totpSecretKey)
		return nil
	}

	switch {
	case errors.IsInvalidData(err):
		req.SetKV(map[string]string{
			"error": "Invalid code format",
		})
		return nil
	case errors.IsUnauthenticated(err):
		req.SetKV(map[string]string{
			"error": "Invalid code",
		})
		return nil

	default:
		// Just in case, delete secret if something unexpected happend
		delete(req.Session.Values, totpSecretKey)
		h.Log.Error("unhandled error", zap.Error(err))
		return err
	}
}

// Displays the QR PNG image
func (h AuthHandlers) mfaTotpConfigQR(req *request.AuthReq) (err error) {
	var (
		issuer      = h.Settings.MultiFactor.TOTP.Issuer
		secret, has = req.Session.Values[totpSecretKey]
	)

	if !has {
		return fmt.Errorf("no secret in session")
	}

	if len(issuer) == 0 {
		issuer = "Corteza"
	}

	account := req.AuthUser.User.Handle
	if len(account) == 0 {
		account = req.AuthUser.User.Email
	}

	URL, err := url.Parse("otpauth://totp")
	if err != nil {
		panic(err)
	}

	URL.Path += "/" + url.PathEscape(issuer) + ":" + url.PathEscape(account)

	params := url.Values{}
	params.Add("secret", secret.(string))
	params.Add("issuer", issuer)

	URL.RawQuery = params.Encode()

	code, err := qr.Encode(URL.String(), qr.Q)
	if err != nil {
		panic(err)
	}

	req.Status = -1
	_, err = req.Response.Write(code.PNG())
	return
}

// Handles MFA TOTP configuration form
//
// Where the TOTP QR & code are displayed and where
func (h AuthHandlers) mfaTotpDisableForm(req *request.AuthReq) (err error) {
	req.Data["form"] = req.PopKV()
	req.Template = TmplMfaTotpDisable
	return nil
}

// Handles MFA OTP form processing
func (h AuthHandlers) mfaTotpDisableProc(req *request.AuthReq) (err error) {
	req.RedirectTo = GetLinks().MfaTotpDisable
	req.SetKV(nil)

	var user *types.User

	// Here is where code validation is done and where the secret is stored
	user, err = h.AuthService.RemoveTOTP(
		req.Context(),
		req.AuthUser.User.ID,
		req.Request.PostFormValue("code"),
	)

	if err == nil {
		req.NewAlerts = append(req.NewAlerts, request.Alert{
			Type: "primary",
			Text: "Two factor authentication with TOTP disabled",
		})

		// Make sure we update User's data in the session
		req.AuthUser.User = user
		req.AuthUser.ResetTOTP()
		req.AuthUser.Save(req.Session)

		h.Log.Info("TOTP disabled")
		req.RedirectTo = GetLinks().Security
		return nil
	}

	switch {
	case errors.IsInvalidData(err):
		req.SetKV(map[string]string{
			"error": "Invalid code format",
		})
		return nil
	case errors.IsUnauthenticated(err):
		req.SetKV(map[string]string{
			"error": "Invalid code",
		})
		return nil

	default:
		h.Log.Error("unhandled error", zap.Error(err))
		return err
	}
}

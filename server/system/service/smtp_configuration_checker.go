package service

import (
	"context"
	"crypto/tls"
	"fmt"
	intAuth "github.com/cortezaproject/corteza/server/pkg/auth"
	"github.com/cortezaproject/corteza/server/pkg/mail"
	"github.com/cortezaproject/corteza/server/pkg/options"
	"github.com/cortezaproject/corteza/server/system/types"
	gomail "gopkg.in/mail.v2"
	htpl "html/template"
	"io/ioutil"
	"strings"
)

type (
	smtpConfigurationChecker struct {
		settings      *types.AppSettings
		ts            TemplateService
		opt           options.AuthOpt
		accessControl smtpCheckAccessController
	}

	smtpCheckAccessController interface {
		CanManageSettings(context.Context) bool
	}
)

func SmtpConfigurationChecker(s *types.AppSettings, ts TemplateService, ac accessController, opt options.AuthOpt) *smtpConfigurationChecker {
	return &smtpConfigurationChecker{
		settings:      s,
		ts:            ts,
		opt:           opt,
		accessControl: ac,
	}
}

// Check SMTP server configurations and send a test email
// to recipients if they're provided
func (svc smtpConfigurationChecker) Check(ctx context.Context, smtpConfigs *types.SmtpConfiguration) (checkResults *types.SmtpCheckResult, err error) {
	if !svc.accessControl.CanManageSettings(ctx) {
		return nil, fmt.Errorf("not allowed to check SMTP configurations")
	}

	var (
		tlsConfig = &tls.Config{}
	)

	checkResults = &types.SmtpCheckResult{}

	if smtpConfigs.Port == 0 {
		smtpConfigs.Port = 25
	}

	//check for validity of the host
	if !mail.IsValidHost(smtpConfigs.Host) {
		checkResults.Host = fmt.Sprintf("%s name is Invalid", smtpConfigs.Host)
	}

	// Applying TLS
	tlsConfig = &tls.Config{ServerName: smtpConfigs.Host}

	if smtpConfigs.TLSInsecure {
		tlsConfig.InsecureSkipVerify = true
	}

	if smtpConfigs.TLSServerName != "" {
		tlsConfig.ServerName = smtpConfigs.TLSServerName
	}

	checkResults.Server = mail.ConfigCheck(smtpConfigs.Host, smtpConfigs.Port, smtpConfigs.Username, smtpConfigs.Password, tlsConfig)

	//send the email there are recipients
	if checkResults.Server == "" {
		if len(smtpConfigs.Recipients) != 0 {
			checkResults.Send, err = svc.smtpSend(ctx, smtpConfigs.Recipients)
		}
	}

	return checkResults, err
}

func (svc smtpConfigurationChecker) smtpSend(ctx context.Context, recipients []string) (expected string, err error) {
	var (
		ntf      = mail.New()
		toHeader string
		// context with service user
		// we need this for retrieving & rendering email templates
		suCtx = intAuth.SetIdentityToContext(ctx, intAuth.ServiceUser())
	)

	if err = svc.procEmailRecipients(ntf, "To", recipients); err != nil {
		return "", err
	}

	toHeader = strings.Join(recipients, ",")
	ntf.SetAddressHeader("To", toHeader, "")

	st, ct, err := svc.findEmailTemplates(suCtx)
	// if we cannot find an email template
	if err != nil {
		ntf.SetHeader("Subject", "SMTP Configuration check")
		ntf.SetBody("text/html", "<h2 style=\"color: #568ba2;text-align: center;\">SMTP configurations check passed</h2>")

		err = mail.Send(ntf)

		if err != nil {
			return err.Error(), nil
		}

		return "", nil
	}

	subjectTmp, contentTmp, err := svc.procEmailTemplate(suCtx, st.ID, ct.ID)
	if err != nil {
		return "", err
	}

	ntf.SetHeader("Subject", string(subjectTmp))
	ntf.SetBody("text/html", string(contentTmp))

	err = mail.Send(ntf)

	if err != nil {
		return err.Error(), nil
	}

	return "", nil
}

func (svc smtpConfigurationChecker) procEmailRecipients(m *gomail.Message, field string, recipients []string) (err error) {
	var (
		email string
		rcpt  string
	)

	if len(recipients) == 0 {
		return
	}

	for _, rcpt = range recipients {

		email = strings.TrimSpace(rcpt)

		// Validate email here
		if !mail.IsValidAddress(email) {
			return fmt.Errorf("invalid recipient email address %s", email)
		}

	}

	m.SetHeader(field, recipients...)
	return nil
}

// procEmailTemplate processes Email address template based on the template's subject ID and content ID
func (svc smtpConfigurationChecker) procEmailTemplate(ctx context.Context, stId uint64, ctId uint64) (subjectTmp []byte, contentTmp []byte, err error) {
	// Prepare payload
	payload := map[string]interface{}{
		"Logo":    htpl.URL(svc.settings.General.Mail.Logo),
		"BaseURL": svc.opt.BaseURL,
	}

	// Render document
	subject, err := svc.ts.Render(ctx, stId, "text/plain", payload, nil)
	if err != nil {
		return nil, nil, err
	}

	content, err := svc.ts.Render(ctx, ctId, "text/plain", payload, nil)
	if err != nil {
		return nil, nil, err
	}

	subjectTmp, err = ioutil.ReadAll(subject)
	if err != nil {
		return nil, nil, err
	}

	contentTmp, err = ioutil.ReadAll(content)
	if err != nil {
		return nil, nil, err
	}

	return subjectTmp, contentTmp, nil
}

func (svc smtpConfigurationChecker) findEmailTemplates(ctx context.Context) (st *types.Template, ct *types.Template, err error) {
	var (
		hdl string
	)

	hdl = "smtp_configuration_check_subject"
	st, err = svc.ts.FindByHandle(ctx, hdl)
	if err != nil {
		return nil, nil, err
	}

	hdl = "smtp_configuration_check_content"
	ct, err = svc.ts.FindByHandle(ctx, hdl)
	if err != nil {
		return nil, nil, err
	}

	return st, ct, nil
}

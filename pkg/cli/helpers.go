package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/cortezaproject/corteza-server/pkg/auth"
	"github.com/cortezaproject/corteza-server/pkg/cli/options"
	"github.com/cortezaproject/corteza-server/pkg/http"
	"github.com/cortezaproject/corteza-server/pkg/mail"
)

func InitGeneralServices(smtpOpt *options.SMTPOpt, jwtOpt *options.JWTOpt, httpClientOpt *options.HttpClientOpt) {
	auth.SetupDefault(jwtOpt.Secret, int(jwtOpt.Expiry/time.Minute))
	mail.SetupDialer(smtpOpt.Host, smtpOpt.Port, smtpOpt.User, smtpOpt.Pass, smtpOpt.From)
	http.SetupDefaults(
		httpClientOpt.HttpClientTimeout,
		httpClientOpt.ClientTSLInsecure,
	)
}

func HandleError(err error) {
	if err == nil {
		return
	}

	_, _ = fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

package cli

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/http"
	"github.com/cortezaproject/corteza-server/internal/mail"
	"github.com/cortezaproject/corteza-server/pkg/cli/flags"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

func InitGeneralServices(logOpt *flags.LogOpt, smtpOpt *flags.SMTPOpt, jwtOpt *flags.JWTOpt, httpClientOpt *flags.HttpClientOpt) {
	// Reset logger's level to whatever we want
	var logLevel = zap.InfoLevel
	_ = logLevel.Set(logOpt.Level)
	logger.DefaultLevel.SetLevel(logLevel)

	auth.SetupDefault(jwtOpt.Secret, jwtOpt.Expiry)
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

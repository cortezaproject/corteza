package cli

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/internal/auth"
	"github.com/cortezaproject/corteza-server/internal/http"
	"github.com/cortezaproject/corteza-server/internal/mail"
	"github.com/cortezaproject/corteza-server/pkg/cli/options"
	"github.com/cortezaproject/corteza-server/pkg/logger"
)

func InitGeneralServices(logOpt *options.LogOpt, smtpOpt *options.SMTPOpt, jwtOpt *options.JWTOpt, httpClientOpt *options.HttpClientOpt) {
	// Reset logger's level to whatever we want
	var logLevel = zap.InfoLevel
	_ = logLevel.Set(logOpt.Level)
	logger.DefaultLevel.SetLevel(logLevel)

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

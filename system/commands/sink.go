package commands

import (
	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/spf13/cobra"
	"time"
)

// Will perform OpenID connect auto-configuration
func Sink() *cobra.Command {
	var (
		expires     string
		origin      string
		contentType string
		method      string
		maxBodySize int64
	)

	cmd := &cobra.Command{
		Use:   "sink",
		Short: "Sink",
	}

	signatureCmd := &cobra.Command{
		Use:   "signature",
		Short: "Creates signature for sink HTTP endpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			var (
				srup = service.SinkRequestUrlParams{
					Method:      method,
					Origin:      origin,
					ContentType: contentType,
					MaxBodySize: maxBodySize,
				}
			)

			if expires != "" {
				// validate expiration date if set
				if exp, err := time.Parse("2006-01-02", expires); err != nil {
					return err
				} else {
					srup.Expires = &exp
				}
			}

			cmd.Printf("%+v\n", srup)

			if su, err := service.DefaultSink.SignURL(srup); err != nil {
				return err
			} else {
				cmd.Println(su)
			}

			return nil
		},
	}

	signatureCmd.Flags().StringVar(
		&origin,
		"origin",
		"",
		"Origin of the request (arbitrary string, optional)")

	signatureCmd.Flags().StringVar(
		&contentType,
		"content-type",
		"",
		"Content type (optional)")

	signatureCmd.Flags().StringVar(
		&expires,
		"expires",
		"",
		"Date of expiration (YYYY-MM-DD, optional)")

	signatureCmd.Flags().StringVar(
		&method,
		"method",
		"GET",
		"HTTP method that will be used")

	signatureCmd.Flags().Int64Var(
		&maxBodySize,
		"max-body-size",
		0,
		"Max allowed body size")

	cmd.AddCommand(
		signatureCmd,
	)

	return cmd
}

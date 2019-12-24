package commands

import (
	"net/url"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/cortezaproject/corteza-server/pkg/auth"
)

// Will perform OpenID connect auto-configuration
func Sink() *cobra.Command {
	var (
		expires     string
		origin      string
		contentType string
		method      string
	)

	cmd := &cobra.Command{
		Use:   "sink",
		Short: "Sink",
	}

	signatureCmd := &cobra.Command{
		Use:   "signature",
		Short: "Creates signature for sink HTTP endpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			method = strings.ToUpper(method)

			if expires != "" {
				// validate expiration date if set
				if _, err := time.Parse("2006-01-02", expires); err != nil {
					return err
				}
			}

			v := url.Values{}
			v.Set("sign", auth.DefaultSigner.Sign(0, method, "/sink", contentType, origin, expires))
			v.Set("expires", expires)
			v.Set("content-type", contentType)
			v.Set("origin", origin)
			v.Set("method", method)

			// @todo add host & schema
			cmd.Println((&url.URL{
				Path:     "/sink",
				RawQuery: v.Encode()}).String())

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

	cmd.AddCommand(
		signatureCmd,
	)

	return cmd
}

package commands

import (
	"context"
	"time"

	"github.com/cortezaproject/corteza-server/system/service"
	"github.com/spf13/cobra"
)

// Will perform OpenID connect auto-configuration
func Sink(ctx context.Context, app serviceInitializer) *cobra.Command {
	var (
		expires string
		srup    = service.SinkRequestUrlParams{}
	)

	cmd := &cobra.Command{
		Use:   "sink",
		Short: "Sink",
	}

	signatureCmd := &cobra.Command{
		Use:     "signature",
		Short:   "Creates signature for sink HTTP endpoint",
		PreRunE: commandPreRunInitService(app),
		RunE: func(cmd *cobra.Command, args []string) error {
			if expires != "" {
				// validate expiration date if set
				if exp, err := time.Parse("2006-01-02", expires); err != nil {
					return err
				} else {
					srup.Expires = &exp
				}
			}

			if su, srup, err := service.DefaultSink.SignURL(srup); err != nil {
				return err
			} else {
				cmd.Println(su)

				cmd.Println("Sink request constraints:")
				if srup.SignatureInPath {
					cmd.Println(" - signature should be part of path")
				} else {
					cmd.Println(" - signature should be part of query-string")
				}

				if srup.Method != "" {
					cmd.Printf(" - expecting request method %q\n", srup.Method)

				}
				if srup.Expires != nil {
					cmd.Printf(" - signature expires at: %s\n", srup.Expires)

				}
				if srup.MaxBodySize > 0 {
					cmd.Printf(" - max request body size is %d Kb\n", srup.MaxBodySize/1024)

				} else {
					cmd.Println(" - body size is not limited")

				}
				if srup.ContentType != "" {
					cmd.Printf(" - expecting content type to be %q\n", srup.ContentType)

				}
				if srup.Path != "" {
					cmd.Printf(" - valid path under /sink: %q\n", srup.Path)

				}
			}

			return nil
		},
	}

	signatureCmd.Flags().StringVar(
		&srup.Origin,
		"origin",
		"",
		"Origin of the request (arbitrary string, optional)")

	signatureCmd.Flags().StringVar(
		&srup.ContentType,
		"content-type",
		"",
		"Content type (optional)")

	signatureCmd.Flags().StringVar(
		&srup.Path,
		"path",
		"",
		"Full sink request path (do not include /sink prefix, add / for just root)")

	signatureCmd.Flags().StringVar(
		&expires,
		"expires",
		"",
		"Date of expiration (YYYY-MM-DD, optional)")

	signatureCmd.Flags().StringVar(
		&srup.Method,
		"method",
		"",
		"HTTP method that will be used (optional)")

	signatureCmd.Flags().Int64Var(
		&srup.MaxBodySize,
		"max-body-size",
		0,
		"Max allowed body size")

	signatureCmd.Flags().BoolVar(
		&srup.SignatureInPath,
		"signature-in-path",
		false,
		"Include signature in a path instead of query string")

	cmd.AddCommand(
		signatureCmd,
	)

	return cmd
}

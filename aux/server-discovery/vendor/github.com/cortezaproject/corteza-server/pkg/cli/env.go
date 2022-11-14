package cli

import (
	"os"
	"path"
	"sort"

	"github.com/cortezaproject/corteza-server/pkg/errors"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// LoadEnv loads env-variables (KEY=VAL) from provided files and paths (by searching for .env file in the given path)
//
// Please note that loaded values
// DO NOT OVERRIDE the existing
// environmental variables
func LoadEnv(pp ...string) error {
	// preparse the input and try to figure out if .env should be appended
	var checked = make([]string, 0, len(pp))
	for _, p := range pp {
		if s, err := os.Stat(p); err != nil {
			return err
		} else if s.IsDir() {
			chk := path.Join(p, ".env")
			if _, err = os.Stat(chk); err == nil {
				// make sure only .env files
				checked = append(checked, chk)
			} else if !errors.Is(err, os.ErrNotExist) {
				return err
			}
		} else {
			checked = append(checked, p)
		}
	}

	if len(checked) == 0 {
		return nil
	}

	return godotenv.Load(checked...)
}

// EnvCommand outputs all loaded env variables
func EnvCommand() *cobra.Command {
	return &cobra.Command{
		Use:  "env",
		Long: "Outputs list (sorted by key) of of all environmental variables. Can be used for diagnosis and debugging",
		Run: func(cmd *cobra.Command, args []string) {
			kv := os.Environ()
			sort.Strings(kv)
			for _, l := range kv {
				cmd.Println(l)
			}
		},
	}
}

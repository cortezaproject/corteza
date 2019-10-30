package types

type (
	Settings struct {
		// UI related settings
		UI struct {
			// Emoji
			// @todo implementation
			NamespaceSwitcher struct {
				Enabled bool
			} `kv:"namespace-switcher"`
		} `kv:"ui"`

		// Message related settings
		Record struct {
			// @todo implementation
			Attachments struct {
				// What is max size (in MB, so: MaxSize x 2^20)
				MaxSize uint `kv:"max-size"`

				// List of mime-types we support,
				Mimetypes []string
			}
		}

		// Page related settings
		Page struct {
			// @todo implementation
			Attachments struct {
				// What is max size (in MB, so: MaxSize x 2^20)
				MaxSize uint `kv:"max-size"`

				// List of mime-types we support,
				Mimetypes []string
			}
		}
	}
)

package types

type (
	// Settings type is structured representation of current compose settings
	//
	// Raw settings keys are hypen (kebab) case, separated with a dot (.) that indicates sub-level
	// JSON properties for settings are NOT converted (lower-cased, etc...)
	// Use `json:"-"` tag to hide settings on REST endpoint
	Settings struct {
		// UI related settings
		UI struct {
			// Emoji
			// @todo implementation
			NamespaceSwitcher struct {
				Enabled     bool
				DefaultOpen bool
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

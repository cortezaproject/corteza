package types

type (
	// Settings type is structured representation of current messaging settings
	//
	// Raw settings keys are hypen (kebab) case, separated with a dot (.) that indicates sub-level
	// JSON properties for settings are NOT converted (lower-cased, etc...)
	// Use `json:"-"` tag to hide settings on REST endpoint
	Settings struct {
		// UI related settings
		UI struct {
			// Emoji
			// @todo implementation
			Emoji struct {
				Enabled bool
			}

			// In-browser notifications
			// @todo implementation
			BrowserNotifications struct {
				Enabled     bool
				Header      string
				MessageTrim uint `kv:"message-trim"`
			} `kv:"browser-notifications"`
		} `kv:"ui"`

		// Message related settings
		Message struct {
			// @todo implementation
			Attachments struct {
				// Completely disable attachments
				Enabled bool

				// What is max size (in MB, so: MaxSize x 2^20)
				MaxSize uint `kv:"max-size"`

				// List of mime-types we support,
				Mimetypes []string

				// Enable/disable individual attachment sources (mobile)
				Source struct {
					Gallery struct{ Enabled bool }
					Camera  struct{ Enabled bool }
				}
			}
		}
	}
)

package options

type (
	AppstoreOpt struct {
		URL    string `env:"APPSTORE_URL"`
		APIKey string `env:"APPSTORE_API_KEY"`
	}
)

func Appstore() *AppstoreOpt {
	o := &AppstoreOpt{}
	fill(o)
	return o
}

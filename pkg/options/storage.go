package options

type (
	StorageOpt struct {
		Path string `env:"STORAGE_PATH"`

		MinioEndpoint  string `env:"MINIO_ENDPOINT"`
		MinioSecure    bool   `env:"MINIO_SECURE"`
		MinioAccessKey string `env:"MINIO_ACCESS_KEY"`
		MinioSecretKey string `env:"MINIO_SECRET_KEY"`
		MinioSSECKey   string `env:"MINIO_SSEC_KEY"`
		MinioBucket    string `env:"MINIO_BUCKET"`
		MinioStrict    bool   `env:"MINIO_STRICT"`
	}
)

func Storage(pfix string) (o *StorageOpt) {
	o = &StorageOpt{
		Path: "var/store",

		// Make minio secure by default
		MinioSecure: true,

		// Run in struct mode:
		//  - do not create un-existing buckets
		MinioStrict: false,
	}

	fill(o)

	return
}

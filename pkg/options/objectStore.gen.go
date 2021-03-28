package options

// This file is auto-generated.
//
// Changes to this file may cause incorrect behavior and will be lost if
// the code is regenerated.
//
// Definitions file that controls how this file is generated:
// pkg/options/objectStore.yaml

type (
	ObjectStoreOpt struct {
		Path           string `env:"STORAGE_PATH"`
		MinioEndpoint  string `env:"MINIO_ENDPOINT"`
		MinioSecure    bool   `env:"MINIO_SECURE"`
		MinioAccessKey string `env:"MINIO_ACCESS_KEY"`
		MinioSecretKey string `env:"MINIO_SECRET_KEY"`
		MinioSSECKey   string `env:"MINIO_SSEC_KEY"`
		MinioBucket    string `env:"MINIO_BUCKET"`
		MinioBucketSep string `env:"MINIO_BUCKET_SEP"`
		MinioStrict    bool   `env:"MINIO_STRICT"`
	}
)

// ObjectStore initializes and returns a ObjectStoreOpt with default values
func ObjectStore() (o *ObjectStoreOpt) {
	o = &ObjectStoreOpt{
		Path:           "var/store",
		MinioSecure:    true,
		MinioBucketSep: "/",
		MinioStrict:    false,
	}

	fill(o)

	// Function that allows access to custom logic inside the parent function.
	// The custom logic in the other file should be like:
	// func (o *ObjectStore) Defaults() {...}
	func(o interface{}) {
		if def, ok := o.(interface{ Defaults() }); ok {
			def.Defaults()
		}
	}(o)

	return
}

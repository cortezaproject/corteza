package options

import (
	"github.com/cortezaproject/corteza-server/codegen/schema"
)

objectStore: schema.#optionsGroup & {
	handle: "object-store"
	title:  "Object (file) storage"

	intro: "The MinIO integration allows you to replace local storage with cloud storage. When configured, `STORAGE_PATH` is not needed."
	options: {
		path: {
			defaultValue: "var/store"
			description:  "Location where uploaded files are stored."
			env:          "STORAGE_PATH"
		}
		minioEndpoint: {
			env: "MINIO_ENDPOINT"
		}
		minioSecure: {
			type:          "bool"
			defaultGoExpr: "true"
			env:           "MINIO_SECURE"
		}
		minioAccessKey: {
			env: "MINIO_ACCESS_KEY"
		}
		minioSecretKey: {
			env: "MINIO_SECRET_KEY"
		}
		minioSSECKey: {
			env: "MINIO_SSEC_KEY"
		}
		minioBucket: {
			defaultValue: "{component}"
			description:  "`component` placeholder is replaced with service name (e.g system)."
			env:          "MINIO_BUCKET"
		}
		minioPathPrefix: {
			description: "`component` placeholder is replaced with service name (e.g system)."
			env:         "MINIO_PATH_PREFIX"
		}
		minioStrict: {
			type: "bool"
			env:  "MINIO_STRICT"
		}
	}
}

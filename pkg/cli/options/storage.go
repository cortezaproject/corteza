package options

type (
	StorageOpt struct {
		Path string `env:"STORAGE_PATH"`
	}
)

func Storage(pfix string) (o *StorageOpt) {
	o = &StorageOpt{
		Path: "var/store",
	}

	fill(o, pfix)

	return
}

package resource

type (
	provider interface {
		Fields() []string
		Next() (map[string]string, error)
	}

	ResourceDataset struct {
		*base

		P provider
	}
)

func NewResourceDataset(name string, p provider) *ResourceDataset {
	r := &ResourceDataset{base: &base{}}
	r.P = p

	r.SetResourceType(DATA_SOURCE_RESOURCE_TYPE)
	r.AddIdentifier(identifiers(name)...)

	return r
}

package resource

type (
	provider interface {
		Fields() []string
		Count() uint64
		Reset() error
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

	r.SetResourceType(DataSourceResourceType)
	r.AddIdentifier(identifiers(name)...)

	return r
}

func (r *ResourceDataset) Resource() interface{} {
	return nil
}

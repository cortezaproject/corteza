package importer

type (
	Interface interface {
		CastSet(set interface{}) error
		Cast(handle string, set interface{}) error
	}

	PermissionImporter interface {
		CastSet(string, string, interface{}) error
		CastResourcesSet(string, interface{}) error
	}
)

func NormalizeHandle(in string) string {
	return in
}

func IsValidHandle(in string) bool {
	return true
}

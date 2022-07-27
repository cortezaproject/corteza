package revisions

type (
	Operation = string
)

const (
	unknown     Operation = ""
	Created               = "created"
	Updated               = "updated"
	SoftDeleted           = "soft-deleted"
	Restored              = "restored"
	HardDeleted           = "hard-deleted"
)

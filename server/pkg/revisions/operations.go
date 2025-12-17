package revisions

type (
	Operation = string
	Status    = string
)

const (
	unknown     Operation = ""
	Created     Operation = "created"
	Updated     Operation = "updated"
	SoftDeleted Operation = "soft-deleted"
	Undeleted   Operation = "undeleted"
	HardDeleted Operation = "hard-deleted"
)

const (
	StatusNone  Status = ""      // Normal revision - Operation tells you what happened
	StatusDraft Status = "draft" // Draft revision from local storage sync
)

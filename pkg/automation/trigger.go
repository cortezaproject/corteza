package automation

import (
	"time"
)

//
//		  _____                                _           _
//		 |  __ \                              | |         | |
//		 | |  | | ___ _ __  _ __ ___  ___ __ _| |_ ___  __| |
//		 | |  | |/ _ \ '_ \| '__/ _ \/ __/ _` | __/ _ \/ _` |
//		 | |__| |  __/ |_) | | |  __/ (_| (_| | ||  __/ (_| |
//		 |_____/ \___| .__/|_|  \___|\___\__,_|\__\___|\__,_|
//					 | |
//					 |_|
//
//
//
// This automation package is kept only to aid in export of deprecated
// trigger & script format from database into files
//
// Package was refactored and replaced by pkg/corredor, pkg/eventbus and pkg/scheduler
//
// Scheduled for removal in 2020.6
//

type (
	Event string

	Trigger struct {
		ID uint64 `json:"triggerID,string" db:"id"`

		// Resource that triggered the event
		//  - "compose:" (unspec, general)
		//  - "compose:record"
		//  - "compose:namespace"
		Resource string `json:"resource" db:"resource"`

		// Event name, arbitrary string
		//  - "manual"
		//  - "interval"
		//  - "deferred"
		//  - "before"
		//  - "after"
		Event string `json:"event" db:"event"`

		// Arbitrary data for trigger condition
		//
		// It is caller's responsibility to encode, decode and verify conditions
		Condition string `json:"condition" db:"event_condition"`

		ScriptID uint64 `json:"scriptID,string" db:"rel_script"`

		// Is trigger enabled or disabled?
		Enabled bool `json:"enabled" db:"enabled"`

		Weight int `json:"weight" db:"weight"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt"`
		CreatedBy uint64     `db:"created_by" json:"createdBy,string" `
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		UpdatedBy uint64     `db:"updated_by" json:"updatedBy,string,omitempty" `
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
		DeletedBy uint64     `db:"deleted_by" json:"deletedBy,string,omitempty" `
	}
)

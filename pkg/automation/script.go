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
	Script struct {
		ID uint64 `json:"scriptID,string" db:"id"`

		NamespaceID uint64 `json:"namespaceID,string" db:"rel_namespace"`

		Name string `json:"name" db:"name"`

		// (URL) Where did we get the source from?
		SourceRef string `json:"sourceRef" db:"source_ref"`

		// Code
		Source string `json:"source" db:"source"`

		// No need to wait for script to return the value
		Async bool `json:"async" db:"async"`

		// Who is running this script?
		// Leave it at 0 for the current user (security invoker) or
		// set ID of specific user (security definer)
		RunAs uint64 `json:"runAs,string" db:"rel_runner"`

		// Where can we run this script? user-agent? corredor service?
		RunInUA bool `json:"runInUA" db:"run_in_ua"`

		// Are you doing something that can take more time?
		// specify timeout (in milliseconds)
		Timeout uint `json:"timeout" db:"timeout"`

		// Is it critical to run this script successfully?
		Critical bool `json:"critical" db:"critical"`

		Enabled bool `json:"enabled" db:"enabled"`

		CreatedAt time.Time  `db:"created_at" json:"createdAt"`
		CreatedBy uint64     `db:"created_by" json:"createdBy,string" `
		UpdatedAt *time.Time `db:"updated_at" json:"updatedAt,omitempty"`
		UpdatedBy uint64     `db:"updated_by" json:"updatedBy,string,omitempty" `
		DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
		DeletedBy uint64     `db:"deleted_by" json:"deletedBy,string,omitempty" `

		// Serves as container for valid Triggers for runnable scripts internal cache
		// and for as a transport on create/update operations
		//
		// Node: on c/u op., we currently just merge current state with the given list,
		// w/o updating the rest of the script's Triggers
		Triggers []*Trigger
	}
)

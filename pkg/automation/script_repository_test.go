// +build integration

package automation

import (
	"os"
	"testing"
	"time"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/internal/test"
)

func TestScriptRepository_findByID(t *testing.T) {
	factory.Database.Add("compose", os.Getenv("COMPOSE_DB_DSN"))

	var (
		scriptID = factory.Sonyflake.NextID()

		err    error
		ss     ScriptSet
		script *Script

		db    = factory.Database.MustGet("compose")
		srepo = ScriptRepository("compose")
	)

	db.Begin()
	defer db.Rollback()

	_, err = srepo.findByID(db, scriptID)
	test.Error(t, err, "expecting error")

	//  runnable script
	script = &Script{ID: scriptID, CreatedAt: time.Now(), Enabled: true}
	err = srepo.create(db, script)
	test.NoError(t, err, "unexpected error: %v")

	// find this script now
	script, err = srepo.findByID(db, scriptID)
	test.NoError(t, err, "script not found")
	test.Assert(t, script != nil && script.ID == scriptID, "script not found (in runnable scripts)")

	// find the script through find func
	ss, _, err = srepo.find(db, ScriptFilter{})
	test.NoError(t, err, "unexpected error: %v")
	test.Assert(t, len(ss) > 0, "could not find the script")
	test.Assert(t, ss.FindByID(scriptID) != nil, "could not find the script")

	err = srepo.update(db, script)
	test.NoError(t, err, "unexpected error when updating the script: %v")
}

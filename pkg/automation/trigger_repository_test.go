// +build integration

package automation

import (
	"os"
	"testing"
	"time"

	"github.com/titpetric/factory"

	"github.com/cortezaproject/corteza-server/internal/test"
)

func TestTriggerRepository_findByID(t *testing.T) {
	factory.Database.Add("compose", os.Getenv("COMPOSE_DB_DSN"))

	var (
		nsID      = factory.Sonyflake.NextID()
		scriptID  = factory.Sonyflake.NextID()
		triggerID = factory.Sonyflake.NextID()

		err     error
		trigger *Trigger
		tt      TriggerSet
		script  *Script

		db    = factory.Database.MustGet("compose")
		trepo = TriggerRepository("compose")
		srepo = ScriptRepository("compose")
	)

	db.Begin()
	defer db.Rollback()

	_, err = db.Exec("insert into compose_namespace (id, name, slug, meta, enabled) values (?, 'test ns', 'test slug', '{}', true)", nsID)
	test.NoError(t, err, "unexpected error")

	_, err = trepo.findByID(db, triggerID)
	test.Error(t, err, "expecting error")

	// should have script present
	script = &Script{ID: scriptID, CreatedAt: time.Now(), NamespaceID: nsID}
	err = srepo.create(db, script)
	test.NoError(t, err, "unexpected error: %v")

	//  runnable trigger
	trigger = &Trigger{ID: triggerID, Resource: "test", Event: "test", CreatedAt: time.Now(), Enabled: true, ScriptID: scriptID}
	err = trepo.replace(db, trigger)
	test.NoError(t, err, "unexpected error: %v")

	// find this trigger now
	trigger, err = trepo.findByID(db, triggerID)
	test.NoError(t, err, "trigger not found")
	test.Assert(t, trigger != nil && trigger.ID == triggerID, "trigger not found")

	tt, err = trepo.findRunnable(db)
	test.NoError(t, err, "unexpected error: %v")
	test.Assert(t, tt.FindByID(triggerID) != nil, "trigger not found (in runnable)")

	// find the trigger through find func
	tt, _, err = trepo.find(db, TriggerFilter{
		ScriptID: script.ID,
		Resource: "test",
		Event:    "test",
	})

	test.NoError(t, err, "unexpected error: %v")
	test.Assert(t, len(tt) == 1, "could not find the trigger")
	test.Assert(t, tt[0].ID == triggerID, "unexpected trigger")

	err = trepo.replace(db, trigger)
	test.NoError(t, err, "unexpected error when replacing the trigger: %v")

	err = trepo.deleteByScriptID(db, scriptID)
	test.NoError(t, err, "unexpected error when updating the trigger: %v")

	// find this trigger now
	trigger, err = trepo.findByID(db, triggerID)
	test.Error(t, err, "expecting error")

}

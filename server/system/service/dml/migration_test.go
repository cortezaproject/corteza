package dml

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	composeTypes "github.com/crusttech/human/server/compose/types"
)

// The internal connection is the only place the source record is still whole.
// If it stops handing identity and provenance to the importer, links cannot be
// remapped and migrated records silently get the publisher's name on them.
func TestInternalConnCarriesIdentityAndProvenance(t *testing.T) {
	var (
		req     = require.New(t)
		created = time.Date(2021, 3, 4, 5, 6, 7, 0, time.UTC)
		updated = time.Date(2022, 8, 9, 10, 11, 12, 0, time.UTC)
	)

	iter := &composeRecordIterator{pos: -1, records: composeTypes.RecordSet{
		{
			ID:        1001,
			OwnedBy:   77,
			CreatedAt: created,
			CreatedBy: 42,
			UpdatedAt: &updated,
			UpdatedBy: 43,
			Values: composeTypes.RecordValueSet{
				{Name: "title", Value: "hello"},
			},
		},
		{
			// never edited: updatedAt must stay nil, not become "now"
			ID:        1002,
			CreatedAt: created,
			CreatedBy: 42,
		},
	}}

	row := &importRow{}

	req.True(iter.Next(t.Context()))
	row.reset()
	req.NoError(iter.Scan(row))

	v, err := row.GetValue("title", 0)
	req.NoError(err)
	req.Equal("hello", v)

	mr := readRow(row)
	req.Equal(uint64(1001), mr.srcID)
	req.Equal(uint64(77), mr.ownedBy)
	req.Equal(created, mr.createdAt)
	req.Equal(uint64(42), mr.createdBy)
	req.NotNil(mr.updatedAt)
	req.Equal(updated, *mr.updatedAt)
	req.Equal(uint64(43), mr.updatedBy)

	req.True(iter.Next(t.Context()))
	row.reset()
	req.NoError(iter.Scan(row))

	mr = readRow(row)
	req.Equal(uint64(1002), mr.srcID)
	req.Nil(mr.updatedAt)

	req.False(iter.Next(t.Context()))
}

// The stamp a migrated record is written with has to be the stamp clearing
// looks for; if the two drift, a retried publish doubles every row.
func TestMigratedRecordsAreRecognisable(t *testing.T) {
	req := require.New(t)

	mr := &migratedRecord{srcID: 1001}
	req.True(isMigrated(&composeTypes.Record{Meta: mr.meta()}))

	req.False(isMigrated(&composeTypes.Record{}))
	req.False(isMigrated(&composeTypes.Record{Meta: map[string]any{"something": "else"}}))
	req.False(isMigrated(nil))
}

// The whole point of the second pass: a link value that named a record in the
// previous namespace must come out naming the record that was migrated from it.
func TestFixupRemapsLinksAndRestoresProvenance(t *testing.T) {
	var (
		req     = require.New(t)
		created = time.Date(2021, 3, 4, 5, 6, 7, 0, time.UTC)

		mig = &Migration{idMap: map[uint64]uint64{
			// old target record -> the id it got in the new namespace
			500: 900,
			501: 901,
		}}

		mm = &migratedModule{
			module: &composeTypes.Module{ID: 7, Handle: "agent_t"},
			records: []*migratedRecord{{
				srcID:     100,
				newID:     200,
				ownedBy:   77,
				createdAt: created,
				createdBy: 42,
				links: composeTypes.RecordValueSet{
					{Name: "link", Value: "500", Place: 0},
					{Name: "link", Value: "501", Place: 1},
				},
			}},
		}

		rec = &composeTypes.Record{
			ID:     200,
			Values: composeTypes.RecordValueSet{{Name: "title", Value: "hello"}},
		}
	)

	upd, unresolved := mig.fixup(mm, map[uint64]*composeTypes.Record{200: rec})

	req.Empty(unresolved)
	req.Len(upd, 1)

	req.Equal("900", upd[0].Values.Get("link", 0).Value)
	req.Equal(uint64(900), upd[0].Values.Get("link", 0).Ref)
	req.Equal("901", upd[0].Values.Get("link", 1).Value)
	req.Equal("hello", upd[0].Values.Get("title", 0).Value)

	req.Equal(uint64(77), upd[0].OwnedBy)
	req.Equal(created, upd[0].CreatedAt)
	req.Equal(uint64(42), upd[0].CreatedBy)
	req.Nil(upd[0].UpdatedAt)
}

// A link whose target never migrated cannot be silently emptied: the publish
// has to stop so the draft can be retried once the target module is mapped.
func TestFixupReportsLinksItCannotRemap(t *testing.T) {
	var (
		req = require.New(t)

		mig = &Migration{idMap: map[uint64]uint64{}}

		mm = &migratedModule{
			module: &composeTypes.Module{ID: 7, Handle: "agent_t"},
			records: []*migratedRecord{{
				srcID: 100,
				newID: 200,
				links: composeTypes.RecordValueSet{
					{Name: "link", Value: "500"},
				},
			}},
		}

		rec = &composeTypes.Record{ID: 200}
	)

	upd, unresolved := mig.fixup(mm, map[uint64]*composeTypes.Record{200: rec})

	req.Len(unresolved, 1)
	req.Contains(unresolved[0], "agent_t.link")
	req.Contains(unresolved[0], "500")

	// the row is still written -- provenance is worth keeping even on a run
	// that will fail, and the retry clears it anyway
	req.Len(upd, 1)
	req.Nil(upd[0].Values.Get("link", 0))
}

// Records the migration did not create belong to whoever typed them into the
// draft; the fixup must not touch them.
func TestFixupLeavesForeignRecordsAlone(t *testing.T) {
	var (
		req = require.New(t)
		mig = &Migration{idMap: map[uint64]uint64{}}
		mm  = &migratedModule{
			module:  &composeTypes.Module{ID: 7, Handle: "agent_t"},
			records: []*migratedRecord{{srcID: 100, newID: 200}},
		}
	)

	upd, unresolved := mig.fixup(mm, map[uint64]*composeTypes.Record{
		999: {ID: 999},
	})

	req.Empty(unresolved)
	req.Empty(upd)
}

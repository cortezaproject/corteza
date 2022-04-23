package crs

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/data"
	"github.com/cortezaproject/corteza-server/pkg/id"
	"github.com/cortezaproject/corteza-server/pkg/logger"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms"
	"github.com/cortezaproject/corteza-server/store/adapters/rdbms/drivers/sqlite"
)

const (
	repeatReads  = 1000
	totalRecords = 100
)

func TestIteratorNG(t *testing.T) {
	ctx := context.Background()
	ctx = logger.ContextWithValue(ctx, logger.MakeDebugLogger())
	cfg, err := sqlite.NewConfig("sqlite3://file::memory:?cache=shared&mode=memory")
	//cfg, err := postgres.NewConfig("postgres://darh@localhost:5432/corteza_2022_3?sslmode=disable&")
	if err != nil {
		t.Errorf("could not connect: %v", err)
	}

	db, err := rdbms.Connect(ctx, logger.Default(), cfg)
	if err != nil {
		t.Errorf("could not connect: %v", err)
	}

	_, err = db.ExecContext(ctx, `
	CREATE TABLE test_tbl (
		id NUMERIC NOT NULL,
		created_at TIMESTAMP,
		updated_at TIMESTAMP,
		phy TEXT,
		meta JSON
	)
	`)
	if err != nil {
		t.Errorf("could not create table: %v", err)
	}

	m := &data.Model{
		Ident: "test_tbl",
		Attributes: data.AttributeSet{
			&data.Attribute{Ident: sysID, Type: &data.TypeID{}, Store: &data.StoreCodecAlias{Ident: "id"}},
			&data.Attribute{Ident: sysCreatedAt, Type: &data.TypeTimestamp{}, Store: &data.StoreCodecAlias{Ident: "created_at"}},
			&data.Attribute{Ident: sysUpdatedAt, Type: &data.TypeTimestamp{}, Store: &data.StoreCodecAlias{Ident: "updated_at"}},
			&data.Attribute{Ident: "foo", Type: &data.TypeText{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
			&data.Attribute{Ident: "bar", Type: &data.TypeNumber{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
			&data.Attribute{Ident: "baz", Type: &data.TypeBoolean{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
			&data.Attribute{Ident: "bbb", Type: &data.TypeUUID{}, Store: &data.StoreCodecStdRecordValueJSON{Ident: "meta"}},
			&data.Attribute{Ident: "phy", Type: &data.TypeText{}, Store: &data.StoreCodecPlain{}},
		},
	}

	d := CRS(m, db)
	_, err = db.Exec("DELETE FROM test_tbl")
	if err != nil {
		t.Errorf("could not truncate: %v", err)
	}

	for i := 0; i < totalRecords; i++ {
		r := &types.Record{ID: id.Next(), CreatedAt: time.Now()}

		r.Values = r.Values.Set(&types.RecordValue{Name: "foo", Value: "FooVal1"})
		r.Values = r.Values.Set(&types.RecordValue{Name: "bar", Value: "34234"})
		r.Values = r.Values.Set(&types.RecordValue{Name: "baz", Value: "true"})
		r.Values = r.Values.Set(&types.RecordValue{Name: "phy", Value: "lly"})

		if err = d.Create(ctx, r); err != nil {
			t.Errorf("could not create: %v", err)
		}
	}

	t.Logf("records (%dx) created", totalRecords)

	i := &iterator{
		ms: d,
	}

	for rep := 0; rep < repeatReads; rep++ {
		for i.Next(ctx) {
			r := new(types.Record)
			if err = i.Scan(r); err != nil {
				t.Errorf("could not scan: %v", err)
			}
		}
	}
}

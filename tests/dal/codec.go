package dal

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza-server/compose/types"
	"github.com/cortezaproject/corteza-server/pkg/dal"
	"github.com/stretchr/testify/require"
)

func RecordCodec(t *testing.T, ctx context.Context, d dal.Connection) {
	var (
		req = require.New(t)

		m = &dal.Model{
			Ident: "crs_test_codec",
			Attributes: dal.AttributeSet{
				&dal.Attribute{Ident: "ID", Type: &dal.TypeID{}, Store: &dal.CodecAlias{Ident: "id"}, PrimaryKey: true},
				&dal.Attribute{Ident: "createdAt", Type: &dal.TypeTimestamp{}, Store: &dal.CodecAlias{Ident: "created_at"}},
				&dal.Attribute{Ident: "updatedAt", Type: &dal.TypeTimestamp{}, Store: &dal.CodecAlias{Ident: "updated_at"}},

				&dal.Attribute{Ident: "vID", Type: &dal.TypeID{}, Store: &dal.CodecRecordValueSetJSON{Ident: "meta"}},
				&dal.Attribute{Ident: "vRef", Type: &dal.TypeRef{}, Store: &dal.CodecRecordValueSetJSON{Ident: "meta"}},
				&dal.Attribute{Ident: "vTimestamp", Type: &dal.TypeTimestamp{}, Store: &dal.CodecRecordValueSetJSON{Ident: "meta"}},
				&dal.Attribute{Ident: "vTime", Type: &dal.TypeTime{}, Store: &dal.CodecRecordValueSetJSON{Ident: "meta"}},
				&dal.Attribute{Ident: "vDate", Type: &dal.TypeDate{}, Store: &dal.CodecRecordValueSetJSON{Ident: "meta"}},
				&dal.Attribute{Ident: "vNumber", Type: &dal.TypeNumber{}, Store: &dal.CodecRecordValueSetJSON{Ident: "meta"}},
				&dal.Attribute{Ident: "vText", Type: &dal.TypeText{}, Store: &dal.CodecRecordValueSetJSON{Ident: "meta"}},
				&dal.Attribute{Ident: "vBoolean_T", Type: &dal.TypeBoolean{}, Store: &dal.CodecRecordValueSetJSON{Ident: "meta"}},
				&dal.Attribute{Ident: "vBoolean_F", Type: &dal.TypeBoolean{}, Store: &dal.CodecRecordValueSetJSON{Ident: "meta"}},
				&dal.Attribute{Ident: "vEnum", Type: &dal.TypeEnum{}, Store: &dal.CodecRecordValueSetJSON{Ident: "meta"}},
				&dal.Attribute{Ident: "vGeometry", Type: &dal.TypeGeometry{}, Store: &dal.CodecRecordValueSetJSON{Ident: "meta"}},
				&dal.Attribute{Ident: "vJSON", Type: &dal.TypeJSON{}, Store: &dal.CodecRecordValueSetJSON{Ident: "meta"}},
				&dal.Attribute{Ident: "vBlob", Type: &dal.TypeBlob{}, Store: &dal.CodecRecordValueSetJSON{Ident: "meta"}},
				&dal.Attribute{Ident: "vUUID", Type: &dal.TypeUUID{}, Store: &dal.CodecRecordValueSetJSON{Ident: "meta"}},
				&dal.Attribute{Ident: "pID", Type: &dal.TypeID{}, Store: &dal.CodecPlain{}},
				&dal.Attribute{Ident: "pRef", Type: &dal.TypeRef{}, Store: &dal.CodecPlain{}},
				&dal.Attribute{Ident: "pTimestamp_TZT", Type: &dal.TypeTimestamp{Timezone: true}, Store: &dal.CodecPlain{}},
				&dal.Attribute{Ident: "pTimestamp_TZF", Type: &dal.TypeTimestamp{Timezone: false}, Store: &dal.CodecPlain{}},
				&dal.Attribute{Ident: "pTime", Type: &dal.TypeTime{}, Store: &dal.CodecPlain{}},
				&dal.Attribute{Ident: "pDate", Type: &dal.TypeDate{}, Store: &dal.CodecPlain{}},
				&dal.Attribute{Ident: "pNumber", Type: &dal.TypeNumber{}, Store: &dal.CodecPlain{}},
				&dal.Attribute{Ident: "pText", Type: &dal.TypeText{}, Store: &dal.CodecPlain{}},
				&dal.Attribute{Ident: "pBoolean_T", Type: &dal.TypeBoolean{}, Store: &dal.CodecPlain{}},
				&dal.Attribute{Ident: "pBoolean_F", Type: &dal.TypeBoolean{}, Store: &dal.CodecPlain{}},
				&dal.Attribute{Ident: "pEnum", Type: &dal.TypeEnum{}, Store: &dal.CodecPlain{}},
				&dal.Attribute{Ident: "pGeometry", Type: &dal.TypeGeometry{}, Store: &dal.CodecPlain{}},
				&dal.Attribute{Ident: "pJSON", Type: &dal.TypeJSON{}, Store: &dal.CodecPlain{}},
				&dal.Attribute{Ident: "pBlob", Type: &dal.TypeBlob{}, Store: &dal.CodecPlain{}},
				&dal.Attribute{Ident: "pUUID", Type: &dal.TypeUUID{}, Store: &dal.CodecPlain{}},
			},
		}

		rIn  = types.Record{ID: 42}
		err  error
		rOut *types.Record

		piTime time.Time
	)

	piTime, err = time.Parse("2006-01-02T15:04:05", "2006-01-02T15:04:05")
	req.NoError(err)
	piTime = piTime.UTC()

	rIn.CreatedAt = piTime
	rIn.UpdatedAt = &piTime

	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vID", Value: "34324"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vRef", Value: "32423"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vTimestamp", Value: "2022-01-01T10:10:10+02:00"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vTime", Value: "04:10:10+04:00"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vDate", Value: "2022-01-01"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vNumber", Value: "2423423"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vText", Value: "lorm ipsum "})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vBoolean_T", Value: "true"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vBoolean_F", Value: "false"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vEnum", Value: "abc"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vGeometry", Value: `{"lat":1,"lng":1}`})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vJSON", Value: `[{"bool":true"}]`})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vBlob", Value: "0110101"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vUUID", Value: "ba485865-54f9-44de-bde8-6965556c022a"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pID", Value: "34324"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pRef", Value: "32423"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pTimestamp_TZF", Value: "2022-02-01T10:10:10"})

	// @todo how (if at all) should we know if underlying DB supports timezone?
	//rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pTimestamp_TZT", Value: "2022-02-01T10:10:10"})

	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pTime", Value: "06:06:06"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pDate", Value: "2022-01-01"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pNumber", Value: "2423423"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pText", Value: "lorm ipsum "})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pBoolean_T", Value: "true"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pBoolean_F", Value: "false"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pEnum", Value: "abc"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pGeometry", Value: `{"lat":1,"lng":1}`})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pJSON", Value: `[{"bool":true"}]`})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pBlob", Value: "0110101"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "pUUID", Value: "ba485865-54f9-44de-bde8-6965556c022a"})
	rIn.Values = rIn.Values.GetClean()

	req.NoError(d.Create(ctx, m, &rIn))

	rOut = new(types.Record)
	req.NoError(d.Lookup(ctx, m, dal.PKValues{"id": rIn.ID}, rOut))

	{
		// normalize timezone on timestamps
		rOut.CreatedAt = rOut.CreatedAt.UTC()
		aux := rOut.UpdatedAt.UTC()
		rOut.UpdatedAt = &aux
	}

	for _, attr := range m.Attributes {
		vIn, err := rIn.GetValue(attr.Ident, 0)
		req.NoError(err)
		vOut, err := rOut.GetValue(attr.Ident, 0)
		req.NoError(err)
		req.Equal(vIn, vOut, "values for attribute %q are not equal", attr.Ident)
	}
}

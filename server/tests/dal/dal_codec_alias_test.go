package dal

import (
	"context"
	"testing"
	"time"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/dal"
)

func Test_dal_codec_alias(t *testing.T) {
	t.Skip("needs refactoring")

	model := &dal.Model{
		Ident: "compose_record_partitioned",
		Attributes: dal.AttributeSet{
			&dal.Attribute{Ident: "ID", Type: &dal.TypeID{}, Store: &dal.CodecAlias{Ident: "id"}, PrimaryKey: true},

			&dal.Attribute{Ident: "intr_vID", Type: &dal.TypeID{}, Store: &dal.CodecAlias{Ident: "vID"}},
			&dal.Attribute{Ident: "intr_vRef", Type: &dal.TypeRef{}, Store: &dal.CodecAlias{Ident: "vRef"}},
			&dal.Attribute{Ident: "intr_vTimestamp", Type: &dal.TypeTimestamp{}, Store: &dal.CodecAlias{Ident: "vTimestamp"}},
			&dal.Attribute{Ident: "intr_vTime", Type: &dal.TypeTime{}, Store: &dal.CodecAlias{Ident: "vTime"}},
			&dal.Attribute{Ident: "intr_vDate", Type: &dal.TypeDate{}, Store: &dal.CodecAlias{Ident: "vDate"}},
			&dal.Attribute{Ident: "intr_vNumber", Type: &dal.TypeNumber{}, Store: &dal.CodecAlias{Ident: "vNumber"}},
			&dal.Attribute{Ident: "intr_vText", Type: &dal.TypeText{}, Store: &dal.CodecAlias{Ident: "vText"}},
			&dal.Attribute{Ident: "intr_vBoolean_T", Type: &dal.TypeBoolean{}, Store: &dal.CodecAlias{Ident: "vBoolean_T"}},
			&dal.Attribute{Ident: "intr_vBoolean_F", Type: &dal.TypeBoolean{}, Store: &dal.CodecAlias{Ident: "vBoolean_F"}},
			&dal.Attribute{Ident: "intr_vEnum", Type: &dal.TypeEnum{}, Store: &dal.CodecAlias{Ident: "vEnum"}},
			&dal.Attribute{Ident: "intr_vGeometry", Type: &dal.TypeGeometry{}, Store: &dal.CodecAlias{Ident: "vGeometry"}},
			&dal.Attribute{Ident: "intr_vJSON", Type: &dal.TypeJSON{}, Store: &dal.CodecAlias{Ident: "vJSON"}},
			&dal.Attribute{Ident: "intr_vBlob", Type: &dal.TypeBlob{}, Store: &dal.CodecAlias{Ident: "vBlob"}},
			&dal.Attribute{Ident: "intr_vUUID", Type: &dal.TypeUUID{}, Store: &dal.CodecAlias{Ident: "vUUID"}},
		},
	}

	var (
		rIn  = types.Record{ID: 42, CreatedAt: time.Now()}
		rOut = types.Record{}
	)

	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "intr_vID", Value: "34324"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "intr_vRef", Value: "32423"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "intr_vTimestamp", Value: "2022-01-01T10:10:10"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "intr_vTime", Value: "04:10:10"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "intr_vDate", Value: "2022-01-01"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "intr_vNumber", Value: "2423423"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "intr_vText", Value: "lorm ipsum "})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "intr_vBoolean_T", Value: "1"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "intr_vBoolean_F", Value: ""})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "intr_vEnum", Value: "abc"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "intr_vGeometry", Value: `{"lat":1,"lng":1}`})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "intr_vJSON", Value: `[{"bool":true"}]`})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "intr_vBlob", Value: "0110101"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "intr_vUUID", Value: "ba485865-54f9-44de-bde8-6965556c022a"})

	bootstrap(t, func(ctx context.Context, t *testing.T, h helper, svc dalService) {
		h.cleanupDal()

		h.a.NoError(svc.ReplaceModel(ctx, model))

		h.a.NoError(svc.Create(ctx, model.ToFilter(), dal.CreateOperations(), &rIn))

		h.a.NoError(svc.Lookup(ctx, model.ToFilter(), dal.LookupOperations(), dal.PKValues{"id": rIn.ID}, &rOut))

		for _, inVal := range rIn.Values {
			outVal := rOut.Values.Get(inVal.Name, 0)
			h.a.Equal(inVal.Value, outVal.Value, "field", inVal.Name)
		}
	})
}

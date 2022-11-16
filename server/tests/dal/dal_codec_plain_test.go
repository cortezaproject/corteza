package dal

import (
	"context"
	"testing"

	"github.com/cortezaproject/corteza/server/compose/types"
	"github.com/cortezaproject/corteza/server/pkg/dal"
)

func Test_dal_codec_plain(t *testing.T) {
	t.Skip("needs refactoring")

	model := &dal.Model{
		Ident: "compose_record_partitioned",
		Attributes: dal.AttributeSet{
			&dal.Attribute{Ident: "ID", Type: &dal.TypeID{}, Store: &dal.CodecAlias{Ident: "id"}, PrimaryKey: true},

			&dal.Attribute{Ident: "vID", Type: &dal.TypeID{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vRef", Type: &dal.TypeRef{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vTimestamp", Type: &dal.TypeTimestamp{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vTime", Type: &dal.TypeTime{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vDate", Type: &dal.TypeDate{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vNumber", Type: &dal.TypeNumber{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vText", Type: &dal.TypeText{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vBoolean_T", Type: &dal.TypeBoolean{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vBoolean_F", Type: &dal.TypeBoolean{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vEnum", Type: &dal.TypeEnum{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vGeometry", Type: &dal.TypeGeometry{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vJSON", Type: &dal.TypeJSON{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vBlob", Type: &dal.TypeBlob{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vUUID", Type: &dal.TypeUUID{}, Store: &dal.CodecPlain{}},
		},
	}

	var (
		rIn  = types.Record{ID: 42}
		rOut = types.Record{}
	)

	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vID", Value: "34324"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vRef", Value: "32423"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vTimestamp", Value: "2022-01-01T10:10:10"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vTime", Value: "04:10:10"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vDate", Value: "2022-01-01"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vNumber", Value: "2423423"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vText", Value: "lorm ipsum "})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vBoolean_T", Value: "1"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vBoolean_F", Value: ""})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vEnum", Value: "abc"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vGeometry", Value: `{"lat":1,"lng":1}`})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vJSON", Value: `[{"bool":true"}]`})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vBlob", Value: "0110101"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vUUID", Value: "ba485865-54f9-44de-bde8-6965556c022a"})

	bootstrap(t, func(ctx context.Context, t *testing.T, h helper, svc dalService) {
		h.a.NoError(svc.ReplaceModel(ctx, model))
		h.a.NoError(svc.Create(ctx, model.ToFilter(), dal.CreateOperations(), &rIn))

		h.a.NoError(svc.Lookup(ctx, model.ToFilter(), dal.LookupOperations(), dal.PKValues{"id": rIn.ID}, &rOut))

		for _, inVal := range rIn.Values {
			outVal := rOut.Values.Get(inVal.Name, 0)
			h.a.Equal(inVal.Value, outVal.Value, "field", inVal.Name)
		}
	})
}

func benchmark_dal_codec_plain(b *testing.B, count int) {
	b.Skip("needs refactoring")

	model := &dal.Model{
		Ident: "compose_record_partitioned",
		Attributes: dal.AttributeSet{
			&dal.Attribute{Ident: "ID", Type: &dal.TypeID{}, Store: &dal.CodecAlias{Ident: "id"}, PrimaryKey: true},

			&dal.Attribute{Ident: "vID", Type: &dal.TypeID{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vRef", Type: &dal.TypeRef{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vTimestamp", Type: &dal.TypeTimestamp{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vTime", Type: &dal.TypeTime{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vDate", Type: &dal.TypeDate{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vNumber", Type: &dal.TypeNumber{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vText", Type: &dal.TypeText{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vBoolean_T", Type: &dal.TypeBoolean{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vBoolean_F", Type: &dal.TypeBoolean{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vEnum", Type: &dal.TypeEnum{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vGeometry", Type: &dal.TypeGeometry{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vJSON", Type: &dal.TypeJSON{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vBlob", Type: &dal.TypeBlob{}, Store: &dal.CodecPlain{}},
			&dal.Attribute{Ident: "vUUID", Type: &dal.TypeUUID{}, Store: &dal.CodecPlain{}},
		},
	}

	var (
		rIn = types.Record{ID: 42}
	)

	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vID", Value: "34324"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vRef", Value: "32423"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vTimestamp", Value: "2022-01-01T10:10:10"})
	rIn.Values = rIn.Values.Set(&types.RecordValue{Name: "vTime", Value: "04:10:10"})
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

	bootstrapBenchmark(b, func(ctx context.Context, b *testing.B, h helper, svc dalService) {
		h.a.NoError(svc.ReplaceModel(ctx, model))

		ctr := uint64(1)
		for rep := 0; rep < b.N; rep++ {
			insert := make([]dal.ValueGetter, 0, 100)
			for i := 0; i < 100; i++ {
				tmp := rIn
				tmp.ID = ctr
				ctr++

				insert = append(insert, &tmp)
			}

			b.StartTimer()
			h.a.NoError(svc.Create(ctx, model.ToFilter(), dal.CreateOperations(), insert...))
			b.StopTimer()
		}
	})
}

func Benchmark_dal_codec_plain_1(b *testing.B)   { benchmark_dal_codec_plain(b, 1) }
func Benchmark_dal_codec_plain_10(b *testing.B)  { benchmark_dal_codec_plain(b, 10) }
func Benchmark_dal_codec_plain_100(b *testing.B) { benchmark_dal_codec_plain(b, 100) }

func Benchmark_dal_codec_plain_1000(b *testing.B)  { benchmark_dal_codec_plain(b, 1000) }
func Benchmark_dal_codec_plain_10000(b *testing.B) { benchmark_dal_codec_plain(b, 10000) }

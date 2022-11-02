package dal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDiff_same(t *testing.T) {
	a := &Model{
		Attributes: AttributeSet{{
			Ident: "F1",
			Type:  TypeText{},
			Store: &CodecPlain{},
		}},
	}

	dd := a.Diff(a)
	require.Len(t, dd, 0)
}

func TestDiff_wrongAttrType(t *testing.T) {
	a := &Model{
		Attributes: AttributeSet{{
			Ident: "F1",
			Type:  TypeText{},
			Store: &CodecPlain{},
		}},
	}
	b := &Model{
		Attributes: AttributeSet{{
			Ident: "F1",
			Type:  TypeBlob{},
			Store: &CodecPlain{},
		}},
	}

	dd := a.Diff(b)
	require.Len(t, dd, 1)
	require.Equal(t, AttributeTypeMissmatch, dd[0].Type)
	require.Equal(t, AttributeChanged, dd[0].Modification)
}

func TestDiff_removedAttr(t *testing.T) {
	a := &Model{
		Attributes: AttributeSet{{
			Ident: "F1",
			Type:  TypeText{},
			Store: &CodecPlain{},
		}, {
			Ident: "F2",
			Type:  TypeText{},
			Store: &CodecPlain{},
		}},
	}
	b := &Model{
		Attributes: AttributeSet{{
			Ident: "F1",
			Type:  TypeText{},
			Store: &CodecPlain{},
		}},
	}

	dd := a.Diff(b)
	require.Len(t, dd, 1)
	require.Equal(t, AttributeMissing, dd[0].Type)
	require.Equal(t, AttributeDeleted, dd[0].Modification)
	require.NotNil(t, dd[0].Original)
	require.Nil(t, dd[0].Inserted)
}

func TestDiff_addedAttr(t *testing.T) {
	a := &Model{
		Attributes: AttributeSet{{
			Ident: "F1",
			Type:  TypeText{},
			Store: &CodecPlain{},
		}},
	}
	b := &Model{
		Attributes: AttributeSet{{
			Ident: "F1",
			Type:  TypeText{},
			Store: &CodecPlain{},
		}, {
			Ident: "F2",
			Type:  TypeText{},
			Store: &CodecPlain{},
		}},
	}

	dd := a.Diff(b)
	require.Len(t, dd, 1)
	require.Equal(t, AttributeMissing, dd[0].Type)
	require.Equal(t, AttributeAdded, dd[0].Modification)
	require.Nil(t, dd[0].Original)
	require.NotNil(t, dd[0].Inserted)
}

func TestDiff_changedCodec(t *testing.T) {
	a := &Model{
		Attributes: AttributeSet{{
			Ident: "F1",
			Type:  TypeText{},
			Store: &CodecPlain{},
		}},
	}
	b := &Model{
		Attributes: AttributeSet{{
			Ident: "F1",
			Type:  TypeText{},
			Store: &CodecRecordValueSetJSON{},
		}},
	}

	dd := a.Diff(b)
	require.Len(t, dd, 1)
	require.Equal(t, AttributeCodecMismatch, dd[0].Type)
	require.Equal(t, AttributeChanged, dd[0].Modification)
}

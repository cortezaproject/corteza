package id

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	TestRes struct {
		Id ID
	}
)

func TestNumericIdentifiers(t *testing.T) {
	ints := []uint64{
		1,
		100,
		231,
		1000,
		1234567890,
		18446744073709551615, // max uint64
	}

	for _, i := range ints {
		id, err := NumID(i)
		require.NoError(t, err)
		require.Equal(t, i, id.Number())
		require.True(t, id.isNum())
		require.False(t, id.isStr())

		tmpRes := fmt.Sprintf(`{"Id":"%d"}`, i)
		var tr TestRes
		err = json.Unmarshal([]byte(tmpRes), &tr)
		require.NoError(t, err)
		require.Equal(t, i, tr.Id.Number())
		require.True(t, tr.Id.isNum())
		require.False(t, tr.Id.isStr())

		bb, err := json.Marshal(&tr)
		require.NoError(t, err)
		require.Equal(t, tmpRes, string(bb))
	}
}

func TestStringIdentifiers(t *testing.T) {
	strs := []string{
		"-1",
		"asdf",
		"019864d2-3d93-7a86-9012-ec15f6358ab0", // uuid v7
		"0e04ae61-6887-4dda-bcef-c514240d1cd3", // uuid v4
		"b0a372c4-6eb4-11f0-8de9-0242ac120002", // uuid v1
		"00000000-0000-0000-0000-000000000000", // nil uuid
	}

	for _, s := range strs {
		id, err := StrID(s)
		require.NoError(t, err)
		require.Equal(t, s, id.String())
		require.False(t, id.isNum())
		require.True(t, id.isStr())

		tmpRes := fmt.Sprintf(`{"Id":"%s"}`, s)
		var tr TestRes
		err = json.Unmarshal([]byte(tmpRes), &tr)
		require.NoError(t, err)
		require.Equal(t, s, tr.Id.String())
		require.False(t, tr.Id.isNum())
		require.True(t, tr.Id.isStr())

		bb, err := json.Marshal(&tr)
		require.NoError(t, err)
		require.Equal(t, tmpRes, string(bb))
	}
}

func TestByteIdentifiers(t *testing.T) {
	tcc := []struct {
		in []byte

		outNum uint64
		outStr string
	}{
		{
			in:     []byte("123"),
			outNum: 123,
			outStr: "",
		}, {
			in:     []byte("asdf"),
			outNum: 0,
			outStr: "asdf",
		}, {
			in:     []byte("019864d2-3d93-7a86-9012-ec15f6358ab0"),
			outNum: 0,
			outStr: "019864d2-3d93-7a86-9012-ec15f6358ab0",
		},
	}

	for _, tc := range tcc {
		t.Run(string(tc.in), func(t *testing.T) {
			id, err := ByteID(tc.in)
			require.NoError(t, err)

			require.Equal(t, tc.outNum, id.Number())
			require.Equal(t, tc.outStr, id.String())
		})
	}
}

func TestInSlice(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		require.False(t, InSlice(MustNumID(1)))
	})

	t.Run("yes", func(t *testing.T) {
		require.True(t, InSlice(MustNumID(2), MustNumID(1), MustNumID(2), MustNumID(3)))
	})

	t.Run("no", func(t *testing.T) {
		require.False(t, InSlice(MustNumID(99), MustNumID(1), MustNumID(2), MustNumID(3)))
	})
}

func TestRemoveFromSlice(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		require.Len(t, RemoveFromSlice(MustNumID(1)), 0)
	})

	t.Run("found", func(t *testing.T) {
		slc := []ID{MustNumID(1), MustNumID(2)}
		slc = RemoveFromSlice(MustNumID(1), slc...)

		require.Len(t, slc, 1)
		require.Equal(t, MustNumID(2), slc[0])
	})

	t.Run("not found", func(t *testing.T) {
		slc := []ID{MustNumID(1), MustNumID(2)}
		slc = RemoveFromSlice(MustNumID(999), slc...)

		require.Len(t, slc, 2)
	})
}

// func TestIdentifierAssignment(t *testing.T) {
// 	ints := []uint64{
// 		1,
// 		100,
// 		231,
// 		1000,
// 		1234567890,
// 		18446744073709551615, // max uint64
// 	}

// 	strs := []string{
// 		"1",
// 		"asdf",
// 		"019864d2-3d93-7a86-9012-ec15f6358ab0", // uuid v7
// 		"0e04ae61-6887-4dda-bcef-c514240d1cd3", // uuid v4
// 		"b0a372c4-6eb4-11f0-8de9-0242ac120002", // uuid v1
// 		"00000000-0000-0000-0000-000000000000", // nil uuid
// 	}

// }

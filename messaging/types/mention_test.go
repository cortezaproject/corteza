package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMentionSet_Diff(t *testing.T) {
	ex := MentionSet{&Mention{ID: 1000}, &Mention{ID: 1001}}
	add, upd, del := ex.Diff(MentionSet{&Mention{ID: 1001}, &Mention{UserID: 1}})

	require.True(t, len(add) == 1 && len(add.FindByUserID(1)) == 1, "Did not find expected mention (UserID:1) for creation")
	require.True(t, len(upd) == 1 && upd.FindByID(1001) != nil, "Did not find expected mention (id:1001) for update")
	require.True(t, len(del) == 1 && del.FindByID(1000) != nil, "Did not find expected mention (id:1000) for removal")
}

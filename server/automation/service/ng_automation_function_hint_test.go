package service

import (
	"testing"

	"github.com/crusttech/human/server/automation/types"
	"github.com/stretchr/testify/require"
)

// The library as it stands: singular/plural pairs and a couple of short names,
// which is what makes a near match worth computing at all.
var hintLibrary = []types.ConstructFunction{
	{Ref: "composeRecordsCreate"},
	{Ref: "composeRecordsUpdate"},
	{Ref: "composeRecordsDelete"},
	{Ref: "composeRecordsEach"},
	{Ref: "loopDo"},
	{Ref: "notificationSend"},
	{Ref: "usersCreate"},
}

func TestUnknownFunctionHintNamesTheNearMatch(t *testing.T) {
	// The observed failure: a TAQ stored with ref "composeRecordCreate" and an
	// issue that named the bad ref and none of the legal ones.
	h := unknownFunctionHint(hintLibrary, "composeRecordCreate")

	require.Contains(t, h, `"composeRecordsCreate"`)
	require.Contains(t, h, "automation_taq_construct_lookup")
	require.NotContains(t, h, `"loopDo"`)
}

func TestUnknownFunctionHintPointsAtTheCatalogueWhenNothingIsClose(t *testing.T) {
	h := unknownFunctionHint(hintLibrary, "sendCarrierPigeon")

	require.NotContains(t, h, "did you mean")
	require.Contains(t, h, "automation_taq_construct_lookup")
}

func TestNearestFunctionRefsDoesNotMatchAShortUnrelatedRef(t *testing.T) {
	// "loopDo" is six characters; a one-edit budget must not drag it in as the
	// near match for anything else short.
	require.Empty(t, nearestFunctionRefs(hintLibrary, "gateway"))
}

func TestNearestFunctionRefsCapsAtThree(t *testing.T) {
	got := nearestFunctionRefs(hintLibrary, "composeRecords")
	require.Len(t, got, 3)
}

// The correction has to lead: a caller reads the first name offered, and the
// whole point is that a slipped plural is one word away.
func TestNearestFunctionRefsRanksTheClosestFirst(t *testing.T) {
	got := nearestFunctionRefs(hintLibrary, "composeRecordCreate")
	require.NotEmpty(t, got)
	require.Equal(t, "composeRecordsCreate", got[0])
}

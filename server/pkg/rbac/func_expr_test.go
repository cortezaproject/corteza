package rbac

import (
	"testing"

	"github.com/crusttech/human/server/pkg/id"
	"github.com/stretchr/testify/require"
)

func TestIsDescendantOf(t *testing.T) {
	t.Run("single values", func(t *testing.T) {
		var (
			ownerL []id.ID
			userL  []id.ID
			pathsL [][]string
		)

		isAboveChecker = func(owner, user id.ID, paths ...string) bool {
			ownerL = append(ownerL, owner)
			userL = append(userL, user)
			pathsL = append(pathsL, paths)
			return true
		}

		gRBAC = &service{}

		isDescendantOf(uint64(101), uint64(102))

		require.Len(t, ownerL, 1)
		require.Len(t, userL, 1)

		require.Equal(t, id.MustNumID(102), ownerL[0])
		require.Equal(t, id.MustNumID(101), userL[0])
	})

	t.Run("multiple owners", func(t *testing.T) {
		var (
			ownerL []id.ID
			userL  []id.ID
			pathsL [][]string
		)

		isAboveChecker = func(owner, user id.ID, paths ...string) bool {
			ownerL = append(ownerL, owner)
			userL = append(userL, user)
			pathsL = append(pathsL, paths)
			return false
		}

		gRBAC = &service{}

		isDescendantOf(uint64(101), []uint64{uint64(102), uint64(103)})

		require.Len(t, ownerL, 2)
		require.Len(t, userL, 2)

		require.Equal(t, id.MustNumID(102), ownerL[0])
		require.Equal(t, id.MustNumID(103), ownerL[1])
		require.Equal(t, id.MustNumID(101), userL[0])
		require.Equal(t, id.MustNumID(101), userL[1])
	})
}

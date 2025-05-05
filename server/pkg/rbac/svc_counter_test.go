package rbac

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWrapperCounter(t *testing.T) {
	// @note since I'm leaving the decayInterval empty we don't need to fiddle
	// with lastAccess timestamps
	svc := &usageCounter[string]{
		index: map[string]counterItem[string]{},

		sigEvictThreshold: 0.5,
		decayFactor:       0.5,
	}

	svc.inc("k1")
	aux := svc.index["k1"]
	require.Equal(t, 1.0, aux.score)

	svc.inc("k2")
	aux = svc.index["k1"]
	require.Equal(t, 1.0, aux.score)
	aux = svc.index["k2"]
	require.Equal(t, 1.0, aux.score)

	svc.inc("k1")
	aux = svc.index["k1"]
	require.Equal(t, 2.0, aux.score)
	aux = svc.index["k2"]
	require.Equal(t, 1.0, aux.score)

	svc.decay()
	aux = svc.index["k1"]
	require.Equal(t, 1.0, aux.score)
	aux = svc.index["k2"]
	require.Equal(t, 0.5, aux.score)

	cleaned := svc.evict()
	require.Len(t, cleaned, 1)
	aux, ok := svc.index["k1"]
	require.True(t, ok)

	aux, ok = svc.index["k2"]
	require.False(t, ok)

	svc.decay()
	aux = svc.index["k1"]
	require.Equal(t, 0.5, aux.score)

	cleaned = svc.evict()
	require.Len(t, cleaned, 1)
	aux, ok = svc.index["k1"]
	require.False(t, ok)
}

func TestCounterCleanRoleKeys(t *testing.T) {
	req := require.New(t)
	svc := &usageCounter[string]{
		index: map[string]counterItem[string]{},

		sigEvictThreshold: 0.5,
		decayFactor:       0.5,

		checkKeyInclusion: func(k string, role uint64) bool {
			return strings.HasPrefix(k, fmt.Sprintf("%d", role))
		},
	}

	svc.inc("12:res/1/2/3")
	svc.inc("12:res/2/2/3")
	svc.inc("12:res/3/2/3")
	svc.inc("13:res/1/2/3")
	svc.inc("14:res/1/2/3")

	svc.cleanRoleKeys(12)
	req.Len(svc.index, 2)

	svc.cleanRoleKeys(13)
	req.Len(svc.index, 1)

	svc.cleanRoleKeys(14)
	req.Len(svc.index, 0)
}

func TestCounterBestPerformers(t *testing.T) {
	req := require.New(t)
	svc := &usageCounter[string]{
		index: map[string]counterItem[string]{},

		sigEvictThreshold: 0.5,
		decayFactor:       0.5,
	}

	svc.inc("12:res/1/2/3")
	svc.inc("12:res/2/2/3")
	svc.inc("12:res/3/2/3")
	svc.inc("13:res/1/2/3")
	svc.inc("14:res/1/2/3")

	// -1 gets all
	out := svc.bestPerformers(-1)
	req.Len(out, 5)

	// 0 gets none
	out = svc.bestPerformers(0)
	req.Len(out, 0)

	// n gets some
	out = svc.bestPerformers(2)
	req.Len(out, 2)

	// too big n gets max
	out = svc.bestPerformers(99)
	req.Len(out, 5)
}

func TestMinHeap(t *testing.T) {
	hp := MinHeap[string]{}

	hp = append(hp, counterItem[string]{score: 4, key: "4"})
	hp = append(hp, counterItem[string]{score: 10, key: "10"})
	hp = append(hp, counterItem[string]{score: 1, key: "1"})
	hp = append(hp, counterItem[string]{score: 4, key: "4"})
	hp = append(hp, counterItem[string]{score: 2, key: "2"})
	hp = append(hp, counterItem[string]{score: 99, key: "99"})
	hp = append(hp, counterItem[string]{score: 12, key: "12"})
	hp = append(hp, counterItem[string]{score: 3, key: "3"})
	hp = append(hp, counterItem[string]{score: 4, key: "4"})
	hp = append(hp, counterItem[string]{score: 5, key: "5"})

	sort.Sort(hp)

	require.Equal(t, "1", hp[0].key)
	require.Equal(t, "2", hp[1].key)
	require.Equal(t, "3", hp[2].key)
	require.Equal(t, "4", hp[3].key)
	require.Equal(t, "4", hp[4].key)
	require.Equal(t, "4", hp[5].key)
	require.Equal(t, "5", hp[6].key)
	require.Equal(t, "10", hp[7].key)
	require.Equal(t, "12", hp[8].key)
	require.Equal(t, "99", hp[9].key)
}

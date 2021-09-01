package cache

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewLRUCache(t *testing.T) {
	req := require.New(t)
	lc := NewLRUCache(2)

	t.Run("Add to cache", func(t *testing.T) {
		b := lc.Add("k1", "V1")
		req.Equal(true, b)
	})

	t.Run("Add to cache but same key", func(t *testing.T) {
		b := lc.Add("k1", "V3")
		req.Equal(false, b)
	})

	t.Run("Get value from cache", func(t *testing.T) {
		r, b := lc.Get("k1")
		req.Equal(true, b)
		req.Equal("V3", r)
		rn, bf := lc.Get("k5")
		req.Equal(false, bf)
		req.Equal("", rn)
	})

	t.Run("Delete value from cache", func(t *testing.T) {
		b := lc.Delete("k1")
		req.Equal(true, b)
		bf := lc.Delete("k5")
		req.Equal(false, bf)
	})

	t.Run("Remove last element", func(t *testing.T) {
		_ = lc.Add("k2", "V2")
		_ = lc.Add("k4", "V4")
		k3 := lc.Add("k3", "V3")
		req.Equal(true, k3)
	})
}

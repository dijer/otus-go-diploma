package cache

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := New(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := New(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := New(3)
		c.Set("a", 1)
		c.Set("b", 2)
		c.Set("c", 3)
		c.Set("d", 4)
		val, ok := c.Get("a")
		require.False(t, ok)
		require.Equal(t, nil, val)
	})

	t.Run("purge logic old elements", func(t *testing.T) {
		c := New(3)
		c.Set("a", 1)
		b1ok, has := c.Get("b")
		require.False(t, has)
		require.Nil(t, b1ok)
		c.Set("b", 2)
		b1ok, has = c.Get("b")
		require.True(t, has)
		require.Equal(t, b1ok, 2)
		c.Set("c", 3)
		c.Set("d", 4)
		has = c.Set("b", -3)
		require.True(t, has)

		b1ok, has = c.Get("b")
		require.True(t, has)
		require.Equal(t, b1ok, -3)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := New(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
			if err != nil {
				t.Errorf("Err: %v", err)
				return
			}
			c.Get(Key(strconv.Itoa(int(n.Int64()))))
		}
	}()

	wg.Wait()
}

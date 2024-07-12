package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

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
		// Write me
		c := NewCache(3)

		// на логику выталкивания элементов из-за размера очереди
		// n = 3, добавили 4 элемента
		_ = c.Set("first", 100)
		_ = c.Set("second", 200)
		_ = c.Set("third", 300)
		_ = c.Set("fourth", 400)

		// 1й из кэша вытолкнулся
		val, ok := c.Get("first")
		require.False(t, ok)
		require.Nil(t, val)

		// на логику выталкивания давно используемых элементов
		// обратились несколько раз к разным элементам: изменили значение, получили значение
		_, _ = c.Get("second")
		_, _ = c.Get("fourth")
		_ = c.Set("third", 3000)
		// добавили 4й элемент
		_ = c.Set("fifth", 50)

		// из первой тройки вытолкнется тот элемент, что был затронут наиболее давно
		val, ok = c.Get("second")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
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
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}

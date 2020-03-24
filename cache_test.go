package easycache

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

/**
 * @Author: eggsy
 * @Description:
 * @File:  cache_test
 * @Version: 1.0.0
 * @Date: 2020-03-24 17:11
 */

func TestEasyCache_Set(t *testing.T) {
	cache := New(time.Duration(5 * time.Second))
	cache.Set("hello", "15s", time.Duration(15*time.Second))
	value, found := cache.Get("hello")
	assert.True(t, found)
	assert.True(t, value == "15s")
}

func TestEasyCache_Delete(t *testing.T) {
	cache := New(time.Duration(5 * time.Second))
	cache.Set("hello", "15s", time.Duration(15*time.Second))
	cache.Delete("hello")
	_, found := cache.Get("hello")
	assert.True(t, !found)
}

// 测试过期
func TestEasyCache_Expires(t *testing.T) {
	cache := New(time.Duration(5 * time.Second))
	cache.Set("hello", "8s", time.Duration(8*time.Second))
	value, found := cache.Get("hello")
	assert.True(t, found)
	assert.True(t, value == "8s")
	time.Sleep(time.Second * 15)
	_, found = cache.Get("hello")
	assert.True(t, !found)
}

func BenchmarkNew(b *testing.B) {
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			New(5 * time.Second).Close()
		}
	})
}

func BenchmarkGet(b *testing.B) {
	c := New(5 * time.Second)
	defer c.Close()
	c.Set("Hello", "World", 0)

	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Get("Hello")
		}
	})
}

func BenchmarkSet(b *testing.B) {
	c := New(5 * time.Second)
	defer c.Close()

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Set("Hello", "World", 0)
		}
	})
}

func BenchmarkDelete(b *testing.B) {
	c := New(5 * time.Second)
	defer c.Close()

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			c.Delete("Hello")
		}
	})
}

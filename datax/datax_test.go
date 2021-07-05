package datax

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSet(t *testing.T) {
	s := NewSet()
	assert.False(t, s.Has("test_key"), "initial set should hot have any key")
	s.Add("test_key")
	assert.True(t, s.Has("test_key"), "set should have key when set")
	assert.False(t, s.TryAdd("test_key"), "try set should false when set a existed key")
	assert.True(t, s.TryAdd("test_key2"), "try set should true when set a non-existed key")
	assert.True(t, s.Has("test_key2"), "set should have key when try set success")
	s.Remove("test_key2")
	assert.False(t, s.Has("test_key2"), "try set should false after delete key")
}

func BenchmarkSetWithPtr(b *testing.B) {
	s := NewSet()
	k := &Empty{}
	for i := 0; i < b.N; i++ {
		s.Add(k)
		s.Has(k)
		s.Remove(k)
	}
}

func BenchmarkMapWithPtr(b *testing.B) {
	s := map[*Empty]Empty{}
	k := &Empty{}
	for i := 0; i < b.N; i++ {
		s[k] = Empty{}
		_, _ = s[k]
		delete(s, k)
	}
}

func BenchmarkSetWithString(b *testing.B) {
	s := NewSet()
	for i := 0; i < b.N; i++ {
		s.Add("test_key")
		s.Has("test_key")
		s.Remove("test_key")
	}
}

func BenchmarkMapWithString(b *testing.B) {
	s := map[string]Empty{}
	for i := 0; i < b.N; i++ {
		s["test_key"] = Empty{}
		_, _ = s["test_key"]
		delete(s, "test_key")
	}
}

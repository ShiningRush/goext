package datax

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSet(t *testing.T) {
	s := Set{}
	assert.False(t, s.Has("test_key"), "initial set should hot have any key")
	s.Set("test_key")
	assert.True(t, s.Has("test_key"), "set should have key when set")
	assert.False(t, s.TrySet("test_key"), "try set should false when set a existed key")
	assert.True(t, s.TrySet("test_key2"), "try set should true when set a non-existed key")
	assert.True(t, s.Has("test_key2"), "set should have key when try set success")
	s.Delete("test_key2")
	assert.False(t, s.Has("test_key2"), "try set should false after delete key")
}

func BenchmarkSetWithPtr(b *testing.B) {
	s := Set{}
	k := &Empty{}
	for i := 0; i < b.N; i++ {
		s.Set(k)
		s.Has(k)
		s.Delete(k)
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
	s := Set{}
	for i := 0; i < b.N; i++ {
		s.Set("test_key")
		s.Has("test_key")
		s.Delete("test_key")
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

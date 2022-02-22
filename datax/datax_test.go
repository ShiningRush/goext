package datax

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	s := NewSet()
	assert.False(t, s.Has("test_key"), "initial set should hot have any key")
	s.Add("test_key", "test_key1")
	assert.True(t, s.Has("test_key"), "set should have key when set")
	assert.True(t, s.Has("test_key1"), "set should have key when set")
	assert.False(t, s.TryAdd("test_key"), "try set should false when set a existed key")
	assert.True(t, s.TryAdd("test_key2"), "try set should true when set a non-existed key")
	assert.True(t, s.Has("test_key2"), "set should have key when try set success")
	assert.Equal(t, 3, s.Len(), "set should have 2 items")
	assert.ElementsMatch(t, []interface{}{"test_key", "test_key1", "test_key2"}, s.All())
	s.Remove("test_key1", "test_key2")
	assert.False(t, s.Has("test_key2"), "set should return false after delete key")
	assert.Equal(t, 1, s.Len(), "set should have 1 items")
}

func TestNewSetOperation(t *testing.T) {
	a := NewSet().Add("1", "2", "3")
	b := NewSet().Add("0", "1", "2", "3", "4")
	c := NewSet().Add("2", "3", "4")
	empty := NewSet()

	assert.True(t, a.IsSubsetOf(a))
	assert.True(t, a.IsSubsetOf(b))
	assert.False(t, a.IsProperSubsetOf(a))
	assert.True(t, a.IsProperSubsetOf(b))
	assert.True(t, b.IsSupersetOf(a))

	assert.True(t, NewSet().Add("2", "3").Equal(a.Intersect(c)))
	assert.True(t, NewSet().Equal(a.Intersect(empty)))

	assert.True(t, NewSet().Add("1", "2", "3", "4").Equal(a.Union(c)))
	assert.True(t, NewSet().Add("1", "2", "3").Equal(a.Union(empty)))

	assert.True(t, NewSet().Add("1").Equal(a.Diff(c)))
	assert.True(t, NewSet().Add("4").Equal(c.Diff(a)))
	assert.True(t, NewSet().Add("2", "3", "4").Equal(c.Diff(empty)))
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

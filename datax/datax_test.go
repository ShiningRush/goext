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

	// diff type test
	var i uint = 10
	s.AddInt32(10)
	assert.False(t, s.Has(i), "set should not query when type diff")
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

func TestHasItem(t *testing.T) {
	tests := []struct {
		caseDesc  string
		giveArray []string
		giveItem  string
		wantRet   bool
	}{
		{
			caseDesc:  "existed",
			giveArray: []string{"fine", "bad"},
			giveItem:  "bad",
			wantRet:   true,
		},
		{
			caseDesc:  "not existed",
			giveArray: []string{"fine", "bad"},
			giveItem:  "not",
			wantRet:   false,
		},
		{
			caseDesc:  "nil",
			giveArray: nil,
			giveItem:  "not",
			wantRet:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.caseDesc, func(t *testing.T) {
			ret := HasItem(tc.giveArray, tc.giveItem)
			assert.Equal(t, tc.wantRet, ret)
		})
	}
}

func TestIsSuperset(t *testing.T) {
	tests := []struct {
		caseDesc string
		giveSrc  []string
		giveDest []string
		wantRet  bool
	}{
		{
			caseDesc: "is",
			giveSrc:  []string{"a", "b", "c"},
			giveDest: []string{"a", "c"},
			wantRet:  true,
		},
		{
			caseDesc: "not is",
			giveSrc:  []string{"a", "b", "c"},
			giveDest: []string{"d", "c"},
			wantRet:  false,
		},
		{
			caseDesc: "empty",
			giveSrc:  nil,
			giveDest: nil,
			wantRet:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.caseDesc, func(t *testing.T) {
			ret := IsSuperset(tc.giveSrc, tc.giveDest)
			assert.Equal(t, tc.wantRet, ret)
		})
	}
}

func TestIsSubset(t *testing.T) {
	tests := []struct {
		caseDesc string
		giveSrc  []string
		giveDest []string
		wantRet  bool
	}{
		{
			caseDesc: "is",
			giveSrc:  []string{"a", "c"},
			giveDest: []string{"a", "b", "c"},
			wantRet:  true,
		},
		{
			caseDesc: "self",
			giveSrc:  []string{"a", "b", "c"},
			giveDest: []string{"a", "b", "c"},
			wantRet:  true,
		},
		{
			caseDesc: "empty",
			giveSrc:  nil,
			giveDest: nil,
			wantRet:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.caseDesc, func(t *testing.T) {
			ret := IsSubset(tc.giveSrc, tc.giveDest)
			assert.Equal(t, tc.wantRet, ret)
		})
	}
}

func TestIsProperSubset(t *testing.T) {
	tests := []struct {
		caseDesc string
		giveSrc  []string
		giveDest []string
		wantRet  bool
	}{
		{
			caseDesc: "is",
			giveSrc:  []string{"a", "c"},
			giveDest: []string{"a", "b", "c"},
			wantRet:  true,
		},
		{
			caseDesc: "not is",
			giveSrc:  []string{"a", "b", "c"},
			giveDest: []string{"a", "b", "c"},
			wantRet:  false,
		},
		{
			caseDesc: "empty",
			giveSrc:  nil,
			giveDest: nil,
			wantRet:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.caseDesc, func(t *testing.T) {
			ret := IsProperSubset(tc.giveSrc, tc.giveDest)
			assert.Equal(t, tc.wantRet, ret)
		})
	}
}

func TestIntersect(t *testing.T) {
	tests := []struct {
		caseDesc string
		giveA    []string
		giveB    []string
		wantRet  []string
	}{
		{
			caseDesc: "sanity",
			giveA:    []string{"a", "b", "c"},
			giveB:    []string{"a", "c", "e"},
			wantRet:  []string{"a", "c"},
		},
		{
			caseDesc: "a-empty",
			giveA:    nil,
			giveB:    []string{"a", "c"},
			wantRet:  nil,
		},
		{
			caseDesc: "b-empty",
			giveA:    []string{"a", "b", "c"},
			giveB:    nil,
			wantRet:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.caseDesc, func(t *testing.T) {
			ret := Intersect(tc.giveA, tc.giveB)
			assert.ElementsMatch(t, tc.wantRet, ret)
		})
	}
}

func TestDiff(t *testing.T) {
	tests := []struct {
		caseDesc string
		giveA    []string
		giveB    []string
		wantRet  []string
	}{
		{
			caseDesc: "sanity",
			giveA:    []string{"a", "b", "c"},
			giveB:    []string{"a", "c", "e"},
			wantRet:  []string{"b"},
		},
		{
			caseDesc: "a-empty",
			giveA:    nil,
			giveB:    []string{"a", "c"},
			wantRet:  nil,
		},
		{
			caseDesc: "b-empty",
			giveA:    []string{"a", "b", "c"},
			giveB:    nil,
			wantRet:  []string{"a", "b", "c"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.caseDesc, func(t *testing.T) {
			ret := Diff(tc.giveA, tc.giveB)
			assert.ElementsMatch(t, tc.wantRet, ret)
		})
	}
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

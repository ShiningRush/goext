package datax

type LoopSetFunc func(item interface{}) (breakLoop bool)

type Empty struct{}

func NewSet() *Set {
	return &Set{
		idx: map[interface{}]Empty{},
	}
}

func NewSetFrom[T any](items ...T) *Set {
	s := NewSet()
	for _, v := range items {
		s.idx[v] = Empty{}
	}
	return s
}

type Set struct {
	idx map[interface{}]Empty
}

func (s *Set) Len() int {
	return len(s.idx)
}

func (s *Set) Loop(f LoopSetFunc) {
	for item := range s.idx {
		f(item)
	}
}

func (s *Set) All() (items []any) {
	s.Loop(func(item any) (breakLoop bool) {
		items = append(items, item)
		return
	})
	return
}

func (s *Set) Has(item any) bool {
	_, ok := s.idx[item]
	return ok
}
func (s *Set) Add(items ...any) *Set {
	for _, v := range items {
		s.idx[v] = Empty{}
	}
	return s
}
func (s *Set) AddString(items ...string) *Set {
	for _, v := range items {
		s.idx[v] = Empty{}
	}
	return s
}
func (s *Set) AddInt(items ...int) *Set {
	for _, v := range items {
		s.idx[v] = Empty{}
	}
	return s
}
func (s *Set) AddInt32(items ...int32) *Set {
	for _, v := range items {
		s.idx[v] = Empty{}
	}
	return s
}
func (s *Set) AddInt64(items ...int64) *Set {
	for _, v := range items {
		s.idx[v] = Empty{}
	}
	return s
}
func (s *Set) AddUint(items ...uint) *Set {
	for _, v := range items {
		s.idx[v] = Empty{}
	}
	return s
}
func (s *Set) AddUint32(items ...uint32) *Set {
	for _, v := range items {
		s.idx[v] = Empty{}
	}
	return s
}
func (s *Set) AddUint64(items ...uint64) *Set {
	for _, v := range items {
		s.idx[v] = Empty{}
	}
	return s
}

func (s *Set) Remove(keys ...interface{}) *Set {
	for _, v := range keys {
		delete(s.idx, v)
	}
	return s
}

func (s *Set) TryAdd(key any) bool {
	if s.Has(key) {
		return false
	}
	s.idx[key] = Empty{}
	return true
}

// Equal test whether s equal to another set
func (s *Set) Equal(another *Set) bool {
	if s.Len() != another.Len() {
		return false
	}

	return s.IsSubsetOf(another)
}

// IsSupersetOf test whether s is superset of another set
func (s *Set) IsSupersetOf(another *Set) bool {
	return another.IsProperSubsetOf(s)
}

// IsSubsetOf test whether s is a subset of another set
func (s *Set) IsSubsetOf(another *Set) bool {
	if s.Len() > another.Len() {
		return false
	}

	findDiff := false
	s.Loop(func(key interface{}) (breakLoop bool) {
		if !another.Has(key) {
			findDiff = true
			breakLoop = true
		}
		return
	})
	return !findDiff
}

// IsProperSubsetOf test whether s is a proper subset of another set
func (s *Set) IsProperSubsetOf(another *Set) bool {
	return s.Len() < another.Len() && s.IsSubsetOf(another)
}

// Intersect find the intersection between s and another set
func (s *Set) Intersect(another *Set) *Set {
	intersection := NewSet()
	s.Loop(func(key interface{}) (breakLoop bool) {
		if another.Has(key) {
			intersection.Add(key)
		}
		return
	})

	return intersection
}

// Union find the union between s and another set
func (s *Set) Union(another *Set) *Set {
	union := NewSet()
	s.Loop(func(key interface{}) (breakLoop bool) {
		union.Add(key)
		return
	})

	another.Loop(func(key interface{}) (breakLoop bool) {
		union.Add(key)
		return
	})

	return union
}

// Diff find the difference between s and another set
func (s *Set) Diff(another *Set) *Set {
	diffSet := NewSet()
	s.Loop(func(item interface{}) (breakLoop bool) {
		if !another.Has(item) {
			diffSet.Add(item)
		}
		return
	})
	return diffSet
}

// HasItem test if the "item" is in "lists"
func HasItem[T comparable](lists []T, item T) bool {
	for _, s := range lists {
		if s == item {
			return true
		}
	}
	return false
}

// IsSuperset test if the "src" is superset of "dest"
func IsSuperset[T any](src, dest []T) bool {
	srcSet := NewSetFrom(src...)
	descSet := NewSetFrom(dest...)
	return srcSet.IsSupersetOf(descSet)
}

// IsSubset test if the "src" is subset of "dest"
func IsSubset[T any](src, dest []T) bool {
	srcSet := NewSetFrom(src...)
	descSet := NewSetFrom(dest...)
	return srcSet.IsSubsetOf(descSet)
}

// IsProperSubset test if the "src" is proper subset of "dest"
func IsProperSubset[T any](src, dest []T) bool {
	srcSet := NewSetFrom(src...)
	descSet := NewSetFrom(dest...)
	return srcSet.IsProperSubsetOf(descSet)
}

// Intersect get intersection of two lists
func Intersect[T any](aArr, bArr []T) (intersection []T) {
	aSet := NewSetFrom(aArr...)
	bSet := NewSetFrom(bArr...)

	aSet.Intersect(bSet).Loop(func(item interface{}) (breakLoop bool) {
		intersection = append(intersection, item.(T))
		return
	})
	return
}

// Diff get differences of two strings
func Diff[T any](aArr, bArr []T) (diff []T) {
	aSet := NewSetFrom(aArr...)
	bSet := NewSetFrom(bArr...)

	aSet.Diff(bSet).Loop(func(item interface{}) (breakLoop bool) {
		diff = append(diff, item.(T))
		return
	})
	return
}

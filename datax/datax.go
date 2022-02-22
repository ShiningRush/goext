package datax

type LoopSetFunc func(item interface{}) (breakLoop bool)

type Empty struct{}

func NewSet() *Set {
	return &Set{
		idx: map[interface{}]Empty{},
	}
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

func (s *Set) All() (items []interface{}) {
	s.Loop(func(item interface{}) (breakLoop bool) {
		items = append(items, item)
		return
	})
	return
}

func (s *Set) Has(item interface{}) bool {
	_, ok := s.idx[item]
	return ok
}
func (s *Set) Add(items ...interface{}) *Set {
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

func (s *Set) Remove(keys ...interface{}) *Set {
	for _, v := range keys {
		delete(s.idx, v)
	}
	return s
}

func (s *Set) TryAdd(key interface{}) bool {
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

package datax

type Empty struct{}

func NewSet() *Set {
	return &Set{
		idx: map[interface{}]Empty{},
	}
}

type Set struct {
	idx map[interface{}]Empty
}

func (s *Set) Has(key interface{}) bool {
	_, ok := s.idx[key]
	return ok
}
func (s *Set) Add(key interface{}) {
	s.idx[key] = Empty{}
}

func (s *Set) TryAdd(key interface{}) bool {
	if s.Has(key) {
		return false
	}
	s.idx[key] = Empty{}
	return true
}

func (s *Set) Remove(key interface{}) {
	delete(s.idx, key)
}

func (s *Set) Len() int {
	return len(s.idx)
}

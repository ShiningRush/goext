package datax

type Empty struct{}

type Set map[interface{}]Empty

func (s Set) Has(key interface{}) bool {
	_, ok := s[key]
	return ok
}
func (s Set) Set(key interface{}) {
	s[key] = Empty{}
}

func (s Set) TrySet(key interface{}) bool {
	if s.Has(key) {
		return false
	}
	s[key] = Empty{}
	return true
}
func (s Set) Delete(key interface{}) {
	delete(s, key)
}

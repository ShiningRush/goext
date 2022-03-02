package intx

import "github.com/shiningrush/goext/datax"

// HasItem test if the "item" is in "ints"
func HasItem(ints []int, item int) bool {
	for _, s := range ints {
		if s == item {
			return true
		}
	}
	return false
}

// IsSuperset test if the "src" is superset of "dest"
func IsSuperset(src, dest []int) bool {
	srcSet := datax.NewSet().AddInt(src...)
	descSet := datax.NewSet().AddInt(dest...)
	return srcSet.IsSupersetOf(descSet)
}

// IsSubset test if the "src" is subset of "dest"
func IsSubset(src, dest []int) bool {
	srcSet := datax.NewSet().AddInt(src...)
	descSet := datax.NewSet().AddInt(dest...)
	return srcSet.IsSubsetOf(descSet)
}

// IsProperSubset test if the "src" is proper subset of "dest"
func IsProperSubset(src, dest []int) bool {
	srcSet := datax.NewSet().AddInt(src...)
	descSet := datax.NewSet().AddInt(dest...)
	return srcSet.IsProperSubsetOf(descSet)
}

// Intersect get intersection of two strings
func Intersect(aArr, bArr []int) (intersection []int) {
	aSet := datax.NewSet().AddInt(aArr...)
	bSet := datax.NewSet().AddInt(bArr...)

	aSet.Intersect(bSet).Loop(func(item interface{}) (breakLoop bool) {
		intersection = append(intersection, item.(int))
		return
	})
	return
}

// Diff get differences of two strings
func Diff(aArr, bArr []int) (diff []int) {
	aSet := datax.NewSet().AddInt(aArr...)
	bSet := datax.NewSet().AddInt(bArr...)

	aSet.Diff(bSet).Loop(func(item interface{}) (breakLoop bool) {
		diff = append(diff, item.(int))
		return
	})
	return
}

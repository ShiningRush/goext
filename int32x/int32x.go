package int32x

import "github.com/shiningrush/goext/datax"

// HasItem test if the "item" is in "ints"
func HasItem(ints []int32, item int32) bool {
	for _, s := range ints {
		if s == item {
			return true
		}
	}
	return false
}

// IsSuperset test if the "src" is superset of "dest"
func IsSuperset(src, dest []int32) bool {
	srcSet := datax.NewSet().AddInt32(src...)
	descSet := datax.NewSet().AddInt32(dest...)
	return srcSet.IsSupersetOf(descSet)
}

// IsSubset test if the "src" is subset of "dest"
func IsSubset(src, dest []int32) bool {
	srcSet := datax.NewSet().AddInt32(src...)
	descSet := datax.NewSet().AddInt32(dest...)
	return srcSet.IsSubsetOf(descSet)
}

// IsProperSubset test if the "src" is proper subset of "dest"
func IsProperSubset(src, dest []int32) bool {
	srcSet := datax.NewSet().AddInt32(src...)
	descSet := datax.NewSet().AddInt32(dest...)
	return srcSet.IsProperSubsetOf(descSet)
}

// Intersect get intersection of two strings
func Intersect(aArr, bArr []int32) (intersection []int32) {
	aSet := datax.NewSet().AddInt32(aArr...)
	bSet := datax.NewSet().AddInt32(bArr...)

	aSet.Intersect(bSet).Loop(func(item interface{}) (breakLoop bool) {
		intersection = append(intersection, item.(int32))
		return
	})
	return
}

// Diff get differences of two strings
func Diff(aArr, bArr []int32) (diff []int32) {
	aSet := datax.NewSet().AddInt32(aArr...)
	bSet := datax.NewSet().AddInt32(bArr...)

	aSet.Diff(bSet).Loop(func(item interface{}) (breakLoop bool) {
		diff = append(diff, item.(int32))
		return
	})
	return
}

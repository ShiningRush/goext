package int32x

import "github.com/shiningrush/goext/datax"

// HasItem test if the "item" is in "ints"
func HasItem(ints []int32, item int32) bool {
	return datax.HasItem(ints, item)
}

// IsSuperset test if the "src" is superset of "dest"
func IsSuperset(src, dest []int32) bool {
	return datax.IsSuperset(src, dest)
}

// IsSubset test if the "src" is subset of "dest"
func IsSubset(src, dest []int32) bool {
	return datax.IsSubset(src, dest)
}

// IsProperSubset test if the "src" is proper subset of "dest"
func IsProperSubset(src, dest []int32) bool {
	return datax.IsProperSubset(src, dest)
}

// Intersect get intersection of two strings
func Intersect(aArr, bArr []int32) (intersection []int32) {
	return datax.Intersect(aArr, bArr)
}

// Diff get differences of two strings
func Diff(aArr, bArr []int32) (diff []int32) {
	return datax.Diff(aArr, bArr)
}

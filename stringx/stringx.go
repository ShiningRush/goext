package stringx

import "github.com/shiningrush/goext/datax"

// HasItem test if the "item" is in "strs"
func HasItem(strs []string, item string) bool {
	return datax.HasItem(strs, item)
}

// IsSuperset test if the "src" is superset of "dest"
func IsSuperset(src, dest []string) bool {
	return datax.IsSuperset(src, dest)
}

// IsSubset test if the "src" is subset of "dest"
func IsSubset(src, dest []string) bool {
	return datax.IsSubset(src, dest)
}

// IsProperSubset test if the "src" is proper subset of "dest"
func IsProperSubset(src, dest []string) bool {
	return datax.IsProperSubset(src, dest)
}

// Intersect get intersection of two strings
func Intersect(aArr, bArr []string) (intersection []string) {
	return datax.Intersect(aArr, bArr)
}

// Diff get differences of two strings
func Diff(aArr, bArr []string) (diff []string) {
	return datax.Diff(aArr, bArr)
}

package stringx

import "github.com/shiningrush/goext/datax"

// HasItem test if the "item" is in "strs"
func HasItem(strs []string, item string) bool {
	for _, s := range strs {
		if s == item {
			return true
		}
	}
	return false
}

// IsSuperset test if the "src" is superset of "dest"
func IsSuperset(src, dest []string) bool {
	srcSet := datax.NewSet().AddString(src...)
	descSet := datax.NewSet().AddString(dest...)
	return srcSet.IsSupersetOf(descSet)
}

// IsSubset test if the "src" is subset of "dest"
func IsSubset(src, dest []string) bool {
	srcSet := datax.NewSet().AddString(src...)
	descSet := datax.NewSet().AddString(dest...)
	return srcSet.IsSubsetOf(descSet)
}

// IsProperSubset test if the "src" is proper subset of "dest"
func IsProperSubset(src, dest []string) bool {
	srcSet := datax.NewSet().AddString(src...)
	descSet := datax.NewSet().AddString(dest...)
	return srcSet.IsProperSubsetOf(descSet)
}

// Intersect get intersection of two strings
func Intersect(aArr, bArr []string) (intersection []string) {
	aSet := datax.NewSet().AddString(aArr...)
	bSet := datax.NewSet().AddString(bArr...)

	aSet.Intersect(bSet).Loop(func(item interface{}) (breakLoop bool) {
		intersection = append(intersection, item.(string))
		return
	})
	return
}

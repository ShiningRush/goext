package stringx

import "github.com/ShiningRush/goext/datax"

// HasItem check if the "item" is in "strs"
func HasItem(strs []string, item string) bool {
	for _, s := range strs {
		if s == item {
			return true
		}
	}
	return false
}

// IsSuperset check if the "src" is superset of "subset"
func IsSuperset(src, subset []string) bool {
	srcSet := datax.NewSet()
	for _, srcStr := range src {
		srcSet.Add(srcStr)
	}
	for _, subsetStr := range subset {
		if !srcSet.Has(subsetStr) {
			return false
		}
	}
	return true
}

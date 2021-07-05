package stringx

func HaveItem(strs []string, item string) bool {
	for _, s := range strs {
		if s == item {
			return true
		}
	}
	return false
}

package tools

func IsContainString(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func IsContainInt(slice []int, e int) int {
	for i, v := range slice {
		if e == v {
			return i
		}
	}
	return -1
}

// i is index of slice
func SliceRemove(slice []int, i int) []int {
	return append(slice[:i], slice[i+1:]...)
}

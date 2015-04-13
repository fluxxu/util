package util

func IndexOfString(slice []string, item string) int {
	for index, iter := range slice {
		if iter == item {
			return index
		}
	}
	return -1
}

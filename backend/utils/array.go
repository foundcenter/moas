package utils

func Contains(list []string, value string) bool {
	for _, s := range list {
		if s == value {
			return true
		}
	}

	return false
}

func IndexOf(list []string, value string) int {
	for i, s := range list {
		if s == value {
			return i
		}
	}

	return -1
}

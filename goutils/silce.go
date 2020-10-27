package goutils

func StringInSlice(element string, slice []string) bool {
	for _, s := range slice {
		if s == element {
			return true
		}
	}
	return false
}

func Int64InSlice(element int64, slice []int64) bool {
	for _, s := range slice {
		if s == element {
			return true
		}
	}

	return false
}

func DuplicateElementInStringSlice(slice *[]string) (string, bool) {
	for i, s := range *slice {
		for ii, v := range *slice {
			if i != ii && s == v {
				return s, true
			}
		}
	}
	return "", false
}

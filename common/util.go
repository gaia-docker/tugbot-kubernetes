package common

func SliceContains(val string, s []string) bool {
	ret := false
	for _, curr := range s {
		if curr == val {
			ret = true
			break
		}
	}

	return ret
}
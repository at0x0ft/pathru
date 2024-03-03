package cmd

func stringArrayEquals(ls1 []string, ls2 []string) bool {
	if ls1 == nil && ls2 == nil {
		return true
	} else if ls1 == nil || ls2 == nil {
		return false
	}
	if len(ls1) != len(ls2) {
		return false
	}
	for i, e1 := range ls1 {
		if e1 != ls2[i] {
			return false
		}
	}
	return true
}

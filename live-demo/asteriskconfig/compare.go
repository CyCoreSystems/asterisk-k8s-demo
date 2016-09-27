package main

func differ(X, Y []string) bool {
	if len(X) != len(Y) {
		return true
	}

	counts := make(map[string]int)
	var balance int
	for _, val := range X {
		counts[val]++
		balance++
	}
	for _, val := range Y {
		if _, ok := counts[val]; ok {
			counts[val]++
			balance--
			if counts[val] < 2 {
				return true
			}
		}
	}

	if balance != 0 {
		return true
	}

	for _, count := range counts {
		if count < 2 {
			return true
		}
	}
	return false
}

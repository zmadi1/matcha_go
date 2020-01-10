package validation

func IsDigit(str string) bool {
	itr := 0
	for itr < len(str) {
		if str[itr] < '0' || str[itr] > '9' {
			return false
		}
		itr++
	}
	return true
}

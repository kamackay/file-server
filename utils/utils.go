package utils

func Async(f func()) {
	go f()
}

func TernaryString(cond bool, positive string, negative string) string {
	if cond {
		return positive
	} else {
		return negative
	}
}

func StrArrIncludes(arr []string, val string) int {
	for x, v := range arr {
		if v == val {
			return x
		}
	}
	return -1
}

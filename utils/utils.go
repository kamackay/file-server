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

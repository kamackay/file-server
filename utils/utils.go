package utils

func Async(f func()) {
	go f()
}

package files

import "os"

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func IsDir(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func GetSize(filename string) int64 {
	if stat, err := os.Stat(filename); err != nil {
		return 0
	} else {
		return stat.Size()
	}
}

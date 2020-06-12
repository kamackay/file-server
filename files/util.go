package files

import (
	"io/ioutil"
	"os"
	"path"
	"regexp"
)

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

func GetFolderSize(folder string) int64 {
	fs, err := ioutil.ReadDir(folder)
	if err != nil {
		return 0
	}
	size := int64(0)
	for _, f := range fs {
		size += GetSize(path.Join(folder, f.Name()))
	}
	return size
}

func CountFolderItems(path string) int {
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		return 0
	}
	count := 0
	for _, f := range fs {
		if !regexp.MustCompile(".*\\.meta").MatchString(f.Name()) {
			// Not a Metadata File
			count++
		}
	}
	return count
}

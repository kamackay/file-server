package main

import "os"

func GetFile(path string) (os.FileInfo, bool, error) {
	if fi, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fi, false, err
		} else {
			return fi, true, err
		}
	} else {
		return fi, true, err
	}
}

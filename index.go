package main

import (
	"gitlab.com/kamackay/filer/server"
	"os"
)

var root = os.Getenv("ROOT_PATH")

func main() {
	if len(root) == 0 {
		root = "/files"
	}

	server.New(root).Start()
}

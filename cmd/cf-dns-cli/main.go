package main

import (
	"os"
	"path/filepath"
)

func main() {
	run(filepath.Base(os.Args[0]), os.Args[1:])
}

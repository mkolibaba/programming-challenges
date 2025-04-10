package puzzle

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const inputFolder = "input"

func GetFile(name string) *os.File {
	_, f, _, _ := runtime.Caller(0)
	result := filepath.Join(filepath.Dir(f), inputFolder, name)
	file, err := os.Open(result)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

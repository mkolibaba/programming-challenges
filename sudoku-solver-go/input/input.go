package input

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const inputFolder = "resources"

var (
	mapStringToInt = map[string]int{
		"0": 0, ".": 0, "1": 1, "2": 2, "3": 3, "4": 4,
		"5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
	}
)

type Parser interface {
	fmt.Stringer
	Parse() [][]int
}

func GetFile(name string) (*os.File, func()) {
	_, f, _, _ := runtime.Caller(0)
	result := filepath.Join(filepath.Dir(f), inputFolder, name)
	file, err := os.Open(result)
	if err != nil {
		log.Fatal(err)
	}

	closeFunc := func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}

	return file, closeFunc
}

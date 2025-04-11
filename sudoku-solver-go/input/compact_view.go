package input

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type CompactViewFile struct {
	filename string
}

func NewCompactViewFile(filename string) *CompactViewFile {
	return &CompactViewFile{filename: filename}
}

func (c CompactViewFile) Parse() [][]int {
	file, closeFunc := GetFile(c.filename)
	defer closeFunc()
	return read(file)
}

func (c CompactViewFile) String() string {
	return fmt.Sprintf("Type: %T\nFile name: %s\n", c, c.filename)
}

type CompactViewString struct {
	str string
}

func NewCompactViewString(str string) *CompactViewString {
	return &CompactViewString{str: str}
}

func (c CompactViewString) Parse() [][]int {
	return read(strings.NewReader(c.str))
}

func (c CompactViewString) String() string {
	return fmt.Sprintf("Type: %T\n", c)
}

func read(r io.Reader) [][]int {
	scanner := bufio.NewScanner(r)
	sudoku := make([][]int, 9)

	for s := 0; scanner.Scan(); s++ {
		row := make([]int, 9)
		for i, v := range scanner.Text() {
			row[i] = mapStringToInt[string(v)]
		}
		sudoku[s] = row
	}

	return sudoku
}

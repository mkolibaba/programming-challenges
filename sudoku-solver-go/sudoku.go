package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func readFromFile(path string) [][]int {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	return read(file)
}

func readFromString(s string) [][]int {
	return read(strings.NewReader(s))
}

func read(r io.Reader) (sudoku [][]int) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		var row []int
		for _, v := range strings.Split(scanner.Text(), "") {
			row = append(row, cellStringToInt(v))
		}
		sudoku = append(sudoku, row)
	}

	return sudoku
}

func solved(sudoku [][]int) bool {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if !checkCell(sudoku, i, j) {
				return false
			}
		}
	}
	return true
}

func checkCell(sudoku [][]int, row, column int) bool {
	value := sudoku[row][column]
	// validate cell
	if value < 1 || value > 9 {
		return false
	}

	// validate row
	for i := 0; i < 9; i++ {
		if sudoku[row][i] == value && i != column {
			return false
		}
	}

	// validate column
	for i := 0; i < 9; i++ {
		if sudoku[i][column] == value && i != row {
			return false
		}
	}

	// validate square
	squareRowStart := row / 3 * 3
	squareColumnStart := column / 3 * 3
	for i := 3 * squareRowStart; i < squareRowStart+3; i++ {
		for j := 3 * squareColumnStart; j < squareColumnStart+3; j++ {
			if sudoku[i][j] == value && i != row && j != column {
				return false
			}
		}
	}

	return true
}

func cellStringToInt(s string) (c int) {
	if s == "." {
		return 0
	}
	c, _ = strconv.Atoi(s)
	return
}

func cellIntToString(c int) string {
	if c == 0 {
		return "."
	}
	return strconv.Itoa(c)
}

var squaresSeparator = "-------|-------|-------"
var toAnyArray = func(a []int) (r []any) {
	for _, v := range a {
		r = append(r, cellIntToString(v))
	}
	return
}

func prettyPrint(sudoku [][]int) string {
	builder := &strings.Builder{}
	for rowIdx, row := range sudoku {
		builder.WriteString(fmt.Sprintf(" %v %v %v | %v %v %v | %v %v %v ", toAnyArray(row)...))
		builder.WriteString("\n")
		if rowIdx == 2 || rowIdx == 5 {
			builder.WriteString(squaresSeparator + "\n")
		}
	}
	return builder.String()
}

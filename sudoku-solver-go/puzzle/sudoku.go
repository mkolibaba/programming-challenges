package puzzle

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

const (
	squaresSeparator = "-------|-------|-------"
)

type Sudoku [][]int

func NewFromFile(filename string) Sudoku {
	return read(GetFile(filename))
}

func NewFromString(input string) Sudoku {
	return read(strings.NewReader(input))
}

func NewFromKaggle(puzzleNo int) Sudoku {
	file := GetFile("kaggle-sudoku.csv")
	scanner := bufio.NewScanner(file)
	for i := 0; i < puzzleNo; i++ {
		scanner.Scan()
	}
	quiz := strings.Split(scanner.Text(), ",")[0]
	if len(quiz) != 9*9 {
		log.Fatalf("kaggle puzzle No %d invalid", puzzleNo)
	}
	var sudoku = make(Sudoku, 9)
	for i := 0; i < len(sudoku); i++ {
		sudoku[i] = make([]int, 9)
		for j, v := range strings.Split(quiz[i*9:i*9+9], "") {
			sudoku[i][j] = cellStringToInt(v)
		}
	}
	return sudoku
}

func (s Sudoku) Solve() {
	i, j := getNextBlank(s)
	if i == -1 && j == -1 {
		return
	}
	solveNext(s, i, j)
}

func (s Sudoku) IsSolved() bool {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if !checkCell(s, i, j) {
				return false
			}
		}
	}
	return true
}

func (s Sudoku) String() string {
	builder := &strings.Builder{}
	for rowIdx, row := range s {
		builder.WriteString(fmt.Sprintf(" %v %v %v | %v %v %v | %v %v %v ", toAnyArray(row)...))
		builder.WriteString("\n")
		if rowIdx == 2 || rowIdx == 5 {
			builder.WriteString(squaresSeparator + "\n")
		}
	}
	return builder.String()
}

func read(r io.Reader) (sudoku Sudoku) {
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
	for i := squareRowStart; i < squareRowStart+3; i++ {
		for j := squareColumnStart; j < squareColumnStart+3; j++ {
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

func toAnyArray(a []int) (r []any) {
	for _, v := range a {
		r = append(r, cellIntToString(v))
	}
	return
}

func solveNext(sudoku [][]int, row, column int) bool {
	for possibleValue := 1; possibleValue <= 9; possibleValue++ {
		sudoku[row][column] = possibleValue
		if checkCell(sudoku, row, column) {
			i, j := getNextBlank(sudoku)
			if i == -1 && j == -1 {
				return true
			}
			if solveNext(sudoku, i, j) {
				return true
			}
		}
	}
	// ни одно из значений не подходит
	sudoku[row][column] = 0
	return false
}

func getNextBlank(sudoku [][]int) (int, int) {
	for i, row := range sudoku {
		for j, cell := range row {
			if cell == 0 {
				return i, j
			}
		}
	}
	return -1, -1
}

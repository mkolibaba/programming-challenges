package puzzle

import (
	"bufio"
	"fmt"
	"github.com/loov/hrtime"
	"io"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	squaresSeparator   = "-------|-------|-------"
	kaggleSudokusCount = 1_000_000
)

type Sudoku struct {
	Puzzle      [][]int
	TimeElapsed time.Duration
	Moves       int
	solved      bool
}

func NewFromFile(filename string) *Sudoku {
	return read(GetFile(filename))
}

func NewFromString(input string) *Sudoku {
	return read(strings.NewReader(input))
}

func NewRandomFromKaggle() *Sudoku {
	return NewFromKaggle(rand.Intn(kaggleSudokusCount))
}

func NewFromKaggle(puzzleNo int) *Sudoku {
	file := GetFile("kaggle-sudoku.csv")
	scanner := bufio.NewScanner(file)
	for i := 0; i < puzzleNo; i++ {
		scanner.Scan()
	}
	quiz := strings.Split(scanner.Text(), ",")[0]
	if len(quiz) != 9*9 {
		log.Fatalf("kaggle puzzle No %d invalid", puzzleNo)
	}
	var sudoku = make([][]int, 9)
	for i := 0; i < len(sudoku); i++ {
		sudoku[i] = make([]int, 9)
		for j, v := range strings.Split(quiz[i*9:i*9+9], "") {
			sudoku[i][j] = cellStringToInt(v)
		}
	}
	return &Sudoku{Puzzle: sudoku}
}

func (s *Sudoku) Solve() {
	start := hrtime.Now()
	// start
	i, j := getNextBlank(s.Puzzle)
	if i == -1 && j == -1 {
		return
	}
	s.solveNext(i, j)
	// end
	s.TimeElapsed = hrtime.Since(start)
}

func (s *Sudoku) IsSolved() bool {
	if s.solved {
		return true
	}

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if !checkCell(s.Puzzle, i, j) {
				return false
			}
		}
	}
	s.solved = true
	return true
}

func (s *Sudoku) String() string {
	builder := &strings.Builder{}
	for rowIdx, row := range s.Puzzle {
		builder.WriteString(fmt.Sprintf(" %v %v %v | %v %v %v | %v %v %v ", toAnyArray(row)...))
		builder.WriteString("\n")
		if rowIdx == 2 || rowIdx == 5 {
			builder.WriteString(squaresSeparator + "\n")
		}
	}
	return builder.String()
}

func read(r io.Reader) *Sudoku {
	var sudoku [][]int
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		var row []int
		for _, v := range strings.Split(scanner.Text(), "") {
			row = append(row, cellStringToInt(v))
		}
		sudoku = append(sudoku, row)
	}

	return &Sudoku{Puzzle: sudoku}
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

func (s *Sudoku) solveNext(row, column int) bool {
	for possibleValue := 1; possibleValue <= 9; possibleValue++ {
		s.Puzzle[row][column] = possibleValue
		s.Moves++
		if checkCell(s.Puzzle, row, column) {
			i, j := getNextBlank(s.Puzzle)
			if i == -1 && j == -1 {
				return true
			}
			if s.solveNext(i, j) {
				return true
			}
		}
	}
	// ни одно из значений не подходит
	s.Puzzle[row][column] = 0
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

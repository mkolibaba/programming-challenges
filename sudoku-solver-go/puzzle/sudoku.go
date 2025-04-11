package puzzle

import (
	"fmt"
	"github.com/loov/hrtime"
	"github.com/mkolibaba/programming-challenges/sudoku-solver-go/input"
	"strings"
	"time"
)

const (
	squaresSeparator = "-------|-------|-------"
)

var (
	mapIntToString = map[int]string{
		0: ".", 1: "1", 2: "2", 3: "3", 4: "4",
		5: "5", 6: "6", 7: "7", 8: "8", 9: "9",
	}
)

type Puzzle [][]int

type Sudoku struct {
	Puzzle      [][]int
	Parser      input.Parser
	TimeElapsed time.Duration
	Moves       int
	solved      bool
}

func NewFromFile(filename string) *Sudoku {
	return NewSudoku(input.NewCompactViewFile(filename))
}

func NewFromString(s string) *Sudoku {
	return NewSudoku(input.NewCompactViewString(s))
}

func NewRandomFromKaggle() *Sudoku {
	return NewSudoku(input.NewKaggle())
}

func NewFromKaggle(puzzleNo int) *Sudoku {
	return NewSudoku(input.NewKaggleN(puzzleNo))
}

func NewSudoku(parser input.Parser) *Sudoku {
	return &Sudoku{Puzzle: parser.Parse(), Parser: parser}
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
	for rowIdx, _ := range s.Puzzle {
		builder.WriteString(fmt.Sprintf(" %s %s %s | %s %s %s | %s %s %s ", s.getRowStringRepresentation(rowIdx)...))
		builder.WriteString("\n")
		if rowIdx == 2 || rowIdx == 5 {
			builder.WriteString(squaresSeparator + "\n")
		}
	}
	return builder.String()
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

func (s *Sudoku) getRowStringRepresentation(n int) []any {
	row := s.Puzzle[n]
	result := make([]any, len(row))
	for i, v := range row {
		result[i] = mapIntToString[v]
	}
	return result
}

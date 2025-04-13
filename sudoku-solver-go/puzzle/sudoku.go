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
	s.solved = s.backtrack()
	s.TimeElapsed = hrtime.Since(start)
}

func (s *Sudoku) IsSolved() bool {
	return s.solved
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

func checkCell(sudoku [][]int, value, row, column int) bool {
	// TODO: как будто можно убрать
	// validate cell
	if value < 1 || value > 9 {
		return false
	}

	// validate row & column
	for i := 0; i < 9; i++ {
		if sudoku[row][i] == value || sudoku[i][column] == value {
			return false
		}
	}

	// validate square
	squareRowStart := row / 3 * 3
	squareColumnStart := column / 3 * 3
	for i := squareRowStart; i < squareRowStart+3; i++ {
		for j := squareColumnStart; j < squareColumnStart+3; j++ {
			if sudoku[i][j] == value {
				return false
			}
		}
	}

	return true
}

func (s *Sudoku) backtrack() bool {
	for row := 0; row < 9; row++ {
		for column := 0; column < 9; column++ {
			if s.Puzzle[row][column] == 0 {
				for possibleValue := 1; possibleValue <= 9; possibleValue++ {
					if checkCell(s.Puzzle, possibleValue, row, column) {
						s.Puzzle[row][column] = possibleValue
						s.Moves++
						if s.backtrack() {
							return true
						}
					}
				}
				// ни одно из значений не подходит
				s.Puzzle[row][column] = 0
				return false
			}
		}
	}
	return true
}

func (s *Sudoku) getRowStringRepresentation(n int) []any {
	row := s.Puzzle[n]
	result := make([]any, len(row))
	for i, v := range row {
		result[i] = mapIntToString[v]
	}
	return result
}

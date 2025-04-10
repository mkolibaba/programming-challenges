package puzzle_test

import (
	"github.com/mkolibaba/programming-challenges/sudoku-solver-go/puzzle"
	"reflect"
	"testing"
)

const (
	unsolvedSudokuFilename = "test_unsolved.txt"
	solvedSudokuFilename   = "test_solved.txt"
	unsolvedSudoku         = `...68..32
..6.74...
..395....
.7....9..
4........
.957..4.8
9........
.8.4..6..
.....35..`
)

var (
	oneCellMissing    = punchCellsFromSolved(Cell{0, 0})
	threeCellsMissing = punchCellsFromSolved(Cell{0, 0}, Cell{0, 1}, Cell{0, 2})
)

type Cell struct {
	row, column int
}

func TestPrettyPrint(t *testing.T) {
	sudoku := puzzle.NewFromFile(unsolvedSudokuFilename)

	t.Logf("\n%s", sudoku)
}

func TestReadFromFile(t *testing.T) {
	sudoku := puzzle.NewFromFile(unsolvedSudokuFilename)

	assertValidSudoku(t, sudoku)
}

func TestReadFromString(t *testing.T) {
	sudoku := puzzle.NewFromString(unsolvedSudoku)

	assertValidSudoku(t, sudoku)
}

func TestSolved(t *testing.T) {
	t.Run("check solved sudoku", func(t *testing.T) {
		sudoku := puzzle.NewFromFile(solvedSudokuFilename)

		assertValidSudoku(t, sudoku)

		if !sudoku.IsSolved() {
			t.Errorf("Expected sudoku is solved, got false")
		}
	})
}

func TestSolve(t *testing.T) {
	cases := []struct {
		name   string
		sudoku *puzzle.Sudoku
	}{
		{
			"solved sudoku is solved",
			punchCellsFromSolved(),
		},
		{
			"solve sudoku when cell (0, 0) is missing",
			oneCellMissing,
		},
		{
			"solve sudoku when cells (0, 0), (0, 1), (0, 2) are missing",
			threeCellsMissing,
		},
		{
			"solve when cell (1, 0) is missing",
			punchCellsFromSolved(Cell{1, 0}),
		},
		{
			"solve brand new sudoku",
			puzzle.NewFromFile(unsolvedSudokuFilename),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.sudoku.Solve()

			if !c.sudoku.IsSolved() {
				t.Errorf("Expected sudoku is solved, got false")
			}
		})
	}
}

func TestSolvedCorrectly(t *testing.T) {
	sudoku := puzzle.NewFromFile(unsolvedSudokuFilename)
	solved := puzzle.NewFromFile(solvedSudokuFilename)

	sudoku.Solve()

	if !reflect.DeepEqual(sudoku.Puzzle, solved.Puzzle) {
		t.Errorf("invalid solution\ngot:\n%s\nwant:\n%s", sudoku, solved)
	}
}

func BenchmarkSolve(b *testing.B) {
	for b.Loop() {
		sudoku := puzzle.NewFromString(unsolvedSudoku)
		sudoku.Solve()
	}
}

func assertValidSudoku(t *testing.T, sudoku *puzzle.Sudoku) {
	t.Helper()
	if len(sudoku.Puzzle) != 9 {
		t.Errorf("Expected 9 rows, got %d", len(sudoku.Puzzle))
	}

	for i := 0; i < 9; i++ {
		if len(sudoku.Puzzle[i]) != 9 {
			t.Errorf("Expected 9 columns in row %d, got %d", i, len(sudoku.Puzzle[i]))
		}
	}

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if sudoku.Puzzle[i][j] < 0 || sudoku.Puzzle[i][j] > 9 {
				t.Errorf("Expected cell (%d, %d) is between 0 and 9, got %d", i, j, sudoku.Puzzle[i][j])
			}
		}
	}
}

func punchCellsFromSolved(cells ...Cell) *puzzle.Sudoku {
	sudoku := puzzle.NewFromFile(solvedSudokuFilename)

	for _, c := range cells {
		sudoku.Puzzle[c.row][c.column] = 0
	}
	return sudoku
}

package puzzle

import (
	"testing"
)

const (
	unsolvedSudokuPath = "../input/1.txt"
	solvedSudokuPath   = "../input/1_solved.txt"
	unsolvedSudoku     = `...68..32
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

func TestReadFromFile(t *testing.T) {
	sudoku := ReadFromFile(unsolvedSudokuPath)

	assertValidSudoku(t, sudoku)
}

func TestReadFromString(t *testing.T) {
	sudoku := readFromString(unsolvedSudoku)

	assertValidSudoku(t, sudoku)
}

func TestCheckCell(t *testing.T) {
	sudoku := ReadFromFile(solvedSudokuPath)

	assertValidSudoku(t, sudoku)

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if !checkCell(sudoku, i, j) {
				t.Errorf("Expected cell [%d][%d] is solved, got false", i, j)
			}
		}
	}
}

func TestSolved(t *testing.T) {
	t.Run("check solved sudoku", func(t *testing.T) {
		sudoku := ReadFromFile(solvedSudokuPath)

		assertValidSudoku(t, sudoku)

		if !solved(sudoku) {
			t.Errorf("Expected sudoku is solved, got false")
		}
	})
}

func TestSolve(t *testing.T) {
	cases := []struct {
		name   string
		sudoku [][]int
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
			ReadFromFile(unsolvedSudokuPath),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			solve(c.sudoku)

			if !solved(c.sudoku) {
				t.Errorf("Expected sudoku is solved, got false")
			}
		})
	}
}

func TestGetNextBlank(t *testing.T) {
	t.Run("should find next blank", func(t *testing.T) {
		blankRow := 0
		blankColumn := 2
		sudoku := punchCellsFromSolved(Cell{blankRow, blankColumn})

		i, j := getNextBlank(sudoku)

		if i != blankRow || j != blankColumn {
			t.Errorf("Expected (%d, %d) as next blank cell, got (%d, %d)", blankRow, blankColumn, i, j)
		}
	})
	t.Run("should find next blank when no blank", func(t *testing.T) {
		sudoku := punchCellsFromSolved()

		i, j := getNextBlank(sudoku)

		if i != -1 || j != -1 {
			t.Errorf("Expected (-1, -1) as next blank cell, got (%d, %d)", i, j)
		}
	})
}

func BenchmarkSolve(b *testing.B) {
	for b.Loop() {
		sudoku := readFromString(unsolvedSudoku)
		solve(sudoku)
	}
}

func assertValidSudoku(t *testing.T, sudoku [][]int) {
	t.Helper()
	if len(sudoku) != 9 {
		t.Errorf("Expected 9 rows, got %d", len(sudoku))
	}

	for i := 0; i < 9; i++ {
		if len(sudoku[i]) != 9 {
			t.Errorf("Expected 9 columns in row %d, got %d", i, len(sudoku[i]))
		}
	}

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if sudoku[i][j] < 0 || sudoku[i][j] > 9 {
				t.Errorf("Expected cell (%d, %d) is between 0 and 9, got %d", i, j, sudoku[i][j])
			}
		}
	}
}

func punchCellsFromSolved(cells ...Cell) [][]int {
	sudoku := ReadFromFile(solvedSudokuPath)

	for _, c := range cells {
		sudoku[c.row][c.column] = 0
	}
	return sudoku
}

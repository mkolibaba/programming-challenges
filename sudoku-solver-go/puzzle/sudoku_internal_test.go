package puzzle

import "testing"

const invalidSudoku = `149685732
856174293
723951846
271346985
432598167
395762418
968217354
587439621
614823579`

func TestCheckCell(t *testing.T) {
	sudoku := NewFromString(invalidSudoku)

	if checkCell(sudoku, 1, 8) {
		t.Errorf("want check cell is false, got true")
	}
}

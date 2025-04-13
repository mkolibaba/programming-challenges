package input

import (
	"fmt"
	"io"
	"log"
	"math/rand"
)

const kaggleSudokusCount = 1_000_000
const csvRowLength = 81 + 1 + 81 + 1
const csvHeaderLength = 18 // "quizzes,solutions" + 1

type Kaggle struct {
	puzzleNo int
}

// https://www.kaggle.com/datasets/bryanpark/sudoku?resource=download
func NewKaggle() *Kaggle {
	return &Kaggle{puzzleNo: rand.Intn(kaggleSudokusCount)}
}

func NewKaggleN(puzzleNo int) *Kaggle {
	return &Kaggle{puzzleNo: puzzleNo}
}

func (k *Kaggle) Parse() [][]int {
	file, closeFunc := GetFile("kaggle-sudoku.csv")
	defer closeFunc()

	skip := int64(csvHeaderLength + k.puzzleNo*csvRowLength)
	bytes := make([]byte, 81)
	file.Seek(skip, io.SeekStart) // TODO: errors
	file.Read(bytes)              // TODO: errors

	quiz := string(bytes)
	if len(quiz) != 9*9 {
		log.Fatalf("kaggle puzzle No %d invalid", k.puzzleNo)
	}
	var sudoku = make([][]int, 9)
	for i := 0; i < len(sudoku); i++ {
		sudoku[i] = make([]int, 9)
		for j, v := range quiz[i*9 : i*9+9] {
			sudoku[i][j] = mapStringToInt[string(v)]
		}
	}

	return sudoku
}

func (k *Kaggle) String() string {
	return fmt.Sprintf("Type: %T\nPuzzle No: %d\n", k, k.puzzleNo)
}

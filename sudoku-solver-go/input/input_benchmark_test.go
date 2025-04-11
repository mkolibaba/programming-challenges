package input

import (
	"testing"
)

const compactViewStringSudoku = `...68..32
..6.74...
..395....
.7....9..
4........
.957..4.8
9........
.8.4..6..
.....35..`

func BenchmarkCompactViewString(b *testing.B) {
	for b.Loop() {
		parser := NewCompactViewString(compactViewStringSudoku)
		parser.Parse()
	}
}

func BenchmarkKaggle(b *testing.B) {
	for b.Loop() {
		parser := NewKaggle()
		parser.Parse()
	}
}

func BenchmarkKaggleN(b *testing.B) {
	for b.Loop() {
		parser := NewKaggleN(550_000)
		parser.Parse()
	}
}

func BenchmarkCompactViewFile(b *testing.B) {
	for b.Loop() {
		parser := NewCompactViewFile("test_unsolved.txt")
		parser.Parse()
	}
}

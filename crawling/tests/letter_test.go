package tests

import (
	"testing"
)

const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func LetterToNum(letter string) int {
	for i, a := range letters {
		if letter == string(a) {
			return i + 1
		}
	}
	return 0
}

func NumToLetter(i int) string {
	return letters[i-1 : i]
}

var testData = []struct {
	letter string
	num    int
}{
	{"A", 1},
	{"B", 2},
	{"C", 3},
	{"D", 4},
	{"E", 5},
	{"F", 6},
	{"G", 7},
	{"H", 8},
	{"I", 9},
	{"J", 10},
	{"K", 11},
	{"L", 12},
	{"M", 13},
	{"N", 14},
	{"O", 15},
	{"P", 16},
	{"Q", 17},
	{"R", 18},
	{"S", 19},
	{"T", 20},
	{"U", 21},
	{"V", 22},
	{"W", 23},
	{"X", 24},
	{"Y", 25},
	{"Z", 26},
}

func TestLetterToNum(t *testing.T) {
	for idx, data := range testData {
		if got, wanted := LetterToNum(data.letter), data.num; got != wanted {
			t.Errorf("%s wanted %d, got %d instead at idx %d",
				data.letter,
				data.num,
				got,
				idx,
			)
		}
	}
}

func TestNumToLetter(t *testing.T) {
	for idx, data := range testData {
		if got, wanted := NumToLetter(data.num), data.letter; got != wanted {
			t.Errorf("%d wanted %s, got %s instead at idx %d",
				data.num,
				data.letter,
				got,
				idx,
			)
		}
	}
}

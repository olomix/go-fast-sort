package go_fast_sort

import (
	"math/rand"
	"testing"
)

const bufSize = 104857596
const itemSize = 12
const stableSeed = 50

func bufInit(t *testing.T) []byte {
	t.Helper()
	buf := make([]byte, bufSize)
	rnd := rand.New(rand.NewSource(stableSeed))
	n, err := rnd.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if n != bufSize {
		t.Fatal(n)
	}
	return buf
}

func TestSortStd(t *testing.T) {
	buf := bufInit(t)
	sortStd(buf, itemSize)
}

func TestSortStdStable(t *testing.T) {
	buf := bufInit(t)
	sortStdStable(buf, itemSize)
}
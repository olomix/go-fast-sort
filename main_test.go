package go_fast_sort

import (
	"bytes"
	"math"
	"math/rand"
	"testing"
)

const bufSize = 104857596

//const bufSize = 28814304

//const bufSize = 12288
const itemSize = 12
const stableSeed = 50

func bufInit(t testing.TB) []byte {
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

func BenchmarkSortStd(b *testing.B) {
	b.ReportAllocs()

	buf := bufInit(b)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sortStd(buf, itemSize)
	}
}

func TestSortStdStable(t *testing.T) {
	buf := bufInit(t)
	sortStdStable(buf, itemSize)
}

func BenchmarkSortStdStable(b *testing.B) {
	b.ReportAllocs()

	buf := bufInit(b)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sortStdStable(buf, itemSize)
	}
}

func BenchmarkSortTimSort(b *testing.B) {
	b.ReportAllocs()

	buf := bufInit(b)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		timSort(buf, itemSize)
	}
}

func TestTimSort(t *testing.T) {
	bufTimSorted := bufInit(t)
	bufStdStableSorted := make([]byte, len(bufTimSorted))
	n := copy(bufStdStableSorted, bufTimSorted)
	if n != bufSize {
		t.Fatal(n)
	}

	//fmtArray(bufTimSorted)
	timSort(bufTimSorted, itemSize)
	//fmtArray(bufTimSorted)
	//fmtArray(bufStdStableSorted)
	sortStdStable(bufStdStableSorted, itemSize)
	//fmtArray(bufStdStableSorted)
	if !bytes.Equal(bufTimSorted, bufStdStableSorted) {
		t.Fatal("not equal")
	}
}

func fmtArray(tb testing.TB, in []byte, start int) {
	tb.Logf("[ start ]")
	for i := 0; i < len(in); i += itemSize {
		id := ui64(in, i)
		f := math.Float32frombits(ui32(in, i+8))
		tb.Logf("%v %v %v", i/itemSize+start, id, f)
	}
	tb.Logf("[ end ]")
}

package go_fast_sort

import (
	"bytes"
	"testing"
)

func TestTimSort2_1(t *testing.T) {
	buf := bufInit(t)

	bufStdStableSorted := make([]byte, len(buf))
	n := copy(bufStdStableSorted, buf)
	if n != bufSize {
		t.Fatal(n)
	}
	sortStdStable(bufStdStableSorted, itemSize)

	TimSort2(NewBytesTimSorter(buf, itemSize, btsLt, nil))

	if !bytes.Equal(buf, bufStdStableSorted) {
		t.Fatal("not equal")
	}
}

func btsLt(a, b []byte) bool {
	return ui64(a, 0) < ui64(b, 0)
}

func BenchmarkSortTimSort2(b *testing.B) {
	b.ReportAllocs()

	buf := bufInit(b)
	ts2 := NewBytesTimSorter(buf, itemSize, btsLt, make([]byte, bufSize/2))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		TimSort2(ts2)
	}
}

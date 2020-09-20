package go_fast_sort

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"math/rand"
	"os"
	"testing"
)

type TS2 struct {
	itemSize int
	arr      []byte
	itemBuf  []byte // buf for one item element
	tmpBuf   []byte
}

func newTS2(buf []byte, itemSize int) *TS2 {
	return &TS2{
		itemSize: itemSize,
		arr:      buf,
		itemBuf:  make([]byte, itemSize),
	}

}

func (T *TS2) Len() int {
	return len(T.arr) / itemSize
}

func (T *TS2) Less(i, j int, tmpI, tmpJ bool) bool {
	var vi, vj uint64
	if tmpI {
		vi = ui64(T.tmpBuf, i*T.itemSize)
	} else {
		vi = ui64(T.arr, i*T.itemSize)
	}
	if tmpJ {
		vj = ui64(T.tmpBuf, j*T.itemSize)
	} else {
		vj = ui64(T.arr, j*T.itemSize)
	}
	return vi < vj
}

func (T *TS2) Move(dst int, src int) {
	if dst == src {
		return
	}
	if src > dst {
		copy(T.itemBuf, T.arr[src*T.itemSize:])
		copy(
			T.arr[(dst+1)*T.itemSize:],
			T.arr[dst*T.itemSize:src*T.itemSize],
		)
		copy(T.arr[dst*T.itemSize:], T.itemBuf)
		return
	}

	panic("implement me")

}

func (T *TS2) Copy(dst, src, ln int, tmpDst, tmpSrc bool) {
	var srcBuf []byte
	if tmpSrc {
		srcBuf = T.tmpBuf
	} else {
		srcBuf = T.arr
	}
	var dstBuf []byte
	if tmpDst {
		dstBuf = T.tmpBuf
	} else {
		dstBuf = T.arr
	}
	copy(
		dstBuf[dst*T.itemSize:(dst+ln)*T.itemSize],
		srcBuf[src*T.itemSize:(src+ln)*T.itemSize],
	)
}

func (T *TS2) EnsureTempCapacity(ln int) {
	ln *= T.itemSize

	if len(T.tmpBuf) >= ln {
		return
	}

	if cap(T.tmpBuf) >= ln {
		T.tmpBuf = T.tmpBuf[:ln]
		return
	}

	T.tmpBuf = make([]byte, ln)
}

func (T *TS2) Swap(i, j int) {
	copy(T.itemBuf, T.arr[j*T.itemSize:])
	copy(
		T.arr[j*T.itemSize:(j+1)*T.itemSize],
		T.arr[i*T.itemSize:],
	)
	copy(T.arr[i*T.itemSize:], T.itemBuf)
}

func bufInit2(t testing.TB, bufSize int) []byte {
	t.Helper()
	buf := make([]byte, bufSize*itemSize)
	rnd := rand.New(rand.NewSource(stableSeed))
	for i := 0; i < bufSize; i++ {
		j := rnd.Intn(bufSize * 10)
		binary.LittleEndian.PutUint64(buf[i*itemSize:], uint64(j))
		binary.LittleEndian.PutUint32(buf[i*itemSize+8:], uint32(i))
	}
	return buf
}

func TestTimSort2_1(t *testing.T) {
	//buf := bufInit2(t, 1024)
	buf := bufInit(t)
	//writeBuf(t, buf, "orig.txt")
	//fmtArray(buf)

	//bufStdStableSorted := make([]byte, len(buf))
	//n := copy(bufStdStableSorted, buf)
	//if n != bufSize {
	//	t.Fatal(n)
	//}
	//sortStdStable(bufStdStableSorted, itemSize)
	//writeBuf(t, bufStdStableSorted, "std.txt")

	TimSort2(newTS2(buf, itemSize))

	for i := 1; i < len(buf)/itemSize; i++ {
		vi := ui64(buf, (i-1)*itemSize)
		vj := ui64(buf, i*itemSize)
		if vi > vj {
			t.Fatalf("at %v %v > %v", i, vi, vj)
		}
	}

	//writeBuf(t, buf, "tim.txt")

	//for i := 0; i < bufSize/itemSize; i++ {
	//	if !bytes.Equal(buf[i*itemSize:(i+1)*itemSize], bufStdStableSorted[i*itemSize:(i+1)*itemSize]) {
	//		j := i - 10
	//		if j < 0 {
	//			j = 0
	//		}
	//		t.Log("buf")
	//		fmtArray(t, buf[j*itemSize:(j+20)*itemSize], j)
	//		t.Log("std sorted")
	//		fmtArray(t, bufStdStableSorted[j*itemSize:(j+20)*itemSize], j)
	//		t.Fatalf("not equal at %v", i)
	//	}
	//}

	//if !bytes.Equal(buf, bufStdStableSorted) {
	//	t.Fatal("not equal")
	//}
}

func BenchmarkSortTimSort2(b *testing.B) {
	b.ReportAllocs()

	buf := bufInit(b)
	ts2 := newTS2(buf, itemSize)
	ts2.tmpBuf = make([]byte, bufSize/2)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		TimSort2(ts2)
	}
}

func writeBuf(t testing.TB, buf []byte, fname string) {
	f, err := os.Create(fname)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()
	for i := 0; i < len(buf)/itemSize; i++ {
		fmt.Fprintf(w, "%v: %v\n", i, ui64(buf, i*itemSize))
	}
}

package go_fast_sort

import (
	"encoding/binary"
	"sort"
)

func ui64(buf []byte, idx int) uint64 {
	return binary.LittleEndian.Uint64(buf[idx : idx+8])
}

func ui32(buf []byte, idx int) uint32 {
	return binary.LittleEndian.Uint32(buf[idx : idx+4])
}

func sortStd(buf []byte, itemSize int) {
	sort.Sort(&sorter{
		itemSize: itemSize,
		buf:      buf,
		tmp:      make([]byte, itemSize),
	})
}

func sortStdStable(buf []byte, itemSize int) {
	sort.Stable(&sorter{
		itemSize: itemSize,
		buf:      buf,
		tmp:      make([]byte, itemSize),
	})
}

func timSort(buf []byte, itemSize int) {
	TimSort(newTimSorter(buf, itemSize))
}

type sorter struct {
	itemSize int
	buf      []byte
	tmp      []byte
}

func (s *sorter) Len() int {
	return len(s.buf) / s.itemSize
}

func (s *sorter) Less(i, j int) bool {
	return ui64(s.buf, i*s.itemSize) < ui64(s.buf, j*s.itemSize)
}

func (s *sorter) Swap(i, j int) {
	copy(s.tmp, s.buf[j*s.itemSize:(j+1)*s.itemSize])
	copy(
		s.buf[j*s.itemSize:(j+1)*s.itemSize],
		s.buf[i*s.itemSize:(i+1)*s.itemSize],
	)
	copy(s.buf[i*s.itemSize:(i+1)*s.itemSize], s.tmp)
}

type timSorter struct {
	itemSize int
	buf      []byte
	tmpElm   []byte // buffer for one element
	tmpBuf   []byte
}

func newTimSorter(buf []byte, itemSize int) TimSorter {
	return &timSorter{
		buf:      buf,
		itemSize: itemSize,
		tmpElm:   make([]byte, itemSize),
		tmpBuf:   make([]byte, len(buf)),
	}
}

func (t *timSorter) LtEq(i, j int) bool {
	return ui64(t.buf, i*t.itemSize) <= ui64(t.buf, j*t.itemSize)
}

func (t *timSorter) Gt(i, j int) bool {
	return ui64(t.buf, i*t.itemSize) > ui64(t.buf, j*t.itemSize)
}

func (t *timSorter) Len() int {
	return len(t.buf) / t.itemSize
}

func (t *timSorter) Move(dst int, src int) {
	if dst == src {
		return
	}
	if src > dst {
		copy(t.tmpElm, t.buf[src*t.itemSize:])
		copy(
			t.buf[(dst+1)*t.itemSize:],
			t.buf[dst*t.itemSize:src*t.itemSize],
		)
		copy(t.buf[dst*t.itemSize:], t.tmpElm)
		return
	}

	panic("implement me")
}

func (t *timSorter) CopyTemp(start int, end int) {
	if end <= start {
		return
	}
	copy(t.tmpBuf, t.buf[start*t.itemSize:end*t.itemSize])
}

func (t *timSorter) TmpLtEq(tmpIdx int, idx int) bool {
	return ui64(t.tmpBuf, tmpIdx*t.itemSize) <= ui64(t.buf, idx*t.itemSize)
}

func (t *timSorter) CopyItem(dst, src int) {
	copy(
		t.buf[dst*t.itemSize:(dst+1)*t.itemSize],
		t.buf[src*t.itemSize:(src+1)*t.itemSize],
	)
}

func (t *timSorter) CopyItemsFromTemp(dst, src, ln int) {
	copy(
		t.buf[dst*t.itemSize:(dst+ln)*t.itemSize],
		t.tmpBuf[src*t.itemSize:(src+ln)*t.itemSize],
	)
}

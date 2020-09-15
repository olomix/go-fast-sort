package go_fast_sort

import (
	"encoding/binary"
	"sort"
)

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

type sorter struct {
	itemSize int
	buf []byte
	tmp []byte
}

func (s *sorter) Len() int {
	return len(s.buf) / s.itemSize
}

func (s *sorter) Less(i, j int) bool {
	return binary.LittleEndian.Uint32(s.buf[i*s.itemSize:i*s.itemSize+8]) <
		binary.LittleEndian.Uint32(s.buf[j*s.itemSize:j*s.itemSize+8])
}

func (s *sorter) Swap(i, j int) {
	copy(s.tmp, s.buf[j:j+s.itemSize])
	copy(s.buf[j:j+s.itemSize], s.buf[i:i+s.itemSize])
	copy(s.buf[i:i+s.itemSize], s.tmp)
}


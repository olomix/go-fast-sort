package go_fast_sort

import (
	"bytes"
	"testing"
)

func TestBytesTimSorter_Move(t *testing.T) {
	testCases := []struct {
		title    string
		in       []byte
		exp      []byte
		itemSize int
		src, dst int
	}{
		{
			title:    "destination less then source",
			itemSize: 2,
			src:      4,
			dst:      1,

			in:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			exp: []byte{1, 2, 9, 10, 3, 4, 5, 6, 7, 8, 11, 12, 13, 14, 15, 16},
		},
		{
			title:    "destination greater then source",
			itemSize: 2,
			src:      1,
			dst:      4,

			in:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			exp: []byte{1, 2, 5, 6, 7, 8, 9, 10, 3, 4, 11, 12, 13, 14, 15, 16},
		},
		{
			title:    "destination equal to source",
			itemSize: 2,
			src:      4,
			dst:      4,

			in:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			exp: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.title, func(t *testing.T) {
			bts := &BytesTimSorter{
				arr:      tc.in,
				itemSize: tc.itemSize,
				itemBuf:  make([]byte, tc.itemSize),
			}
			bts.Move(tc.dst, tc.src)
			if !bytes.Equal(tc.exp, bts.arr) {
				t.Fatal(bts.arr)
			}
		})
	}
}

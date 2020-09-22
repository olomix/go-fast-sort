package go_fast_sort

type BytesTimSorter struct {
	itemSize int
	arr      []byte
	itemBuf  []byte // buf for one item element
	tmpBuf   []byte
	lt       func(a, b []byte) bool
}

// Implementation of TimSorter2 interface for bytes array.
// To sort bytes array we should know length of one bytes item, and how
// to compare two bytes items.
// If tmpBuf is provided and it's length at least half of the length of
// buffer to sort, then use provided tmpBuf, else allocate new buf of required
// length
func NewBytesTimSorter(
	buf []byte, itemSize int, lt func(a, b []byte) bool,
	tmpBuf []byte,
) *BytesTimSorter {
	tmpLn := len(buf) / itemSize / 2 * itemSize
	if len(tmpBuf) < tmpLn {
		if cap(tmpBuf) < tmpLn {
			tmpBuf = make([]byte, tmpLn)
		} else {
			tmpBuf = tmpBuf[:tmpLn]
		}
	}

	return &BytesTimSorter{
		itemSize: itemSize,
		arr:      buf,
		itemBuf:  make([]byte, itemSize),
		tmpBuf:   tmpBuf,
		lt:       lt,
	}
}

func (bts *BytesTimSorter) Len() int {
	return len(bts.arr) / bts.itemSize
}

func (bts *BytesTimSorter) Less(i, j int, tmpI, tmpJ bool) bool {
	var a, b []byte
	if tmpI {
		a = bts.tmpBuf[i*bts.itemSize : (i+1)*bts.itemSize]
	} else {
		a = bts.arr[i*bts.itemSize : (i+1)*bts.itemSize]
	}
	if tmpJ {
		b = bts.tmpBuf[j*bts.itemSize : (j+1)*bts.itemSize]
	} else {
		b = bts.arr[j*bts.itemSize : (j+1)*bts.itemSize]
	}
	return bts.lt(a, b)
}

func (bts *BytesTimSorter) Move(dst int, src int) {
	if dst == src {
		return
	}

	copy(bts.itemBuf, bts.arr[src*bts.itemSize:])
	if src > dst {
		copy(
			bts.arr[(dst+1)*bts.itemSize:],
			bts.arr[dst*bts.itemSize:src*bts.itemSize],
		)
		copy(bts.arr[dst*bts.itemSize:], bts.itemBuf)
	} else {
		copy(
			bts.arr[src*bts.itemSize:],
			bts.arr[(src+1)*bts.itemSize:(dst+1)*bts.itemSize],
		)
	}
	copy(bts.arr[dst*bts.itemSize:], bts.itemBuf)
}

func (bts *BytesTimSorter) Copy(dst, src, ln int, tmpDst, tmpSrc bool) {
	var srcBuf []byte
	if tmpSrc {
		srcBuf = bts.tmpBuf
	} else {
		srcBuf = bts.arr
	}
	var dstBuf []byte
	if tmpDst {
		dstBuf = bts.tmpBuf
	} else {
		dstBuf = bts.arr
	}
	copy(
		dstBuf[dst*bts.itemSize:(dst+ln)*bts.itemSize],
		srcBuf[src*bts.itemSize:(src+ln)*bts.itemSize],
	)
}

func (bts *BytesTimSorter) EnsureTempCapacity(ln int) {
	ln *= bts.itemSize

	if len(bts.tmpBuf) >= ln {
		return
	}

	if cap(bts.tmpBuf) >= ln {
		bts.tmpBuf = bts.tmpBuf[:ln]
		return
	}

	bts.tmpBuf = make([]byte, ln)
}

func (bts *BytesTimSorter) Swap(i, j int) {
	copy(bts.itemBuf, bts.arr[j*bts.itemSize:])
	copy(
		bts.arr[j*bts.itemSize:(j+1)*bts.itemSize],
		bts.arr[i*bts.itemSize:],
	)
	copy(bts.arr[i*bts.itemSize:], bts.itemBuf)
}

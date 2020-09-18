package go_fast_sort

type TimSorter interface {
	LtEq(i, j int) bool
	Gt(i, j int) bool
	Len() int
	// Move element from position src to position dst shifting all elements
	// in between.
	Move(dst int, src int)
	CopyTemp(start int, end int)
	// index tmpIdx point to temp storage, indxe idx to regular storage
	TmpLtEq(tmpIdx int, idx int) bool
	CopyItem(dst, src int)
	CopyItemsFromTemp(dst, src, ln int)
}

type run struct {
	ptr  int
	size int
}

func TimSort(s TimSorter) {
	if s.Len() <= 1 {
		// already sorted
		return
	}

	stack := make([]run, 0, 64)
	minRunSize := 64
	currentRun := run{
		ptr:  0,
		size: 1,
	}
	ln := s.Len()
	idx := 1
	for {
		if idx == ln {
			stack = append(stack, currentRun)
			break
		}

		if s.LtEq(currentRun.ptr+currentRun.size-1, idx) {
			currentRun.size++
			idx++
			continue
		}

		if currentRun.size < minRunSize {
			firstUnsorted := currentRun.size
			currentRun.size = minRunSize
			if currentRun.ptr+currentRun.size > ln {
				currentRun.size = ln - currentRun.ptr
			}
			if firstUnsorted < currentRun.size {
				sortRun(s, currentRun, firstUnsorted)
				idx += currentRun.size - firstUnsorted
			}
		}

		stack = append(stack, currentRun)
		stack = normalizeStack(s, stack)

		if idx == ln {
			break
		}

		currentRun.ptr = idx
		currentRun.size = 1
		idx++
	}

	mergeStack(s, stack)
}

func mergeStack(s TimSorter, stack []run) {
	if len(stack) <= 1 {
		return
	}
	for len(stack) > 1 {
		// look for two consecutive runs with minimum summary length
		idx := 0
		sz := stack[idx].size + stack[idx+1].size
		for i := 1; i < len(stack)-1; i++ {
			newSz := stack[i].size + stack[i+1].size
			if newSz < sz {
				idx = i
				sz = newSz
			}
		}
		mergeRuns(s, stack[idx], stack[idx+1])
		stack[idx].size += stack[idx+1].size
		stack = append(stack[:idx+1], stack[idx+2:]...)
	}
}

func normalizeStack(s TimSorter, stack []run) []run {
LOOP:
	for {
		if len(stack) <= 1 {
			break
		}
		if len(stack) == 2 {
			if stack[1].size > stack[0].size {
				mergeRuns(s, stack[0], stack[1])
				stack[0].size += stack[1].size
				stack = stack[:1]
			}
			break
		}
		for i := len(stack) - 3; i >= 0; i-- {
			if stack[i].size > stack[i+1].size+stack[i+2].size &&
				stack[i+1].size > stack[i+2].size {
				continue
			}
			if stack[i].size < stack[i+2].size {
				mergeRuns(s, stack[i], stack[i+1])
				stack[i].size = stack[i].size + stack[i+1].size
				stack = append(stack[:i+1], stack[i+2:]...)
			} else {
				mergeRuns(s, stack[i+1], stack[i+2])
				stack[i+1].size += stack[i+2].size
				stack = append(stack[:i+2], stack[i+3:]...)
			}
			continue LOOP
		}
		break
	}
	return stack
}

func mergeRuns(s TimSorter, l, r run) {
	if l.size < 1 {
		panic("left run has zero size")
	}
	if r.size < 1 {
		panic("right run has zero size")
	}

	var elm1 int // the pointer to first element in left run that is greater
	// then first element in right run
	for elm1 = l.ptr; elm1 < l.ptr+l.size; elm1++ {
		if s.LtEq(elm1, r.ptr) {
			continue
		}
		break
	}
	if elm1 == l.ptr+l.size {
		// two runs already merged
		return
	}

	s.CopyTemp(elm1, l.ptr+l.size)

	lastElm := findLastElm(s, l, r)

	cursorL := 0
	lnL := l.ptr + l.size - elm1
	cursorR := r.ptr
	lnR := lastElm - r.ptr
	dst := elm1

	for lnL > 0 && lnR > 0 {
		if s.TmpLtEq(cursorL, cursorR) {
			s.CopyItemsFromTemp(dst, cursorL, 1)
			cursorL++
			lnL--
		} else {
			s.CopyItem(dst, cursorR)
			cursorR++
			lnR--
		}
		dst++
	}

	if lnL > 0 {
		s.CopyItemsFromTemp(dst, cursorL, lnL)
	}

	for lnR > 0 {
		panic("right run should be exhausted")
	}
}

// find last element in right run that is greader then last element is left run
// example:
// left: [1,2,3,4]
// right: [3,4,5,6]
// return index of value 5 from right run
func findLastElm(s TimSorter, l, r run) int {
	var lastElm int
	for lastElm = r.ptr + r.size; lastElm > r.ptr; lastElm-- {
		if s.LtEq(l.ptr+l.size-1, lastElm-1) {
			continue
		}
		break
	}
	return lastElm
}

// sort current run where index of first unsorted element is unsortedIdx
func sortRun(s TimSorter, currentRun run, unsortedIdx int) {
	// asserted that unsortedIdx >= 1
	for unsortedIdx < currentRun.size {
		i := unsortedIdx - 1
		for {
			if i < 0 {
				break
			}
			if s.Gt(currentRun.ptr+i, currentRun.ptr+unsortedIdx) {
				i--
			} else {
				break
			}
		}
		i++
		s.Move(currentRun.ptr+i, currentRun.ptr+unsortedIdx)

		unsortedIdx++
	}
}

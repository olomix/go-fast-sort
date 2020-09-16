package go_fast_sort

type TimSorter interface {
	LtEq(i, j int) bool
	Gt(i, j int) bool
	Len() int
	// Move element from position src to position dst shifting all elements
	// in between.
	Move(dst int, src int)
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
		normalizeStack(s, stack)

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
	panic("not implemented")
}

func normalizeStack(s TimSorter, stack []run) {
	panic("not implemented")
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

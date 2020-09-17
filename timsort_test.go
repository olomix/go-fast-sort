package go_fast_sort

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

type intSorter struct {
	arr []int
	tmp []int
}

func (is *intSorter) CopyItem(dst, src int) {
	is.arr[dst] = is.arr[src]
}

func (is *intSorter) CopyItemsFromTemp(dst, src, ln int) {
	if ln == 1 {
		is.arr[dst] = is.tmp[src]
	} else {
		copy(is.arr[dst:dst+ln], is.tmp[src:src+ln])
	}
}

func (is *intSorter) TmpLtEq(tmpIdx int, idx int) bool {
	return is.tmp[tmpIdx] <= is.arr[idx]
}

func (is *intSorter) CopyTemp(start int, end int) {
	if end <= start {
		return
	}
	if cap(is.tmp) < end-start {
		is.tmp = make([]int, end-start)
	} else {
		is.tmp = is.tmp[0 : end-start]
	}
	copy(is.tmp, is.arr[start:end])
}

func (is *intSorter) LtEq(i, j int) bool {
	return is.arr[i] <= is.arr[j]
}

func (is *intSorter) Gt(i, j int) bool {
	return is.arr[i] > is.arr[j]
}

func (is *intSorter) Len() int {
	return len(is.arr)
}

func (is *intSorter) Move(dst int, src int) {
	if dst == src {
		return
	}
	if src > dst {
		tmp := is.arr[src]
		copy(is.arr[dst+1:], is.arr[dst:src])
		is.arr[dst] = tmp
	}
}

type testS1 struct {
	key   int
	value int
}

type structSorter struct {
	arr []testS1
	tmp []testS1
}

func (ss *structSorter) CopyItem(dst, src int) {
	ss.arr[dst] = ss.arr[src]
}

func (ss *structSorter) CopyItemsFromTemp(dst, src, ln int) {
	if ln == 1 {
		ss.arr[dst] = ss.tmp[src]
	} else {
		copy(ss.arr[dst:dst+ln], ss.tmp[src:src+ln])
	}
}

func (ss *structSorter) TmpLtEq(tmpIdx int, idx int) bool {
	return ss.tmp[tmpIdx].key <= ss.arr[idx].key
}

func (ss *structSorter) CopyTemp(start int, end int) {
	if end <= start {
		return
	}
	if cap(ss.tmp) < end-start {
		ss.tmp = make([]testS1, end-start)
	} else {
		ss.tmp = ss.tmp[0 : end-start]
	}
	copy(ss.tmp, ss.arr[start:end])
}

func (ss *structSorter) Less(i, j int) bool {
	return ss.arr[i].key < ss.arr[j].key
}

func (ss *structSorter) Swap(i, j int) {
	ss.arr[i], ss.arr[j] = ss.arr[j], ss.arr[i]
}

func (ss *structSorter) LtEq(i, j int) bool {
	return ss.arr[i].key <= ss.arr[j].key
}

func (ss *structSorter) Gt(i, j int) bool {
	return ss.arr[i].key > ss.arr[j].key
}

func (ss *structSorter) Len() int {
	return len(ss.arr)
}

func (ss *structSorter) Move(dst int, src int) {
	if dst == src {
		return
	}
	if src > dst {
		tmp := ss.arr[src]
		copy(ss.arr[dst+1:], ss.arr[dst:src])
		ss.arr[dst] = tmp
	}
}

func TestSortRun(t *testing.T) {
	is := intSorter{arr: []int{10, 1, 2, 3, 4, 7, 5, 6, 6}}
	sortRun(&is, run{ptr: 3, size: 5}, 2)
	expected := []int{10, 1, 2, 3, 4, 5, 6, 7, 6}
	if !reflect.DeepEqual(is.arr, expected) {
		t.Fatal(is.arr)
	}
}

func TestSortRun2(t *testing.T) {
	const arrSz = 200
	const startRun = 50
	const runSz = 100
	is := &intSorter{arr: make([]int, arrSz)}
	rnd := rand.New(rand.NewSource(50))
	for i := range is.arr {
		is.arr[i] = rnd.Int()
	}

	firstUnsorted := 0
	prev := is.arr[startRun]
	for i := startRun + 1; i < startRun+runSz; i++ {
		if prev <= is.arr[i] {
			prev = is.arr[i]
			continue
		}
		firstUnsorted = i - startRun
		break
	}

	sortRun(is, run{ptr: startRun, size: runSz}, firstUnsorted)

	prev = is.arr[startRun]
	for i := startRun + 1; i < startRun+runSz; i++ {
		if prev <= is.arr[i] {
			prev = is.arr[i]
			continue
		}
		t.Fatalf("at %v prev %v, current %v", i, prev, is.arr[i])
	}

	sorted := sort.SliceIsSorted(
		is.arr[startRun:startRun+runSz],
		func(i, j int) bool {
			return is.arr[startRun+i] < is.arr[startRun+j]
		},
	)
	if !sorted {
		t.Fatal("array is not sorted")
	}
}

// test if sortRun is stable
func TestSortRunIsStable(t *testing.T) {
	const arrSz = 200
	const startRun = 50
	const runSz = 100
	is := &structSorter{arr: make([]testS1, arrSz)}
	rnd := rand.New(rand.NewSource(50))
	for i := range is.arr {
		is.arr[i].key = rnd.Int()
		is.arr[i].value = rnd.Int()
	}

	// duplicate is
	is2 := &structSorter{arr: make([]testS1, arrSz)}
	for i := range is.arr {
		is2.arr[i] = is.arr[i]
	}

	firstUnsorted := 0
	prev := is.arr[startRun].key
	for i := startRun + 1; i < startRun+runSz; i++ {
		if prev <= is.arr[i].key {
			prev = is.arr[i].key
			continue
		}
		firstUnsorted = i - startRun
		break
	}

	// sort original arr
	sortRun(is, run{ptr: startRun, size: runSz}, firstUnsorted)

	// sort golden arr by std method
	is3 := &structSorter{arr: is2.arr[startRun : startRun+runSz]}
	sort.Stable(is3)

	for i := range is.arr {
		if is.arr[i] != is2.arr[i] {
			t.Fatalf(
				"at %v elemnts differ %v != %v",
				i, is.arr[i], is2.arr[i],
			)
		}
	}
}

func Test1(t *testing.T) {
	var i int
	for i = 3; i < 5; i++ {

	}
	t.Log(i)
}

func TestFindLastElm(t *testing.T) {
	testCases := []struct {
		title    string
		arr      []int
		expected int
		l, r     run
	}{
		{
			title:    "first",
			arr:      []int{1, 2, 3, 4, 3, 4, 5, 6},
			expected: 5,
			l:        run{ptr: 0, size: 4},
			r:        run{ptr: 4, size: 4},
		},
		{
			title:    "no such element",
			arr:      []int{11, 12, 13, 14, 3, 4, 5, 6},
			expected: 8,
			l:        run{ptr: 0, size: 4},
			r:        run{ptr: 4, size: 4},
		},
		{
			title:    "first element",
			arr:      []int{11, 12, 13, 14, 23, 24, 25, 26},
			expected: 4,
			l:        run{ptr: 0, size: 4},
			r:        run{ptr: 4, size: 4},
		},
	}
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.title, func(t *testing.T) {
			is := &intSorter{arr: tc.arr}
			x := findLastElm(is, tc.l, tc.r)
			if x != tc.expected {
				t.Fatalf("%v: expected %v got %v", tc.title, tc.expected, x)
			}
		})
	}
}

func TestMergeRuns(t *testing.T) {
	testCases := []struct {
		title    string
		arr      []testS1
		expected []testS1
		l, r     run
	}{
		{
			title: "one element in left",
			arr: []testS1{
				{1, 11}, {2, 12}, {3, 13}, {4, 15},
				{3, 16}, {4, 17}, {5, 18}, {6, 11},
			},
			expected: []testS1{
				{1, 11}, {2, 12}, {3, 13}, {3, 16},
				{4, 15}, {4, 17}, {5, 18}, {6, 11},
			},
			l: run{ptr: 0, size: 4},
			r: run{ptr: 4, size: 4},
		},
		{
			title: "already sorted",
			arr: []testS1{
				{1, 11}, {2, 12}, {3, 13}, {4, 15}, {5, 16}, {6, 17}, {7, 18},
				{8, 11},
			},
			expected: []testS1{
				{1, 11}, {2, 12}, {3, 13}, {4, 15}, {5, 16}, {6, 17}, {7, 18},
				{8, 11},
			},
			l: run{ptr: 0, size: 4},
			r: run{ptr: 4, size: 4},
		},
		{
			title: "all left moves to temp",
			arr: []testS1{
				{20, 1},
				{11, 11}, {13, 12}, {15, 13}, {17, 14},
				{5, 15}, {14, 16}, {15, 17}, {16, 18},
				{20, 2},
			},
			expected: []testS1{
				{20, 1}, {5, 15}, {11, 11}, {13, 12}, {14, 16}, {15, 13},
				{15, 17}, {16, 18}, {17, 14}, {20, 2},
			},
			l: run{ptr: 1, size: 4},
			r: run{ptr: 5, size: 4},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		if tc.title != "all left moves to temp" {
			continue
		}
		t.Run(tc.title, func(t *testing.T) {
			mergeRuns(&structSorter{arr: tc.arr}, tc.l, tc.r)
			if !reflect.DeepEqual(tc.expected, tc.arr) {
				t.Fatal(tc.arr)
			}
		})
	}
}

func TestMergeStack(t *testing.T) {
	arr := []testS1{
		{1, 11}, {2, 12}, {7, 13}, {9, 15}, {10, 16}, {11, 17},
		{3, 18}, {4, 19}, {8, 20}, {10, 21},
		{3, 22}, {10, 23}, {11, 24}, {12, 25}, {13, 26},
	}
	expected := []testS1{
		{1, 11}, {2, 12}, {3, 18}, {3, 22}, {4, 19}, {7, 13}, {8, 20}, {9, 15},
		{10, 16}, {10, 21}, {10, 23}, {11, 17}, {11, 24}, {12, 25}, {13, 26},
	}
	stack := []run{
		{ptr: 0, size: 6},
		{ptr: 6, size: 4},
		{ptr: 10, size: 5},
	}
	mergeStack(&structSorter{arr: arr}, stack)
	if !reflect.DeepEqual(arr, expected) {
		t.Fatal(arr)
	}
}

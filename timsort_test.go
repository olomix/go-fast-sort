package go_fast_sort

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

type intSorter struct {
	arr []int
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

package sorting

import "sort"

// By returns a sorting wrapper, that would sort by x, doing lexicographical comparisons.
// Example:
// sorting.By(sort.Reverse(sort.Float64Slice(weights))).Sort(sort.StringSlice(names))
// Sorts parallel arrays weights:floats, names:strings by decreasing weights.
func By(x ...sort.Interface) by {
	if len(x) == 0 {
		panic("give me something to sort")
	}
	return by{x}
}

// Sort by `by` reordering what accordingly.
func (b by) Sort(what ...Swapper) {
	if len(what) == 0 {
		sort.Sort(Lex(b.ss...))
	}
	swappers(what).By(b.ss...)
}

// StableSort by `by` reordering what accordingly.
func (b by) StableSort(what ...Swapper) {
	if len(what) == 0 {
		sort.Stable(Lex(b.ss...))
	}
	swappers(what).By(b.ss...)
}

// SortIdx sorts x, returning a corresponding permutation of indices.
// sorting.SortIdx(sort.StringSlice{"b", "c", "a"}) returns []int{2,0,1} and the slice
// is sorted: {"a", "b", "c",
func SortIdx(x ...sort.Interface) []int {
	b := Lex(x...)
	idx := make([]int, b.Len())
	for i := range idx {
		idx[i] = i
	}
	By(b).Sort(sort.IntSlice(idx))
	return idx
}

// Sort returns a sorting wrapper, that would would swap (reorder) what when sorting.
// Example:
// sorting.Sort(sort.Float64Slice(weights)).By(sort.StringSlice(names))
// Sorts parallel arrays weights:floats, names:strings by increasing names.
func Sort(what ...Swapper) swappers {
	return swappers(what)
}

// Swapper is a sub-interface of sort.Interface, that can swap two elements of the collection
type Swapper interface {
	// Swap the elements with indices i and j.
	Swap(i, j int)
}

// SwapFunc is a single-function implementation of the Swapper interface.
type SwapFunc func(i, j int)

// Swap the elements with indices i and j.
func (s SwapFunc) Swap(i, j int) {
	s(i, j)
}

// Lex returns a sort.Interface that corresonds to the lexicographical ordering of x...
func Lex(x ...sort.Interface) sort.Interface {
	switch len(x) {
	case 0:
		panic("give me something to sort")
	case 1:
		return x[0]
	default:
		return lexSorters(x)
	}
}

type by struct {
	ss []sort.Interface
}

// Sort by x, reordering swappers accordingly.
func (s swappers) By(x ...sort.Interface) {
	if len(s) == 1 {
		sorter{Interface: Lex(x...), what: s[0]}.Sort()
		return
	}
	multiSorter{Interface: Lex(x...), what: s}.Sort()
}

type lexSorters []sort.Interface

func (s lexSorters) Swap(i, j int) {
	for _, q := range s {
		q.Swap(i, j)
	}
}

func (s lexSorters) Len() int {
	return s[0].Len()
}

func (s lexSorters) Less(i, j int) bool {
	for _, x := range s {
		if x.Less(i, j) {
			return true
		}
		if x.Less(j, i) {
			return false
		}
	}
	return false
}

// Sorts by Interface, ordering what accordingly.
type sorter struct {
	what Swapper
	sort.Interface
}

func (s sorter) Swap(i, j int) {
	s.Interface.Swap(i, j)
	s.what.Swap(i, j)
}

func (s sorter) Sort() {
	sort.Sort(s)
}

type multiSorter struct {
	what []Swapper
	sort.Interface
}

func (s multiSorter) Swap(i, j int) {
	s.Interface.Swap(i, j)
	for _, q := range s.what {
		q.Swap(i, j)
	}
}

func (s multiSorter) Sort() {
	sort.Sort(s)
}

type swappers []Swapper

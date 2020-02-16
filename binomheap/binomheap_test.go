package binomheap

import "testing"

func TestBionm(t *testing.T) {
	heap := &BinomialHeap{}

	for i := 1; i < 10; i++ {
		if i == 6 || i == 8 {
			heap.ToString()
			heap.Delete(i - 3)
			heap.ToString()
		} else {
			heap.Insert(i)
		}
	}

}

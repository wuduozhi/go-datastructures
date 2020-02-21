package bitree

import "testing"

func TestNewBITree(t *testing.T) {

	freq := [...]int{2, 1, 1, 3, 2, 3, 4, 5, 6, 7, 8, 9}
	max := 9
	bitree := NewBITree(max)
	sum := 0
	for _, value := range freq {
		sum += bitree.GetSum(max-1) - bitree.GetSum(value-1)
		bitree.Update(value-1, 1)
	}

	t.Log(sum)
}

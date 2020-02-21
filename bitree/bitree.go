package bitree

type BITree struct {
	biArray []int
	Size    int
}

func (bit *BITree) Update(index, val int) {
	index++

	for index < len(bit.biArray) {
		bit.biArray[index] += val
		index += index & (-index)
	}
}

func (bit *BITree) GetSum(index int) int {
	index++

	if index > bit.Size {
		panic("index is out of range")
	}

	var sum int
	for index > 0 {
		sum += bit.biArray[index]
		index -= index & (-index)
	}
	return sum
}

func NewBITree(n int) *BITree {
	return &BITree{
		biArray: make([]int, n+1, n+1),
		Size:    n + 1,
	}
}

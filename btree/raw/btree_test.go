package raw

import (
	"fmt"
	"testing"
)

func TestBTree(t *testing.T) {
	tree := newBTree()

	step := 1
	base := 1000000
	for i := base; i < (base + base / 2); i++ {
		tree.Insert(i)
		tree.Insert(base - step)
		step += 1
	}
	tree.ToString()
	tree.Remove(7)
	tree.ToString()
	tree.Remove(6)
	tree.Remove(11)
	tree.Remove(10)
	tree.Remove(9)
	tree.Remove(8)
	tree.Remove(12)
	tree.Remove(13)
	tree.ToString()
	fmt.Println()
	fmt.Println(tree.Search(99))
}

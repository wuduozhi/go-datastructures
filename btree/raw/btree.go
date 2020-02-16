package raw

import (
	"fmt"
)

const t = 5

type BTree struct {
	root *BTreeNode
	t    int
}

func (tree *BTree) Insert(k int) {
	if tree.root == nil {
		tree.root = newBTreeNode(tree.t, true)
		tree.root.keys[0] = k
		tree.root.n = 1
	} else {
		root := tree.root
		if root.n == (2*t - 1) {
			s := newBTreeNode(tree.t, false)
			s.childes[0] = root
			s.SplitChild(0, root)

			i := 0
			if s.keys[0] < k {
				i++
			}
			s.childes[i].InsertNonFull(k)
			tree.root = s
		} else {
			root.InsertNonFull(k)
		}
	}
}

func (tree *BTree) Remove(k int){
	if tree.root == nil {
		fmt.Println("The tree is empty")
		return
	}

	tree.root.Remove(k)
	if tree.root.n == 0 {
		tmpNode := tree.root
		if tmpNode.leaf{
			tree.root = nil
		}else{
			tree.root = tmpNode.childes[0]
		}
		tmpNode = nil
	}
}

func (tree *BTree) ToString() {
	if tree.root == nil {
		fmt.Println("nothing")
		return
	}

	tree.root.Traverse(0)
}

func (tree *BTree) Search(k int) *BTreeNode {
	if tree.root == nil {
		return nil
	}
	return tree.root.Search(k)
}

func newBTree() *BTree {
	return &BTree{t: t}
}

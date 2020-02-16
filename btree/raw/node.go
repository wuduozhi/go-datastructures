package raw

import (
	"fmt"
)

const defaultValue = 0

var defaultNode *BTreeNode = nil

type BTreeNode struct {
	keys    []int // An array of keys
	t       int   // Minimum degree (defines the range for number of keys)
	n       int
	childes []*BTreeNode // An array of child pointers
	leaf    bool
}

func newBTreeNode(t int, leaf bool) *BTreeNode {
	return &BTreeNode{
		t:       t,
		leaf:    leaf,
		childes: make([]*BTreeNode, 2*t, 2*t),
		keys:    make([]int, 2*t-1, 2*t-1),
	}
}

// Function to search key k in subtree rooted with this node
func (root *BTreeNode) Search(k int) *BTreeNode {
	var i = 0
	for i < root.n && k > root.keys[i] {
		i++
	}

	if i < root.n && root.keys[i] == k {
		return root
	}

	if root.leaf {
		return nil
	}

	return root.childes[i].Search(k)
}

func (root *BTreeNode) Traverse(deep int) {
	printN(deep)
	fmt.Println(root.keys)
	if root.leaf == false {
		for _, node := range root.childes {
			if node != nil {
				node.Traverse(deep + 1)
			}
		}
	}

}

func printN(n int) {
	for i := 0; i < n; i++ {
		fmt.Print("-")
	}
}

// A utility function to split the child y of this node
// Note that y must be full when this function is called
func (root *BTreeNode) SplitChild(i int, y *BTreeNode) {
	// Create a new node which is going to store (t-1) keys of y
	z := newBTreeNode(y.t, y.leaf)
	t := root.t
	z.n = t - 1

	// Copy the last (t-1) keys of y to z
	for j := 0; j < t-1; j++ {
		z.keys[j] = y.keys[j+t]
		y.keys[j+t] = 0
	}

	// Copy the last t children of y to z
	if y.leaf == false {
		for j := 0; j < t; j++ {
			z.childes[j] = y.childes[j+t]
			y.childes[j+t] = nil
		}
	}

	// Reduce the number of keys in y
	y.n = t - 1

	// Since this node is going to have a new child,
	// create space of new child
	for j := root.n; j >= i+1; j-- {
		root.childes[j+1] = root.childes[j]
	}
	root.childes[i+1] = z

	// A key of y will move to this node. Find location of
	// new key and move all greater keys one space ahead
	for j := root.n - 1; j >= i; j-- {
		root.keys[j+1] = root.keys[j]
	}
	root.keys[i] = y.keys[t-1]
	y.keys[t-1] = 0

	root.n += 1
}

// A utility function to insert a new key in this node
// The assumption is, the node must be non-full when this
// function is called
func (root *BTreeNode) InsertNonFull(k int) {
	i := root.n - 1

	if root.leaf {
		for i >= 0 && k < root.keys[i] {
			root.keys[i+1] = root.keys[i]
			i--
		}

		root.keys[i+1] = k
		root.n += 1
	} else {
		for i >= 0 && k < root.keys[i] {
			i--
		}

		if root.childes[i+1].n == (2*root.t - 1) {
			root.SplitChild(i+1, root.childes[i+1])

			if root.keys[i+1] < k {
				i++
			}
		}

		root.childes[i+1].InsertNonFull(k)
	}
}

// A function to remove the key k from the sub-tree rooted with this node
func (root *BTreeNode) Remove(k int) {
	idx := root.findKey(k)
	if idx < root.n && root.keys[idx] == k {
		if root.leaf {
			root.removeFromLeaf(idx)
		} else {
			root.removeFromNonLeaf(idx)
		}
	} else {
		if root.leaf {
			fmt.Printf("Key %v does exist in the tree.\n", k)
			return
		}

		// The key to be removed is present in the sub-tree rooted with this node
		// The flag indicates whether the key is present in the sub-tree rooted
		// with the last child of this node
		flag := false
		if idx == root.n {
			flag = true
		}

		if root.childes[idx].n < root.t {
			root.fill(idx)
		}

		// If the last child has been merged, it must have merged with the previous
		// child and so we recurse on the (idx-1)th child. Else, we recurse on the
		// (idx)th child which now has at least t keys
		if flag && idx > root.n {
			root.childes[idx-1].Remove(k)
		} else {
			root.childes[idx].Remove(k)
		}

	}
}

// A function to fill child C[idx] which has less than t-1 keys
func (root *BTreeNode) fill(idx int) {
	// If the previous child(C[idx-1]) has more than t-1 keys, borrow a key
	// from that child
	if idx != 0 && root.childes[idx-1].n >= root.t {
		root.borrowFromPrev(idx)
	} else if idx != root.n && root.childes[idx+1].n >= root.t {
		// If the next child(C[idx+1]) has more than t-1 keys, borrow a key
		// from that child
		root.borrowFromNext(idx)
	} else {
		// Merge C[idx] with its sibling
		// If C[idx] is the last child, merge it with with its previous sibling
		// Otherwise merge it with its next sibling
		if idx != root.n {
			root.merge(idx)
		} else {
			root.merge(idx - 1)
		}
	}
}

func (root *BTreeNode) borrowFromPrev(idx int) {
	preNode := root.childes[idx-1]
	curNode := root.childes[idx]

	for i := curNode.n; i > 0; i-- {
		curNode.keys[i] = curNode.keys[i-1]
		curNode.childes[i+1] = curNode.childes[i]
	}

	// Setting child's first key equal to keys[idx-1] from the current node
	curNode.keys[0] = root.keys[idx-1]
	root.keys[idx-1] = preNode.keys[preNode.n-1]

	curNode.childes[0] = preNode.childes[preNode.n]
	preNode.keys[preNode.n-1] = defaultValue
	preNode.childes[preNode.n] = defaultNode

	curNode.n += 1
	preNode.n -= 1
}

func (root *BTreeNode) borrowFromNext(idx int) {
	nextNode := root.childes[idx-1]
	curNode := root.childes[idx]

	// Setting child's first key equal to keys[idx-1] from the current node
	curNode.keys[curNode.n] = root.keys[idx]
	root.keys[idx] = nextNode.keys[0]

	curNode.childes[curNode.n+1] = nextNode.childes[0]

	for i := 0; i < nextNode.n-1; i++ {
		nextNode.keys[i] = nextNode.keys[i+1]
		nextNode.childes[i] = nextNode.childes[i+1]
	}
	nextNode.keys[nextNode.n-1] = defaultValue
	nextNode.childes[nextNode.n] = defaultNode

	curNode.n += 1
	nextNode.n -= 1
}

// A function to remove the idx-th key from this node - which is a leaf node
func (root *BTreeNode) removeFromLeaf(idx int) {
	for i := idx + 1; i < root.n; i++ {
		root.keys[i-1] = root.keys[i]
	}
	root.keys[root.n-1] = defaultValue
	root.n -= 1
}

// A function to remove the idx-th key from this node - which is a non-leaf node
func (root *BTreeNode) removeFromNonLeaf(idx int) {
	k := root.keys[idx]
	// If the child that precedes k (C[idx]) has atleast t keys,
	// find the predecessor 'pred' of k in the subtree rooted at
	// C[idx]. Replace k by pred. Recursively delete pred
	// in C[idx]
	if root.childes[idx].n >= root.t {
		pred := root.getPred(idx)
		root.keys[idx] = pred
		root.childes[idx].Remove(pred)
	} else if root.childes[idx+1].n >= root.t {
		succ := root.getSucc(idx)
		root.keys[idx] = succ
		root.childes[idx+1].Remove(succ)
	} else {
		root.merge(idx)
		root.childes[idx].Remove(k)
	}
}

// A function to merge C[idx] with C[idx+1]
// C[idx+1] is freed after merging
func (root *BTreeNode) merge(idx int) {
	child := root.childes[idx]
	sibling := root.childes[idx+1]

	// Pulling a key from the current node and inserting it into (t-1)th
	// position of C[idx]
	child.keys[child.n] = root.keys[idx]

	// Copying the keys from C[idx+1] to C[idx] at the end
	// Copying the child pointers from C[idx+1] to C[idx]
	for i := 0; i < sibling.n; i++ {
		child.keys[i+child.t] = sibling.keys[i]
		child.childes[i+child.t] = sibling.childes[i]
	}
	child.childes[2*child.t-1] = sibling.childes[sibling.n]

	// Moving all keys after idx in the current node one step before -
	// to fill the gap created by moving keys[idx] to C[idx]
	for i := idx + 1; i < root.n; i++ {
		root.keys[i-1] = root.keys[i]
		root.childes[i] = root.childes[i+1]
	}

	root.keys[root.n-1] = defaultValue
	root.childes[root.n] = defaultNode
	root.n -= 1

	child.n += sibling.n + 1
	sibling = nil
}

func (root *BTreeNode) getPred(idx int) int {
	cur := root.childes[idx]
	for cur.leaf == false {
		cur = cur.childes[cur.n]
	}

	return cur.keys[cur.n-1]
}

func (root *BTreeNode) getSucc(idx int) int {
	cur := root.childes[idx+1]
	for cur.leaf == false {
		cur = cur.childes[0]
	}

	return cur.keys[0]
}

func (root *BTreeNode) findKey(k int) int {
	var i = 0
	for i < root.n && root.keys[i] < k {
		i++
	}
	return i
}
